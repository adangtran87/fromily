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
	Admin  bool
	Cmd    func(s *discordgo.Session, m *discordgo.MessageCreate, sub string)
	Subset *CommandSet
	Help   string
}

var Commands = CommandSet{
	Prefix: "!",
	Commands: CommandMap{
		"ping": &Command{
			Admin:  false,
			Name:   "ping",
			Cmd:    ping,
			Subset: nil,
			Help:   "Reply with Pong!",
		},
		"dictator": &Command{
			Admin: false,
			Name:  "dictator",
			Cmd:   dictator,
			Subset: &CommandSet{
				Prefix: "",
				Commands: CommandMap{
					"set": &Command{
						Admin:  true,
						Name:   "dictator set",
						Cmd:    dictator_set,
						Subset: nil,
						Help:   "Set a dictator",
					},
				},
			},
			Help: "Return the current dictator",
		},
		"dpoints": &Command{
			Admin: false,
			Name:  "dpoints",
			Cmd:   dpoints,
			Subset: &CommandSet{
				Prefix: "",
				Commands: CommandMap{
					"give": &Command{
						Admin:  true,
						Name:   "dpoints give",
						Cmd:    dpoints_give,
						Subset: nil,
						Help:   "Give the dpoints :eyes:",
					},
				},
			},
			Help: "Return the number of dpoints for the user",
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

func dpoints(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
	//@TODO Support querying other people
	if sub != "" {
		return
	}
	data := Backend.GetUserData(m.GuildID, m.Author.ID)
	if data == nil {
		return
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("DPoints: %s", data.DPoints))

	if len(data.DPointLog) > 0 {
		out.WriteString("\n\nLast 5 Records:\n")
	}

	for _, record := range data.DPointLog {
		out.WriteString(fmt.Sprintf("- %s: %s \"%s\"\n", record.Date, record.Points, record.Reason))
	}

	s.ChannelMessageSend(m.ChannelID, out.String())
}

func dpoints_give(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
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
		command := cs.Commands[cmd]
		if command.Subset == nil || len(cmdSlice) == 1 {
			fmt.Println("Running Command: ", command.Name)

			var remainder string
			if len(cmdSlice) == 1 {
				remainder = ""
			} else {
				remainder = cmdSlice[1]
			}

			_, adminStatus := Backend.IsAdmin(m.GuildID, m.Author.ID)
			if command.Admin == true && adminStatus == false {
				s.ChannelMessageSend(m.ChannelID, ":no_entry_sign: Ah ah ah, you didn't say the magic word. :no_entry_sign:")
			} else {
				go command.Cmd(s, m, remainder)
			}
		} else {
			go command.Subset.Dispatch(s, m, "", cmdSlice[1])
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, GetResp("cmd:unknown", cmd))
	}
}
