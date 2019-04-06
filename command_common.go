package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

/**
 * CmdType
 *
 * Structure for commands
 *	- funcPtr function to handle the command
 *  - help		help text for the command
 */
type CmdType struct {
	funcPtr func(s *discordgo.Session, m *discordgo.MessageCreate)
	help    string
}

func help(s *discordgo.Session, m *discordgo.MessageCreate, cmd map[string]CmdType) {
	var str strings.Builder
	for key, value := range cmd {
		str.WriteString("`")
		str.WriteString(config.Prefix)
		str.WriteString(key)
		str.WriteString("`")
		str.WriteString(": ")
		str.WriteString(value.help)
		str.WriteString("\n")
	}
	s.ChannelMessageSend(m.ChannelID, str.String())
}
