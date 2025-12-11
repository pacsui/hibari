package handlers

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func ConfessionVoteDelete(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.ChannelID == DiscordBotConfigValues.ConfessionChannel {
		if r.Emoji.Name == "❌" {
			// do some deletion logic same as starboard calculation? if needed or ask to ping mods?
		}
	}
}

func SendConfessionMessage(s *discordgo.Session, message string) {
	// Send message, (message Content isnt required to keep it truly anon)
	embedGen := discordgo.MessageEmbed{
		Title: "Anon Confession",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.pacsui.me/imgs/hibari_look.png",
		},
		Description: message,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ping Mods to have this deleted",
		},
	}
	s.ChannelMessageSendEmbed(
		DiscordBotConfigValues.ConfessionChannel,
		&embedGen,
	)
}

func ConfessionMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	whichChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Warn("no channel? how can this be?")
	}
	if strings.HasPrefix(m.Content, C("confess")) {
		if whichChannel.GuildID == "" {
			// No GuildID, its sent in DM so allowed
			messageToSend, _ := strings.CutPrefix(m.Content, C("confess"))
			s.MessageReactionAdd(m.ChannelID, m.ID, "okies:1415595699214618666")
			go SendConfessionMessage(s, messageToSend)

		} else {
			// Message sent in channel somewhere, ask them to send in DM
			s.MessageReactionAdd(m.ChannelID, m.ID, "‼️")
			delLater, err := s.ChannelMessageSend(
				m.ChannelID,
				"Send Confessions in Direct Message `>>confess ...` "+m.Author.Mention(),
			)
			go s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.Warn(err.Error())
			}
			time.Sleep(5 * time.Second)
			s.ChannelMessageDelete(delLater.ChannelID, delLater.ID)
		}
	}
}
