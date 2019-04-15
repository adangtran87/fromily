package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var roleBlacklist = []string{
	`@everyone`,
}

// For a particular guild, query all roles and count the number of members with
// said role
func GetRoleCount(s *discordgo.Session, guild string) (roleSummary []*RoleInfo) {
	time := time.Now()
	roles, err := s.GuildRoles(guild)
	if err != nil {
		fmt.Printf("RUMI_ERR: Unable to get roles for Guild: %s", guild)
	}

	roleCount := make(map[string]uint)
	for _, role := range roles {
		roleCount[role.ID] = 0
	}

	last_member := ""
	for {
		members, _ := s.GuildMembers(guild, last_member, 1000)
		if err != nil {
			fmt.Printf("RUMI_ERR: Unable to get members for Guild: %s", guild)
		}
		if len(members) == 0 {
			break
		}

		for _, member := range members {
			for _, role := range member.Roles {
				roleCount[role] = roleCount[role] + 1
			}
		}

		// Get last member
		last_member = members[len(members)-1].User.ID
	}

	for _, role := range roles {
		if !IsRoleBlacklisted(role.Name) {
			var roleInfo RoleInfo
			roleInfo.GuildID = guild
			roleInfo.ID = role.ID
			roleInfo.Name = role.Name
			roleInfo.Count = roleCount[role.ID]
			roleInfo.Time = time
			roleSummary = append(roleSummary, &roleInfo)
		}
	}
	return roleSummary
}

// Discord Utility functions

// Return empty string if invalid user
func DUTIL_ExtractUserMention(user string) string {
	regex := regexp.MustCompile(`<@(\d+)>`)
	parsedUser := regex.FindStringSubmatch(user)
	if parsedUser == nil {
		return ""
	} else {
		return parsedUser[1]
	}
}

// Validate if string is a valid userid; basically is it all digits?
// Return empty string if invalid user
func DUTIL_ValidateUser(user string) string {
	_, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		return ""
	}
	return user
}

func DUTIL_UpdateMember(m *discordgo.Member) {
	// Add users to backend
	user := NewUser{
		Id:   m.User.ID,
		Name: m.User.Username,
	}

	if Backend.UserExists(m.User.ID) == false {
		if Backend.AddUser(m.GuildID, &user) == false {
			fmt.Println("Error creating user: ", user.Id)
		}
	} else {
		go Backend.UpdateUser(&user)
	}

	// Create userdata
	if Backend.UserDataExists(m.GuildID, m.User.ID) == false {
		if Backend.AddUserData(m.GuildID, &user) == false {
			fmt.Println("Error creating userdata:", m.GuildID, m.User.ID)
		}
	}
}

func DUTIL_UpdateGuildInfo(s *discordgo.Session, guild *discordgo.Guild) {
	guildInfo, _ := s.Guild(guild.ID)
	fmt.Printf("%s:%s\n", guildInfo.Name, guildInfo.ID)

	if Backend.ServerExists(guild.ID) {
	} else {
		fmt.Println("Creating guild: ", guild.ID)

		server := NewServer{
			Id:   guild.ID,
			Name: guild.Name,
		}
		if Backend.AddServer(&server) == false {
			fmt.Println("Error creating server: ,", server.Name)
		}
	}

	// Set Admins
	for _, admin := range config.Admins {
		go Backend.SetAdmin(guild.ID, admin)
	}

	// Set user data
	for _, member := range guildInfo.Members {
		DUTIL_UpdateMember(member)
	}
}
