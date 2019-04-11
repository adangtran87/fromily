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

type Client struct {
	Token   string
	BaseUrl string
}

type User struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type Server struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
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
