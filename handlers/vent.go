package handlers

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

// SendConfessionMessage : sends a message in the channel
func SendVentMessage(s *discordgo.Session, message string, imgURL string, colHex int) {
	// Send message, (message Content isnt required to keep it truly anon)
	embedGen := discordgo.MessageEmbed{
		Title: "Anon Vent",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://pub-63132de0d0bd4788949b876a878bc482.r2.dev/misc-image/hibari/hibari-concern.png",
		},
		Description: message,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "@Mod for reporting/Help - Color coding to person changes every hour",
		},
		Color: colHex,
	}
	if imgURL != "" {
		embedGen.Image = &discordgo.MessageEmbedImage{
			URL: imgURL,
		}
	}
	_, err := s.ChannelMessageSendEmbed(
		DiscordBotConfigValues.VentChannel,
		&embedGen,
	)
	if err != nil {
		log.Errorf("unable to send vent : %s", err.Error())
	}
}

// ConfessionMessageHandler : check if the message in dm and call SendConfessionMessage
func VentMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	whichChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Warn("no channel? how can this be?")
	}

	if strings.HasPrefix(m.Content, C("vent")) {
		// non Guild members aren't allowed to send message in confession
		_, err = s.GuildMember(DiscordBotConfigValues.DiscordConfig.GuildID, m.Author.ID)
		if err != nil {
			log.Warnf("Non guild member has triggered %s", C("vent"))
			return
		}
		if whichChannel.GuildID == "" {
			messageToSend, _ := strings.CutPrefix(m.Content, C("vent"))

			s.MessageReactionAdd(m.ChannelID, m.ID, "➕")

			if strings.ContainsAny(messageToSend, "@") {
				s.ChannelMessageSend(m.ChannelID, "> Mentions aren't supported btw!")
			}

			time.Sleep(time.Second * time.Duration(rand.Intn(10))) //await random time before sending
			imgs := ""
			if len(m.Attachments) > 0 {
				if strings.Contains(m.Attachments[0].ContentType, "image") {
					log.Debug("Vent has an image... forwarding->")
					imgs = m.Attachments[0].URL
				}

			}
			go SendVentMessage(s, messageToSend, imgs, GenerateColorHex(m.Author.ID))

		} else {
			// Message sent in channel somewhere, ask them to send in DM
			s.MessageReactionAdd(m.ChannelID, m.ID, "‼️")
			delLater, err := s.ChannelMessageSend(
				m.ChannelID,
				"Send Vent in Direct Message `>>vent ...` "+m.Author.Mention(),
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
