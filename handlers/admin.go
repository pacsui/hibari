package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func AdminHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != DiscordBotConfigValues.ModChannel {
		return
	}
	CommandSent := strings.TrimPrefix(m.Message.Content, C("admin"))
	switch CommandSent {
	case "announce":
		handleAnnouncement(s, m)
	case "purge":
		handlePurge(s, m)
	case "send":
		arg := strings.TrimPrefix(CommandSent, "send")
		handleSend(s, m, arg)
	default:
		return
	}
}

func handleAnnouncement(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Warn("not impl")
}

func handleSend(s *discordgo.Session, m *discordgo.MessageCreate, a string) {
	// expecting args to be ChannelID::Message to send
	splitted := strings.Split(a, "::")
	if len(splitted) == 2 {
		s.ChannelMessageSend(splitted[0], splitted[1])
	} else {
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
	}
}

func handlePurge(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Warn("not impl")
}
