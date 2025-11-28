package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func handleStarBoardStuff(s *discordgo.Session, chanMsgType ChanMsgKeyType, todo int) {
	discMsg, err := s.ChannelMessage(chanMsgType.ChanID, chanMsgType.MsgID)
	if err != nil {
		return
	}
	log.Debugf("Reaction Added to message %s", discMsg.Content)
	starMoji, sparkleMoji := 0, 0
	for _, reaction := range discMsg.Reactions {
		if reaction.Emoji.Name == "⭐" {
			starMoji = reaction.Count
		}
		if reaction.Emoji.Name == "✨" {
			sparkleMoji = reaction.Count
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
				"『 %sx%d 』 %s in %s",
				emoji,
				count,
				getMessage.Author.Mention(),
				chanName.Mention(),
			),
			Embed: &toEmbed,
		},
	)
}
