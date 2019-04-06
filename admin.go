package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var GlobalAdmins = []string{
	"174315979701092352", // Ferko
}

func IsAdmin(user string) bool {
	if InStringArray(user, GlobalAdmins) {
		return true
	}
	// TODO: Check for admins on servers
	return false
}

var adminMap = map[string]CmdType{
	"pong":     CmdType{pong, "Reply with ping!"},
	"roleinfo": CmdType{GetRoleInfo, "Get role information"},
}

// Pong command replies with "Ping!"
func pong(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Ping!")
}

func IsRoleBlacklisted(roleName string) bool {
	return InStringArray(roleName, roleBlacklist)
}

// Get the number of roles and number of members for each role
func GetRoleInfo(s *discordgo.Session, m *discordgo.MessageCreate) {

	guild := m.GuildID

	roleSummary := GetRoleCount(s, guild)

	var strList []string
	for _, summary := range roleSummary {
		str := fmt.Sprintf("**%s**: %d", summary.Name, summary.Count)
		strList = append(strList, str)
	}
	outputStr := strings.Join(strList[:], "\n")

	s.ChannelMessageSend(m.ChannelID, outputStr)
}

func AdminDispatch(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if admin before processing admin command
	if IsAdmin(m.Author.ID) {
		adminCmd := admin_regex.FindStringSubmatch(m.Content)
		if adminCmd != nil {
			// Command must be first match
			adminCmdStr := strings.ToLower(adminCmd[1])
			fmt.Printf("Admin command: %s\n", adminCmdStr)
			if adminCmdStr == "help" {
				go help(s, m, adminMap)
			} else if _, ok := adminMap[adminCmdStr]; ok {
				go adminMap[adminCmdStr].funcPtr(s, m)
			} else {
				s.ChannelMessageSend(m.ChannelID, GetResp("cmd:unknown", adminCmdStr))
			}
		}
	}
}
