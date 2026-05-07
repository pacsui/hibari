package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const FILENAME string = "member.go"

func ImportMemberHandlers() []Handler {
	return []Handler{
		{
			Name:     "OnMemberJoinHandler",
			Function: OnMemberJoin,
			File:     FILENAME,
		},
		{
			Name:     "OnMemberLeaveHandler",
			Function: OnMemberLeave,
			File:     FILENAME,
		},
		{
			Name:     "OnMessageDelete",
			Function: OnMessageDelete,
			File:     FILENAME,
		},
		{
			Name:     "OnMessageEditDiff",
			Function: OnMessageEdit,
			File:     FILENAME,
		},
	}
}

func OnMemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	s.ChannelMessageSend(DiscordBotConfigValues.ModLogChannel, fmt.Sprintf("%s has joined server", m.Mention()))
}

func OnMemberLeave(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	s.ChannelMessageSend(DiscordBotConfigValues.ModLogChannel, fmt.Sprintf("%s has left server", m.Mention()))
}

func OnMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	content := "not cached"
	author := "not cached user"

	if m.BeforeDelete != nil {
		content = m.BeforeDelete.Content
		if m.BeforeDelete.Author != nil {
			author = m.BeforeDelete.Author.Mention()
		}
	}
	s.ChannelMessageSend(DiscordBotConfigValues.ModLogChannel, fmt.Sprintf("```%s``` was deleted from <#%s> by %s", content, m.ChannelID, author))
}

func OnMessageEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author.ID == s.State.User.ID || (strings.Contains(m.Content, "gif") && strings.Contains(m.Content, "https")) {
		return
	}
	dmp := diffmatchpatch.New()
	if m.BeforeUpdate == nil {
		log.Warnf("no before update edited? : %s", m.Content)
		return
	}
	diffs := dmp.DiffMain(m.BeforeUpdate.Content, m.Content, false)
	s.ChannelMessageSend(DiscordBotConfigValues.ModLogChannel, fmt.Sprintf("%s\n> edited - %s [Link](https://discord.com/channels/%s/%s/%s)", buildAnsiDiff(diffs), m.Author.Mention(), m.GuildID, m.ChannelID, m.ID))
}

func buildAnsiDiff(diffs []diffmatchpatch.Diff) string {
	var sb strings.Builder
	sb.WriteString("```ansi\n")

	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffDelete:
			sb.WriteString(fmt.Sprintf("\x1b[31m%s\x1b[0m", d.Text))
		case diffmatchpatch.DiffInsert:
			sb.WriteString(fmt.Sprintf("\x1b[32m%s\x1b[0m", d.Text))
		case diffmatchpatch.DiffEqual:
			sb.WriteString(d.Text)
		}
	}

	sb.WriteString("\n```")
	return sb.String()
}
