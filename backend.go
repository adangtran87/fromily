/**
 * Server helper functions
 *
 * Interfaces with fromilyclient types
 */
package main

import (
	"fmt"
	"strconv"

	"github.com/adangtran87/fromily/fromilyclient"
)

type UserInfo struct {
	Id   uint64
	Name string
}

type UserMapType map[uint64]*UserInfo

type NewUser struct {
	Id   string
	Name string
}

/**
 * ServerInfo
 *
 * Quickly accessible data structure for server information
 */
type ServerInfo struct {
	Id       uint64
	Name     string
	Dictator string
	UserMap  UserMapType
}

type NewServer struct {
	Id   string
	Name string
}

type ServerMapType map[uint64]*ServerInfo

type ServerBackend struct {
	Client *fromilyclient.Client
	Info   ServerMapType
}

func (b *ServerBackend) Init() {
	b.Info = ServerMapType{}
}

// Loop through an array of pointers
func (b *ServerBackend) RefreshInfo() bool {
	// Check if guild exists
	servers, err := b.Client.GetServers()
	if err != nil {
		fmt.Println("Error retrieveing servers,", err)
		return false
	}

	for _, s := range servers {
		var serverData = ServerInfo{}
		serverData.Id = s.Id
		serverData.Name = s.Name
		serverData.UserMap = UserMapType{}

		str := strconv.FormatUint(s.Dictator, 10)
		serverData.Dictator = str

		for _, userdata := range s.Users {
			serverData.UserMap[userdata.User.Id] = &UserInfo{Name: userdata.User.Name}
		}
		b.Info[s.Id] = &serverData
	}
	return true
}

func (b *ServerBackend) AddServer(n *NewServer) bool {
	sId, err := strconv.ParseUint(n.Id, 10, 64)
	if err != nil {
		fmt.Println("Error converting guild ID into str,", err)
		return false
	}

	// Add server to ServerMapHash
	var serverData = ServerInfo{}
	serverData.Id = sId
	serverData.Name = n.Name
	serverData.UserMap = UserMapType{}
	serverData.Dictator = "0"

	server := fromilyclient.Server{
		Id:   sId,
		Name: n.Name,
	}
	// Create server on server
	err = b.Client.CreateServer(&server)
	if err != nil {
		return false
	} else {
		// Commit to map
		b.Info[sId] = &serverData
		return true
	}
}

func (b *ServerBackend) ServerExists(s string) bool {
	guildId, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return false
	}
	_, ok := b.Info[guildId]
	return ok
}

// Function is expected to work with Discord queries
// So do the conversion inside
func (b *ServerBackend) DictatorExists(s string) bool {
	guildId, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		fmt.Println("Error converting user into str,", err)
		return false
	}

	if server, ok := b.Info[guildId]; ok == true {
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
