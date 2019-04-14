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

func (b *ServerBackend) GetServerInfo(s string) (*ServerInfoType, bool) {
	server, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, false
	}
	info, ok := b.ServerInfo[server]
	return info, ok
}

func (b *ServerBackend) ServerExists(s string) bool {
	_, ok := b.GetServerInfo(s)
	return ok
}

func (b *ServerBackend) UserExists(user string) bool {
	uId, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		return false
	}
	_, ok := b.UserInfo[uId]
	return ok
}

func (b *ServerBackend) UserDataExists(server, user string) bool {
	serverInfo, ok := b.GetServerInfo(server)
	if ok == false {
		return false
	}

	userId, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		return false
	}

	_, ok = serverInfo.UserMap[userId]
	return ok
}

func (b *ServerBackend) AddUserData(server string, user *NewUser) bool {
	if b.UserDataExists(server, user.Id) == true {
		return false
	}

	serverInfo, ok := b.GetServerInfo(server)
	if ok == false {
		return false
	}

	userId, err := strconv.ParseUint(user.Id, 10, 64)
	if err != nil {
		return false
	}

	userserverdata := fromilyclient.UserServerData{
		User:   userId,
		Server: serverInfo.Id,
	}
	err = b.Client.CreateUserServerData(&userserverdata)
	if err != nil {
		fmt.Println("Error creating userdata on server,", err)
		return false
	}

	//@FIXME Shouldn't have to create this and should be able to pass it
	//       through...
	userInfo := UserInfoType{
		Id:   userId,
		Name: user.Name,
	}

	// Save to map
	serverInfo.UserMap[userId] = &userInfo
	return true
}

func (b *ServerBackend) AddUser(server string, user *NewUser) bool {
	// Don't add user if userdata exists (because it means user exists)
	if b.UserDataExists(server, user.Id) == true {
		return false
	}

	ok := b.ServerExists(server)
	if ok == false {
		return false
	}

	userId, err := strconv.ParseUint(user.Id, 10, 64)
	if err != nil {
		return false
	}
	userInfo := UserInfoType{
		Id:   userId,
		Name: user.Name,
	}

	// UserServerData does not exist so check if user exists
	if b.UserExists(user.Id) == true {
		return false
	}

	// If user doesn't exist, create user on server
	newUser := fromilyclient.User{
		Id:   userId,
		Name: user.Name,
	}
	err = b.Client.CreateUser(&newUser)
	if err != nil {
		print("Error creating user,", err)
		return false
	}
	b.UserInfo[userId] = &userInfo

	// Create userdata
	ok = b.AddUserData(server, user)

	return ok
}

/*******************************************************************************
 * Dictator
*******************************************************************************/
func (b *ServerBackend) DictatorExists(server string) bool {
	serverInfo, ok := b.GetServerInfo(server)
	if ok == false {
		return false
	}

	if serverInfo.Dictator != "0" {
		return true
	} else {
		return false
	}
}

func (b *ServerBackend) SetDictator(server, user string) bool {
	serverInfo, ok := b.GetServerInfo(server)
	if ok == false {
		// This server does not exist
		return false
	}

	if serverInfo.Dictator == user {
		// User is already Dictator
		return false
	}

	if b.UserDataExists(server, user) == false {
		// Userdata does not exist for this server
		return false
	}

	if b.UserExists(user) == false {
		// User does not exist
		return false
	}

	userId, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		return false
	}

	// Set Dictator
	data, err := b.Client.GetServer(serverInfo.Id)
	if err != nil {
		return false
	}
	data.Dictator = userId
	// Clear out to send less data
	data.Users = []fromilyclient.UserData{}
	err = b.Client.UpdateServer(data)
	if err != nil {
		return false
	}
	return true
}
