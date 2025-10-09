package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func IsTheyAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return false
}

func modLogError(str string, s *discordgo.Session) {
	// s.ChannelMessageSend("1397965285851529296", str)
	log.Warn(str)
}

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
			modLogError(err.Error(), s)
			return
		}
		_, _ = s.ChannelMessageSend(thread.ID, threadName)
		m.ChannelID = thread.ID
	}
}

func HandleMessageInThreads(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID != DiscordBotConfigValues.Channels.PictureChannel {
		return
	}
	chanName, _ := s.Channel(m.ChannelID)
	if len(m.Attachments) > 0 {
		del_in, err := s.ChannelMessageSend(m.ChannelID, "Creating thread...")
		defer s.ChannelMessageDelete(del_in.ChannelID, del_in.ID)
		if err != nil {
			modLogError(err.Error(), s)
		}
		CreateThreadForMessage(s, m)
		time.Sleep(10 * time.Second)

	} else {
		if IsTheyAdmin(s, m) {
			return
		}
		s.MessageReactionAdd(m.ChannelID, m.ID, "pno")
		m_snd, _ := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:   GetRandomThreadRequest(m.Author.DisplayName()),
			Reference: m.Reference(),
			Flags:     discordgo.MessageFlags(discordgo.MessageFlagsEphemeral),
		})

		defer s.ChannelMessageDelete(m.ChannelID, m.ID)
		defer s.ChannelMessageDelete(m_snd.ChannelID, m_snd.ID)
		modLogError(fmt.Sprintf(
			"Deleted Message from %s \n```%s```\nReason: Threads not used! in #%s",
			m.Author.GlobalName,
			m.Content,
			chanName.Name,
		), s)
		time.Sleep(5 * time.Second)
	}
}
