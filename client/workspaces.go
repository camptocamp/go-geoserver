package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// WorkspaceReference is a reference to a Workspace
type WorkspaceReference struct {
	Name string `xml:"name"`
}

// Workspaces is a list of workspace reference
type Workspaces struct {
	XMLName xml.Name              `xml:"workspaces"`
	List    []*WorkspaceReference `xml:"workspace"`
}

// Workspace is a Geoserver object
type Workspace struct {
	XMLName  xml.Name `xml:"workspace"`
	Name     string   `xml:"name"`
	Isolated bool     `xml:"isolated"`
}

// WorkspaceRef is a reference to a GeoServer workspace
type WorkspaceRef struct {
	Name string `xml:"name,omitempty"`
}

// GetWorkspaces returns the list of the workspaces
func (c *Client) GetWorkspaces() (workspaces []*Workspace, err error) {
	statusCode, body, err := c.doRequest("GET", "/workspaces", nil)
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

	var data Workspaces
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return workspaces, nil
	}

	for _, workspaceRef := range data.List {
		workspaces = append(workspaces, &Workspace{
			Name: workspaceRef.Name,
		})
	}

	return
}

// GetWorkspace return a single workspace based on its name
func (c *Client) GetWorkspace(name string) (workspace *Workspace, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s", name), nil)
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
	case 200:
		break
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}

	var data Workspace
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return workspace, nil
	}

	workspace = &data

	return
}

// CreateWorkspace creates a workspace
func (c *Client) CreateWorkspace(workspace *Workspace, isDefault bool) (err error) {
	payload, _ := xml.Marshal(workspace)
	statusCode, body, err := c.doRequest("POST", fmt.Sprintf("/workspaces?default=%t", isDefault), bytes.NewBuffer(payload))
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
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}

// UpdateWorkspace updates a workspace
func (c *Client) UpdateWorkspace(name string, workspace *Workspace) (err error) {
	payload, _ := xml.Marshal(workspace)

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/workspaces/%s", name), bytes.NewBuffer(payload))
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

// DeleteWorkspace deletes a workspace
func (c *Client) DeleteWorkspace(name string, recurse bool) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/workspaces/%s?recurse=%t", name, recurse), nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 403:
		err = fmt.Errorf("workspace is not empty")
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
