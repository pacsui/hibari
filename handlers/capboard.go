package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

func CapBoardHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd, r *redis.Client) {
	switch m.Emoji.Name {
	case "🧢":
		log.Debug("Processing Cap Reaction")
		CapBoardProcessing(s, m, r)
	}
}

func CapBoardCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, r *redis.Client) {
	if strings.HasPrefix(m.Content, C("caps")) {
		if len(m.Mentions) <= 0 {
			embedGen := CapBoardEmbedPreview(s, r, m.Author.ID)
			if embedGen == nil {
				log.Errorf("unable to generate embed? value nils?")
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embedGen)
		}
	}
}

func CapBoardEmbedPreview(s *discordgo.Session, r *redis.Client, UID string) *discordgo.MessageEmbed {
	disUser, err := s.User(UID)
	if err != nil {
		log.Errorf("unable to fetch disUser. %s", err.Error())
		return nil
	}

	CapsGiven, CapsRecv := GetCapGiven(r, UID), GetCapRecv(r, UID)

	emb := discordgo.MessageEmbed{
		Title: fmt.Sprintf("CapStats for %s", disUser.GlobalName),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: disUser.AvatarURL("16"),
		},
		Color:       0xFF5656,
		Description: fmt.Sprintf("```given : %s\ngotten : %s\n```", CapsGiven, CapsRecv),
	}
	return &emb
}

func GetCapGiven(r *redis.Client, UID string) string {
	caps := r.Get(context.TODO(), "cap:"+UID+":giv").Val()
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return ""
	// }
	return caps
}
func GetCapRecv(r *redis.Client, UID string) string {
	caps := r.Get(context.TODO(), "cap:"+UID+":recv").Val()
	// if err != redis.Nil {
	// 	log.Error(err.Error())
	// 	return ""
	// }
	return caps
}
func GetCapGivenRecent(r *redis.Client, UID string) string {
	log.Debug(UID)
	caps, err := r.Get(context.TODO(), "cap:"+UID+":giv_prev").Result()
	if err != redis.Nil {
		log.Error(err.Error())
		return ""
	}
	return caps
}
func GetCapRecvRecent(r *redis.Client, UID string) string {
	log.Debug(UID)
	caps, err := r.Get(context.TODO(), "cap:"+UID+":recv_prev").Result()
	if err != redis.Nil {
		log.Error(err.Error())
		return ""
	}
	return caps
}

func CapBoardProcessing(s *discordgo.Session, m *discordgo.MessageReactionAdd, r *redis.Client) {
	TakerMessage, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	log.Debug(m.ChannelID)
	if err != nil {
		log.Errorf("this shoudn't happen?")
		return
	}
	GiverUser := m.MessageReaction.UserID
	Taker := TakerMessage.Author.ID
	// if GiverUser == Taker {
	// 	log.Debugf("Self Capping Ignored! %s", GiverUser)
	// 	return
	// }
	log.Debugf("%s capped %s", GiverUser, Taker)
	{
		if err := r.Incr(context.TODO(), "cap:"+GiverUser+":giv").Err(); err != nil {
			log.Error(err.Error())
			r.Del(context.Background(), "cap:"+GiverUser+":giv")
			return
		}
		if err := r.Set(context.TODO(), "cap:"+GiverUser, ":giv_prev:"+Taker, 0).Err(); err != nil {
			log.Error(err.Error())
			r.Del(context.Background(), "cap:"+GiverUser, ":giv_prev:"+Taker)
			return
		}

		if err := r.Incr(context.TODO(), "cap:"+Taker+":recv").Err(); err != nil {
			log.Error(err.Error())
			r.Del(context.Background(), "cap:"+Taker+":recv")
			return
		}
		if err := r.Set(context.TODO(), "cap:"+Taker, ":recv_prev:"+GiverUser, 0).Err(); err != nil {
			log.Error(err.Error())
			r.Del(context.Background(), "cap:"+Taker, ":recv_prev:"+GiverUser)
			return
		}
	}

}
