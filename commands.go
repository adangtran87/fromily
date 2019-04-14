package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandMap map[string]*Command

type CommandSet struct {
	Name     string
	Prefix   string
	Commands CommandMap
}

type Command struct {
	Name   string
	Cmd    func(s *discordgo.Session, m *discordgo.MessageCreate, sub string)
	Subset *CommandSet
	Help   string
}

var Commands = CommandSet{
	Prefix: "!",
	Commands: CommandMap{
		"ping": &Command{
			Name:   "ping",
			Cmd:    ping,
			Subset: nil,
			Help:   "Reply with Pong!",
		},
		"dictator": &Command{
			Name: "dictator",
			Cmd:  dictator,
			Subset: &CommandSet{
				Prefix: "",
				Commands: CommandMap{
					"set": &Command{
						Name:   "dictator set",
						Cmd:    dictator_set,
						Subset: nil,
						Help:   "Set a dictator",
					},
				},
			},
			Help: "Return the current dictator",
		},
	},
}

// Ping command replies with "Pong!"
func ping(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
	if sub != "" {
		return
	} else {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func dictator(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
	dictator := Backend.GetDictator(m.GuildID)
	if dictator == "" {
		s.ChannelMessageSend(m.ChannelID, "Server has no dictator!")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("All hail, %s!", dictator))
	}
}

func dictator_set(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
	cmdSlice := strings.SplitN(sub, " ", 2)
	if len(cmdSlice) > 1 {
		return
	}
	userid, ok := DUTIL_ExtractUser(sub)
	if ok == false {
		return
	} else {
		if Backend.SetDictator(m.GuildID, userid) {
			dictator(s, m, "")
		}
	}
}

func (cs *CommandSet) Dispatch(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, sub string) {
	// Separate command and text
	cmdSlice := strings.SplitN(sub, " ", 2)

	var cmd string
	if prefix != "" {
		regex := regexp.MustCompile(`^` + prefix + `(\w+)`)
		cmdstr := regex.FindStringSubmatch(cmdSlice[0])
		cmd = cmdstr[1]
		if cmdstr == nil {
			return
		}
	} else {
		cmd = cmdSlice[0]
	}

	if cmd == "help" {
		// go help(s, m, &Commands)
	} else if _, ok := cs.Commands[cmd]; ok {
		if cs.Commands[cmd].Subset == nil || len(cmdSlice) == 1 {
			fmt.Println("Running Command: ", cs.Commands[cmd].Name)

			var remainder string
			if len(cmdSlice) == 1 {
				remainder = ""
			} else {
				remainder = cmdSlice[1]
			}

			go cs.Commands[cmd].Cmd(s, m, remainder)
		} else {
			go cs.Commands[cmd].Subset.Dispatch(s, m, "", cmdSlice[1])
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, GetResp("cmd:unknown", cmd))
	}
}
