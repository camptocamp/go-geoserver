package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"slices"
)

// Users is a list of User
type Users struct {
	XMLName xml.Name `xml:"users"`
	List    []*User  `xml:"user"`
}

// User is GeoServer Resource
type User struct {
	XMLName  xml.Name `xml:"user"`
	Name     string   `xml:"userName"`
	Password string   `xml:"password"`
	Enabled  bool     `xml:"enabled"`
}

// GetUsers returns all the users
func (c *Client) GetUsers(serviceName string) (users Users, err error) {
	var endpoint string

	if serviceName == "" {
		endpoint = "/usergroup/users"
	} else {
		endpoint = fmt.Sprintf("/usergroup/service/%s/users", serviceName)
	}

	statusCode, body, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 200:
		break
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}

	var data Users

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return users, err
	}

	return data, err
}

// GetUser return a single user based on its name
func (c *Client) GetUser(serviceName, userName string) (user *User, err error) {
	users, usersErr := c.GetUsers(serviceName)

	if usersErr != nil {
		return user, usersErr
	}

	userIdx := slices.IndexFunc(users.List, func(user *User) bool { return user.Name == userName })

	if userIdx == -1 {
		return user, fmt.Errorf("User not found")
	}

	user = users.List[userIdx]
	return
}

// CreateUser creates a user
func (c *Client) CreateUser(service string, user *User) (err error) {
	var endpoint string

	if service == "" {
		endpoint = "/usergroup/users"
	} else {
		endpoint = fmt.Sprintf("/usergroup/service/%s/users", service)
	}

	user.XMLName = xml.Name{
		Local: "user",
	}
	payload, err := xml.Marshal(user)
	if err != nil {
		return
	}
	statusCode, body, err := c.doFullyTypedRequest("POST", endpoint, bytes.NewBuffer(payload), "application/xml", "")

	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 201:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s - %s", statusCode, body, string(payload))
		return
	}
}

// UpdateUser creates a user
func (c *Client) UpdateUser(service, userName string, user *User) (err error) {
	var endpoint string

	if service == "" {
		endpoint = fmt.Sprintf("/usergroup/user/%s", userName)
	} else {
		endpoint = fmt.Sprintf("/usergroup/service/%s/user/%s", service, userName)
	}

	user.XMLName = xml.Name{
		Local: "user",
	}
	payload, err := xml.Marshal(user)
	if err != nil {
		return
	}
	statusCode, body, err := c.doFullyTypedRequest("POST", endpoint, bytes.NewBuffer(payload), "application/xml", "")

	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 404:
		err = fmt.Errorf("not found")
		return
	case 405:
		err = fmt.Errorf("forbidden")
		return
	case 200:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}

// DeleteUser deletes user from GeoServer
func (c *Client) DeleteUser(service, userName string) (err error) {
	var endpoint string

	if service == "" {
		endpoint = fmt.Sprintf("/usergroup/user/%s", userName)
	} else {
		endpoint = fmt.Sprintf("/usergroup/service/%s/user/%s", service, userName)
	}

	statusCode, body, err := c.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 403:
		err = fmt.Errorf("service name is not empty")
		return
	case 404:
		err = fmt.Errorf("not found")
		return
	case 405:
		err = fmt.Errorf("forbidden")
		return
	case 200:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}
