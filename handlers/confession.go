package handlers

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

var (
	// MsgStore stores the linux epoch of user chan string
	MsgStore = map[string]int64{}
)

// ConfessionVoteDelete : NOP
func ConfessionVoteDelete(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	msgLink, _ := s.ChannelMessage(r.ChannelID, r.MessageID)

	// Check if already sent by using emoji lock lol
	for _, emoji := range msgLink.Reactions {
		if emoji.Emoji.Name == "👍🏽" {
			return
		}
	}

	if r.ChannelID == DiscordBotConfigValues.ConfessionChannel {
		if r.Emoji.Name == "❌" {
			s.MessageReactionAdd(msgLink.ChannelID, msgLink.ID, "👍🏽")
			s.ChannelMessageSend("1397965285851529296", "Confession was reported by "+r.Member.Mention()+"\n[Message URL]("+MessageURL(msgLink.ChannelID, msgLink.ID)+")")
		}
	}
}

// SendConfessionMessage : sends a message in the channel
func SendConfessionMessage(s *discordgo.Session, message string, imgURL string) {
	// Send message, (message Content isnt required to keep it truly anon)
	embedGen := discordgo.MessageEmbed{
		Title: "Anon Confession",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.pacsui.me/imgs/hibari_look.png",
		},
		Description: message,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "react ❌ for reporting",
		},
	}
	if imgURL != "" {
		embedGen.Image = &discordgo.MessageEmbedImage{
			URL: imgURL,
		}
	}
	s.ChannelMessageSendEmbed(
		DiscordBotConfigValues.ConfessionChannel,
		&embedGen,
	)
}

// ConfessionMessageHandler : check if the message in dm and call SendConfessionMessage
func ConfessionMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	whichChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Warn("no channel? how can this be?")
	}
	if strings.HasPrefix(m.Content, C("confess")) {
		if whichChannel.GuildID == "" {
			// No GuildID, its sent in DM so allowed
			messageToSend, _ := strings.CutPrefix(m.Content, C("confess"))
			uepoch, ok := MsgStore[m.ChannelID]
			if ok {
				delta := time.Now().Unix() - uepoch
				if delta < 60 {
					// well not allowed to send message cuz rate limited
					s.MessageReactionAdd(m.ChannelID, m.ID, "⏳")
					log.Debugf("%ds < 60s", delta)
					return
				}
			}
			MsgStore[m.ChannelID] = time.Now().Unix() // store current epoch for the channel
			s.MessageReactionAdd(m.ChannelID, m.ID, "okies:1415595699214618666")
			time.Sleep(time.Second * time.Duration(rand.Intn(10))) //await random time before sending
			imgs := ""
			if len(m.Attachments) > 0 {
				if strings.Contains(m.Attachments[0].ContentType, "image") {
					log.Debug("Confession has an image... forwarding->")
					imgs = m.Attachments[0].URL
				}

			}
			go SendConfessionMessage(s, messageToSend, imgs)

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
