package main

import (
	"fmt"
	"regexp"
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
func DUTIL_ExtractUser(user string) (string, bool) {
	regex := regexp.MustCompile(`<@(\d+)>`)
	parsedUser := regex.FindStringSubmatch(user)
	if parsedUser == nil {
		return "", false
	} else {
		return parsedUser[1], true
	}
}
