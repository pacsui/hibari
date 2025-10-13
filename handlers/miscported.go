package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func OnMessageOldCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		log.Warnf("Self triggered? --> %s", m.ID)
		return
	}
	if !strings.HasPrefix(m.Content, ">>") {
		return
	}
	switch {
	case strings.HasPrefix(m.Content, ">>av"), strings.HasPrefix(m.Content, ">>avatar"):
		HandleAvatarEmbedReply(s, m)
	case strings.HasPrefix(m.Content, ">>sayas"):
		HandleImpersonation(s, m)
	default:
		log.Warn("Command not found!?")
	}
}

func HandleImpersonation(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Warnf("Not Implemented!")
}

func HandleAvatarEmbedReply(s *discordgo.Session, m *discordgo.MessageCreate) {
	WhoUser := m.Author
	if len(m.Mentions) > 0 {
		WhoUser = m.Mentions[0]
	}

	toEmbed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    WhoUser.GlobalName,
			IconURL: WhoUser.AvatarURL(""),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: WhoUser.AvatarURL(""),
		},
		Color: 0xFFBB22,
	}
	s.ChannelMessageSendEmbedReply(
		m.ChannelID,
		&toEmbed,
		m.MessageReference,
	)
}
