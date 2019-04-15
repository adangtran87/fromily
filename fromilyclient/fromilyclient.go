/*******************************************************************************
 * Fromily Client
 *
 * Client API to fromily-server
 *******************************************************************************/
package fromilyclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DateType time.Time

func (d *DateType) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateType(t)
	return nil
}

func (d *DateType) Format(s string) string {
	t := time.Time(*d)
	return t.Format(s)
}

type Client struct {
	Token   string
	BaseUrl string
}

type User struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type UserData struct {
	User User `json:"user"`
}

type Server struct {
	Id       uint64     `json:"id"`
	Name     string     `json:"name"`
	Dictator uint64     `json:"dictator,omitempty"`
	Users    []UserData `json:"userdata,omitempty"`
}

type DPointRecord struct {
	Points int32     `json:"points"`
	Reason string    `json:"reason"`
	Date   *DateType `json:"date,omitempty"`
}

type UserServerData struct {
	User       uint64          `json:"user"`
	Server     uint64          `json:"server"`
	Dpoints    int32           `json:"dpoints,omitempty"`
	DPoint_log []*DPointRecord `json:"dpoint_log,omitempty"`
}

// Create the Client object
func New(token string) *Client {
	return &Client{
		Token:   token,
		BaseUrl: "http://localhost:8000/v1/",
	}
}

func (s *Client) doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	} else {
		return nil, fmt.Errorf("%s", body)
	}
}

func (s *Client) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (s *Client) post(url string, j []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", s.Token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = s.doRequest(req)
	return err
}

func (s *Client) put(url string, j []byte) error {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", s.Token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	_, err = s.doRequest(req)
	return err
}

// Client APIs
func (s *Client) GetServers() ([]*Server, error) {
	url := fmt.Sprintf(s.BaseUrl + "servers/")
	bytes, err := s.get(url)
	var data []*Server
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Client) GetServer(server uint64) (*Server, error) {
	url := fmt.Sprintf(s.BaseUrl+"servers/%s/", strconv.FormatUint(server, 10))
	bytes, err := s.get(url)
	var data Server
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Client) GetUsers() ([]*User, error) {
	url := fmt.Sprintf(s.BaseUrl + "users/")
	bytes, err := s.get(url)
	var data []*User
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Client) GetServerData(server uint64) ([]*UserServerData, error) {
	url := fmt.Sprintf(s.BaseUrl+"userserverdata/?server=%s", strconv.FormatUint(server, 10))
	bytes, err := s.get(url)
	var data []*UserServerData
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//@TODO Keep these as strings?
func (s *Client) GetUserServerData(server, user string) (*UserServerData, error) {
	url := fmt.Sprintf(s.BaseUrl+"userserverdata/?user=%s&server=%s", user, server)
	bytes, err := s.get(url)
	var data []UserServerData
	err = json.Unmarshal(bytes, &data)
	return &data[0], err
}

func (s *Client) CreateServer(server *Server) error {
	url := fmt.Sprintf(s.BaseUrl + "servers/")
	j, err := json.Marshal(server)
	if err != nil {
		return err
	}
	return s.post(url, j)
}

func (s *Client) UpdateServer(server *Server) error {
	url := fmt.Sprintf(s.BaseUrl+"servers/%s/", strconv.FormatUint(server.Id, 10))
	j, err := json.Marshal(server)
	if err != nil {
		return err
	}
	return s.put(url, j)
}

func (s *Client) CreateUser(user *User) error {
	url := fmt.Sprintf(s.BaseUrl + "users/")
	j, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.post(url, j)
}

func (s *Client) UpdateUser(user *User) error {
	url := fmt.Sprintf(s.BaseUrl+"users/%s/", strconv.FormatUint(user.Id, 10))
	j, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.put(url, j)
}

func (s *Client) CreateUserServerData(data *UserServerData) error {
	url := fmt.Sprintf(s.BaseUrl + "userserverdata/")
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.post(url, j)
}

func (s *Client) CreateDPointRecord(server string, user string, data *DPointRecord) error {
	url := fmt.Sprintf(s.BaseUrl+"dpoints/?user=%s&server=%s", user, server)
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.post(url, j)
}

func (s *Client) GetLeaderboard(server string) ([]*UserServerData, error) {
	url := fmt.Sprintf(s.BaseUrl+"userserverdata/leaderboard/?server=%s", server)
	bytes, err := s.get(url)
	var data []*UserServerData
	err = json.Unmarshal(bytes, &data)
	return data, err
}
