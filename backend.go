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

type UserInfoType struct {
	Id   uint64
	Name string
}

type UserMapType map[uint64]*UserInfoType

type NewUser struct {
	Id   string
	Name string
}

/**
 * ServerInfo
 *
 * Quickly accessible data structure for server information
 */
type ServerInfoType struct {
	Id       uint64
	Name     string
	Dictator string
	UserMap  UserMapType
}

type NewServer struct {
	Id   string
	Name string
}

type ServerMapType map[uint64]*ServerInfoType

type ServerBackend struct {
	Client     *fromilyclient.Client
	ServerInfo ServerMapType
	UserInfo   UserMapType
}

func (b *ServerBackend) Init() {
	b.ServerInfo = ServerMapType{}
	b.UserInfo = UserMapType{}
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
		var serverData = ServerInfoType{}
		serverData.Id = s.Id
		serverData.Name = s.Name
		serverData.UserMap = UserMapType{}

		str := strconv.FormatUint(s.Dictator, 10)
		serverData.Dictator = str

		for _, userdata := range s.Users {
			serverData.UserMap[userdata.User.Id] = &UserInfoType{Name: userdata.User.Name}
		}
		b.ServerInfo[s.Id] = &serverData
	}

	users, err := b.Client.GetUsers()
	if err != nil {
		fmt.Println("Error retrieveing users,", err)
		return false
	}

	for _, u := range users {
		b.UserInfo[u.Id] = &UserInfoType{
			Id:   u.Id,
			Name: u.Name,
		}
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
	var serverData = ServerInfoType{}
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
		b.ServerInfo[sId] = &serverData
		return true
	}
}

func (b *ServerBackend) ServerExists(s string) bool {
	guildId, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return false
	}
	_, ok := b.ServerInfo[guildId]
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

	if server, ok := b.ServerInfo[guildId]; ok == true {
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
