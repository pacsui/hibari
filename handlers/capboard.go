package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func CapBoardHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	switch m.Emoji.Name {
	case "🧢":
		log.Debug("Processing Cap Reaction")
		CapBoardProcessing(s, m)
	}
}

func CapBoardProcessing(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	TakerMessage, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Errorf("this shoudn't happen?")
		return
	}
	GiverUser := m.Emoji.User
	Taker := TakerMessage.Author
	log.Debugf("%s capped %s", GiverUser.ID, Taker.ID)
}
