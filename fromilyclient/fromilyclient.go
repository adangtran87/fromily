/*******************************************************************************
 * Fromily Client
 *
 * Client API to fromily-server
 *******************************************************************************/
package fromilyclient

import (
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
	UserData []UserData `json:"userdata,omitempty"`
}

type DPointRecord struct {
	Points int32    `json:"points"`
	Reason string   `json:"reason"`
	Date   DateType `json:"date"`
}

type UserServerData struct {
	User       uint64          `json:"user"`
	Server     uint64          `json:"server"`
	Dpoints    uint64          `json:"dpoints,omitempty"`
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
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

// Client APIs
func (s *Client) GetServers() ([]*Server, error) {
	url := fmt.Sprintf(s.BaseUrl + "servers/")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data []*Server
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Client) GetServerData(server uint64) ([]*UserServerData, error) {
	url := fmt.Sprintf(s.BaseUrl+"userserverdata/?server=%s", strconv.FormatUint(server, 10))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data []*UserServerData
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
