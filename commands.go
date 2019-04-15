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
	userid := DUTIL_ExtractUserMention(sub)
	if userid == "" {
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

	var fields = []*discordgo.MessageEmbedField{}

	if len(data.DPointLog) > 0 {
		for _, record := range data.DPointLog {
			fieldEntry := &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("__%s__", record.Date),
				Value: fmt.Sprintf("**%s**: *%-s*", record.Points, record.Reason),
			}
			fields = append(fields, fieldEntry)
		}
	}

	embed := &discordgo.MessageEmbed{
		Color:  0xd4d2ff, //@TODO Get color from command set
		Title:  fmt.Sprintf("**Total**: %s", data.DPoints),
		Fields: fields,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func dpoints_give(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
	cmdSlice := strings.SplitN(sub, " ", 3)

	parameters := len(cmdSlice)
	if parameters == 1 {
		// Not enough parameters
		return
	}

	// Check if dictator
	if Backend.IsDictator(m.GuildID, m.Author.ID) == false {
		s.ChannelMessageSend(m.ChannelID, "Begone pleb.")
	}

	// Parse user; if there is a mention extract user from it
	// If there is not a mention, validate ID
	numMentions := len(m.Mentions)
	var user string
	if numMentions > 1 {
		// Do not allow more than one mention in this command.
		return
	} else if numMentions == 1 {
		user = DUTIL_ExtractUserMention(cmdSlice[0])
	} else {
		user = DUTIL_ValidateUser(cmdSlice[0])
	}
	if user == "" {
		// Invalid user
		return
	}

	if Backend.UserDataExists(m.GuildID, user) == false {
		// Not a valid user for this server
		return
	}

	userinfo := Backend.GetUser(user)
	if userinfo == nil {
		// User does not exist for some reason
		return
	}

	var reason string
	if parameters == 2 {
		reason = ""
	} else {
		reason = cmdSlice[2]
	}

	record := DPointRecord{
		Points: cmdSlice[1],
		Reason: reason,
	}

	if Backend.AddDPointRecord(m.GuildID, user, &record) == false {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Could not give points to %s.", userinfo.Name))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Gave %s points to %s.", record.Points, userinfo.Name))
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
