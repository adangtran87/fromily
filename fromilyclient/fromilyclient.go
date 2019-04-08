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
)

type FromilyClient struct {
	Token   string
	BaseUrl string
}

type UserData struct {
	ServerId uint64 `json:"server"`
	Dpoint   uint32 `json:"dpoint"`
}

type ServerData struct {
	UserId uint64 `json:"user"`
	Dpoint uint32 `json:"dpoint"`
}

type FromilyUser struct {
	Id   uint64      `json:"id"`
	Name string      `json:"user_str"`
	Data []*UserData `json:"userdata"`
}

type FromilyServer struct {
	Id   uint64        `json:"id"`
	Name string        `json:"server_str"`
	Data []*ServerData `json:"serverdata"`
}

// Create the Client object
func New(token string) *FromilyClient {
	return &FromilyClient{
		Token:   token,
		BaseUrl: "http://localhost:8000/v1/",
	}
}

func (s *FromilyClient) doRequest(req *http.Request) ([]byte, error) {
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
func (s *FromilyClient) GetServers() ([]*FromilyServer, error) {
	url := fmt.Sprintf(s.BaseUrl + "servers/")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data []*FromilyServer
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
