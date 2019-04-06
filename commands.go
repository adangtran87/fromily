package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var cmdMap = map[string]CmdType{
	"ping": CmdType{ping, "Reply with pong!"},
}

// Ping command replies with "Pong!"
func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func CommandDispatch(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmd := prefix_regex.FindStringSubmatch(m.Content)
	if cmd != nil {
		// Command must be first match
		cmdStr := strings.ToLower(cmd[1])
		fmt.Printf("Command: %s\n", cmdStr)
		if cmdStr == "help" {
			go help(s, m, cmdMap)
		} else if _, ok := cmdMap[cmdStr]; ok {
			go cmdMap[cmdStr].funcPtr(s, m)
		} else {
			s.ChannelMessageSend(m.ChannelID, GetResp("cmd:unknown", cmdStr))
		}
	}
}
