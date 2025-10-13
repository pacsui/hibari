package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func CreateThreadForMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	threadName := m.Content
	if len(threadName) < 1 {
		threadName = m.Author.DisplayName() + " | nil ...?"
	}
	if ch, err := s.State.Channel(m.ChannelID); err != nil || !ch.IsThread() {
		thread, err := s.MessageThreadStartComplex(m.ChannelID, m.ID, &discordgo.ThreadStart{
			Name:      threadName[:min(20, len(threadName)-1)],
			Invitable: true,
		})
		if err != nil {
			Hlog(s, err.Error())
			return
		}
		_, _ = s.ChannelMessageSend(thread.ID, "> Automated Created thread! for "+m.Author.Mention())
		m.ChannelID = thread.ID
	}
}

func CheckChannelIDMatches(chanID string) bool {
	matches := 0
	for _, channel := range DiscordBotConfigValues.Channels {
		if channel == chanID {
			log.Debugf("Channel Matched : %s", chanID)
			matches += 1
		}
	}
	return matches > 0
}

func HandleMessageInChannelPool(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !CheckChannelIDMatches(m.ChannelID) {
		return
	}

	if len(m.Attachments) > 0 {
		del_in, err := s.ChannelMessageSend(m.ChannelID, "Creating thread...")
		defer s.ChannelMessageDelete(del_in.ChannelID, del_in.ID)
		if err != nil {
			Hlog(s, err.Error())
		}
		CreateThreadForMessage(s, m)
		time.Sleep(3 * time.Second)
	} else {
		m_snd, _ := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   GetRandomThreadRequest(m.Author.DisplayName()),
			Reference: m.Reference(),
			Flags:     discordgo.MessageFlags(discordgo.MessageFlagsEphemeral),
		})

		defer s.ChannelMessageDelete(m.ChannelID, m.ID)
		defer s.ChannelMessageDelete(m_snd.ChannelID, m_snd.ID)
		time.Sleep(5 * time.Second)
	}
}
