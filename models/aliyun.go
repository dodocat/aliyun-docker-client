package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
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

type ProjectUpdateParam struct {
	Description string                 `json:"description"`
	Template    string                 `json:"template"`
	Version     string                 `json:"version"`
	Environment map[string]interface{} `json:"environment"`
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
	BaseURL    *url.URL
	UserAgent  string
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
		err = errors.New("status:" + resp.Status)
	}
	return
}

// GetProject Get a project by name
func (c *AliClient) GetProject(name string) (project Project, err error) {
	rel := &url.URL{Path: fmt.Sprint("/projects/%s", name)}
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
	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(&project)
	} else {
		err = errors.New("status:" + resp.Status)
	}
	return
}

func (c *AliClient) UpdateProject(projectName string, updateBody ProjectUpdateParam) (err error) {
	return
}

func (c *AliClient) UpdateService(projectName string, serviceName string, newImage string) (updatedService Service, err error) {
	project, err := c.GetProject(projectName)
	if err != nil {
		return
	}

	template := project.Template
	v := viper.New()
	v.SetConfigType("yml")
	v.ReadConfig(bytes.NewReader([]byte(template)))
	oldImage := v.GetString(fmt.Sprintf("services.%s.image", serviceName))
	if len(oldImage) <= 0 {
		err = errors.New("image path error.")
		return
	}
	project.Template = strings.Replace(template, oldImage, newImage, -1)

	return
}
