package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

type MsgChannelKey struct {
	msgID     string
	channelID string
}

func NewKeyFromString(stringKey string) MsgChannelKey {
	splitted := strings.Split(stringKey, ":")
	return MsgChannelKey{splitted[0], splitted[1]}
}

func (m *MsgChannelKey) GetKey() string {
	return m.channelID + ":" + m.msgID
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: DiscordBotConfigValues.Redis.Endpoint,
		// Username: "default",
		Password: DiscordBotConfigValues.Redis.Password,
		DB:       DiscordBotConfigValues.Redis.DbNum,
	})
}

func KeyGen(channelID string, msgID string) string {
	return fmt.Sprintf("%s:%s", channelID, msgID)
}

func RedisIter(s *discordgo.Session) {
	var cursor uint64
	scnCmd := redisClient.Scan(cursor, "todo:*", 100)
	keys, _, err := scnCmd.Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		keyVal := redisClient.Get(key).Val()
		log.Debugf("Redis: %s -> %s", key, keyVal)
		if val, _ := strconv.Atoi(keyVal); val >= DiscordBotConfigValues.Star.Threshold && val != DiscordBotConfigValues.Redis.DoneVal {
			ScheduleCrossPost(key, s)
		}
	}
}

func SendMessageOnKey(key string, s *discordgo.Session) {
	split_key := strings.Split(key, ":")
	ChanID, MsgID := split_key[2], split_key[1]
	getMessage, err := s.ChannelMessage(ChanID, MsgID)
	if err != nil {
		Hlog(s, err.Error())
		log.Debug(err.Error())
	}
	chanName, err := s.Channel(getMessage.ChannelID)

	if err != nil {
		Hlog(s, fmt.Sprintf("MessageID : %s , ChannelID : %s\nErr: %s", ChanID, MsgID, err.Error()))
		return
	}

	toEmbed := discordgo.MessageEmbed{
		URL: "https://discord.com/channels/" + DiscordBotConfigValues.DiscordConfig.GuildIDs[0] + "/" + getMessage.ChannelID + "/" + getMessage.ID,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    getMessage.Author.GlobalName,
			IconURL: getMessage.Author.AvatarURL(""),
			URL:     "https://discord.com/channels/" + DiscordBotConfigValues.DiscordConfig.GuildIDs[0] + "/" + getMessage.ChannelID + "/" + getMessage.ID,
		},
		Description: getMessage.Content,
	}
	if len(getMessage.Attachments) > 0 {
		// if(getMessage.Attachments[0].ContentType)
		toEmbed.Image = &discordgo.MessageEmbedImage{
			URL: getMessage.Attachments[0].URL,
		}
	}
	s.ChannelMessageSendComplex(
		DiscordBotConfigValues.StarBoardChannel,
		&discordgo.MessageSend{
			Content: fmt.Sprintf(
				"『 x%s 』 %s in %s",
				redisClient.Get(key),
				getMessage.Author.Mention(),
				chanName.Mention(),
			),
			Embed: &toEmbed,
		},
	)
}

func ScheduleCrossPost(key string, s *discordgo.Session) {
	log.Debugf("Crosspost Scheduled! for key %s", key)
	redisClient.Set(key, DiscordBotConfigValues.Redis.DoneVal, 32*time.Hour)
	go SendMessageOnKey(key, s)
}

func ScheduleCrossPostDeletion(key string) {
	log.Debugf("Crosspost Deletion Scheduled! for key %s", key)

	redisClient.Set(key, "", 1) // set big num?
}

func PollingServiceToCrossPost(done chan struct{}, s *discordgo.Session) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		log.Debug("Polling Redis!")
		RedisIter(s)
	}
}

func HandleStarBoardAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	rKey := MsgChannelKey{
		m.ChannelID,
		m.MessageID,
	}
	if m.Emoji.Name == "⭐" || m.Emoji.Name == "✨" || m.Emoji.Name == "❤️" {
		log.Debugf("Incr : %s", rKey.GetKey())
		if val, _ := strconv.Atoi(redisClient.Get(rKey.GetKey()).Val()); val <= DiscordBotConfigValues.Star.Threshold {
			redisClient.Incr("todo:" + rKey.GetKey())
		}

	}
}

func HandleStarBoardDel(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	rKey := MsgChannelKey{
		m.ChannelID,
		m.MessageID,
	}

	if m.Emoji.Name == "⭐" {
		redisClient.Decr("del:" + rKey.GetKey())
	}
}
