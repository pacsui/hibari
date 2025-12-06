package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func handleStarBoardStuff(s *discordgo.Session, chanMsgType ChanMsgKeyType, todo int) {
	time.Sleep(10 * time.Second) // Sanity sleep if more than Threshold is being hit!
	discMsg, err := s.ChannelMessage(chanMsgType.ChanID, chanMsgType.MsgID)
	if err != nil {
		return
	}
	log.Debugf("Reaction Added to message %s", discMsg.Content)
	starMoji, sparkleMoji := 0, 0
	for _, reaction := range discMsg.Reactions {
		switch reaction.Emoji.Name {
		case "⭐":
			starMoji = reaction.Count
		case "✨":
			sparkleMoji = reaction.Count
		case "🔁":
			return
		}
	}
	log.Debugf("⭐ : %d ; ✨ : %d", starMoji, sparkleMoji)
	if sparkleMoji >= DiscordBotConfigValues.Star.Threshold || starMoji >= DiscordBotConfigValues.Star.Threshold {
		whichEmoji := "⭐"
		whichCount := starMoji
		if sparkleMoji > starMoji {
			whichEmoji = "✨"
			whichCount = sparkleMoji
		}
		SendMessageOnKey(chanMsgType, s, whichEmoji, whichCount)
	}
}

func HandleStarBoardAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	postedChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}
	if postedChannel.NSFW {
		log.Warnf("Starboard emoji triggered in NSFW Channel")
		return
	}

	for _, filteredList := range DiscordBotConfigValues.StarBoardFilteredChannels {
		if m.ChannelID == filteredList {
			log.Warn("Starboard emoji triggered in Filtered Channels")
			return
		}
	}

	if m.Emoji.Name == "⭐" || m.Emoji.Name == "✨" {
		handleStarBoardStuff(s, ChanMsgKeyType{m.ChannelID, m.MessageID}, 1)
	}
}

func HandleStarBoardDel(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.Emoji.Name == "⭐" || m.Emoji.Name == "✨" {
		handleStarBoardStuff(s, ChanMsgKeyType{m.ChannelID, m.MessageID}, -1)
	}
}

func SendMessageOnKey(c ChanMsgKeyType, s *discordgo.Session, emoji string, count int) {
	getMessage, err := s.ChannelMessage(c.ChanID, c.MsgID)
	if err != nil {
		Hlog(s, err.Error())
		log.Debug(err.Error())
	}
	chanName, err := s.Channel(getMessage.ChannelID)

	if err != nil {
		Hlog(s, fmt.Sprintf("MessageID : %s , ChannelID : %s\nErr: %s", c.ChanID, c.MsgID, err.Error()))
		return
	}

	msgURL := "https://discord.com/channels/" + DiscordBotConfigValues.DiscordConfig.GuildIDs[0] + "/" + getMessage.ChannelID + "/" + getMessage.ID

	toEmbed := discordgo.MessageEmbed{
		URL: msgURL,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    getMessage.Author.GlobalName,
			IconURL: getMessage.Author.AvatarURL(""),
			URL:     "https://discord.com/channels/" + DiscordBotConfigValues.DiscordConfig.GuildIDs[0] + "/" + getMessage.ChannelID + "/" + getMessage.ID,
		},
		Description: getMessage.Content + "\n\n[Message Link](" + msgURL + ")",
	}
	if len(getMessage.Attachments) > 0 {
		// if(getMessage.Attachments[0].ContentType)
		toEmbed.Image = &discordgo.MessageEmbedImage{
			URL: getMessage.Attachments[0].URL,
		}
	}
	s.MessageReactionAdd(c.ChanID, c.MsgID, "🔁")
	s.ChannelMessageSendComplex(
		DiscordBotConfigValues.StarBoardChannel,
		&discordgo.MessageSend{
			Content: fmt.Sprintf(
				"/ %s `x%d` / %s - %s",
				emoji,
				count,
				getMessage.Author.Mention(),
				chanName.Mention(),
			),
			Embed: &toEmbed,
		},
	)
}
