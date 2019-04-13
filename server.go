/**
 * Server helper functions
 *
 * Interfaces with fromilyclient types
 */
package main

import (
	"strconv"

	"github.com/adangtran87/fromily/fromilyclient"
)

type UserInfo struct {
	Name string
}

type UserMapType map[uint64]*UserInfo

/**
 * ServerInfo
 *
 * Quickly accessible data structure for server information
 */
type ServerInfo struct {
	Name     string
	Dictator string
	UserMap  UserMapType
}

type ServerMapType map[uint64]*ServerInfo

// map of what is on the server
var ServerMap = ServerMapType{}

// Loop through an array of pointers
func (m ServerMapType) ProcessDataIntoServerMap(servers []*fromilyclient.Server) {
	for _, s := range servers {
		var serverData = ServerInfo{}
		serverData.Name = s.Name
		serverData.UserMap = UserMapType{}

		str := strconv.FormatUint(s.Dictator, 10)
		serverData.Dictator = str

		for _, userdata := range s.Users {
			serverData.UserMap[userdata.User.Id] = &UserInfo{Name: userdata.User.Name}
		}
		m[s.Id] = &serverData
	}
}

func (m ServerMapType) AddServer(s *fromilyclient.Server) error {
	// Add server to ServerMapHash
	var serverData = ServerInfo{}
	serverData.Name = s.Name
	serverData.UserMap = UserMapType{}

	str := strconv.FormatUint(s.Dictator, 10)
	serverData.Dictator = str

	ServerMap[s.Id] = &serverData
	// Create server on server
	err := Fromily.CreateServer(s)
	return err
}

// Function is expected to work with Discord queries
// So do the conversion inside
func (m ServerMapType) DictatorExists(s string) bool {
	guildId, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return false
	}
	if server, ok := m[guildId]; ok == true {
		if server.Dictator != "0" {
			return true
		} else {
			return false
		}
	} else {
		// No server so return false
		return false
	}
}

func (m ServerMapType) ServerExists(s string) bool {
	guildId, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return false
	}
	_, ok := ServerMap[guildId]
	return ok
}
