package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Project aliyun cs project
type Project struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Template     string                 `json:"template"`
	Version      string                 `json:"version"`
	Created      string                 `json:"created"`
	Updated      string                 `json:"updated"`
	DesiredState string                 `json:"desired_state"`
	CurrentState string                 `json:"current_state"`
	Environment  map[string]interface{} `json:"environment"`
	Services     []Service              `json:"services"`
}

// Service aliyun cs service
type Service struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Project      string                 `json:"project"`
	Description  string                 `json:"description"`
	Created      string                 `json:"created"`
	Updated      string                 `json:"updated"`
	DesiredState string                 `json:"desired_state"`
	CurrentState string                 `json:"current_state"`
	Definition   map[string]interface{} `json:"definition"`
	Extensions   map[string]interface{} `json:"extensions"`
	Containers   map[string]interface{} `json:"containers"`
}

// Client aliyun client
type AliClient struct {
	BaseURL   *url.URL
	UserAgent string

	HttpClient *http.Client
}

// ListProject Get Project list
func (c *AliClient) ListProject() (projects []Project, err error) {

	rel := &url.URL{Path: "/projects/"}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	//	data, err := ioutil.ReadAll(resp.Body)
	//	println(string(data))

	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(&projects)
	} else {
		err = errors.New("status code:" + resp.Status)
	}
	return
}
