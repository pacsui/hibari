package handlers

import (
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func AmIMentioned(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			log.Debug("I am mentioned!")
			return true
		}
	}
	return false
}

func OnMessageOldCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		log.Warnf("Self triggered? --> %s", m.ID)
		return
	}
	if strings.Contains(strings.ToLower(m.Content[:min(len(m.Content), 50)]), "hello") || AmIMentioned(s, m) {
		reply := DiscordBotConfigValues.HelloReply[rand.Intn(len(DiscordBotConfigValues.HelloReply))]
		s.ChannelMessageSend(m.ChannelID, reply)
	}

	if strings.Contains(strings.ToLower(m.Content[:min(len(m.Content), 50)]), "job") {
		go s.MessageReactionAdd(m.ChannelID, m.ID, "💀")
	}

	if !strings.HasPrefix(m.Content, ">>") {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, ">>av"), strings.HasPrefix(m.Content, ">>avatar"):
		go HandleAvatarEmbedReply(s, m)
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
	log.Debugf("Handling Avatar command for %s", m.Content)
	WhoUser := m.Author
	if len(m.Mentions) > 0 {
		log.Debugf("Mentioned! %s", m.Mentions[0])
		WhoUser = m.Mentions[0]
	}

	toEmbed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    WhoUser.GlobalName,
			IconURL: WhoUser.AvatarURL("16"),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: WhoUser.AvatarURL("512"),
		},
		Color: 0x77FF77,
	}
	// s.ChannelMessageSendEmbedReply(
	// 	m.ChannelID,
	// 	&toEmbed,
	// 	m.MessageReference,
	// )
	s.ChannelMessageSendEmbed(
		m.ChannelID,
		&toEmbed,
	)
}
