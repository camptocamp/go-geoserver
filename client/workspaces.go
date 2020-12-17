package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// WorkspaceRef is a reference to a Workspace
type WorkspaceRef struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

// Workspace is a Geoserver object
type Workspace struct {
	Name           string `json:"name"`
	Isolated       bool   `json:"isolated"`
	DataStores     string `json:"datastores"`
	CoverageStores string `json:"coveragestores"`
	WmsStores      string `json:"wmsstores"`
}

// GetWorkspaces returns the list of the workspaces
func (c *Client) GetWorkspaces() (workspaces []*Workspace, err error) {
	statusCode, body, err := c.doRequest("GET", "/workspaces", nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("Unauthorized")
		return
	case 200:
		break
	default:
		err = fmt.Errorf("Unknown error: %d - %s", statusCode, body)
		return
	}

	var data map[string]map[string][]*WorkspaceRef
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return workspaces, nil
	}

	for _, workspaceRef := range data["workspaces"]["workspace"] {
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
		err = fmt.Errorf("Unauthorized")
		return
	case 404:
		err = fmt.Errorf("Not Found")
		return
	case 200:
		break
	default:
		err = fmt.Errorf("Unknown error: %d - %s", statusCode, body)
		return
	}

	var data map[string]*Workspace
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return workspace, nil
	}

	workspace = data["workspace"]

	return
}

// CreateWorkspace creates a workspace
func (c *Client) CreateWorkspace(workspace *Workspace, isDefault bool) (err error) {
	payload, _ := json.Marshal(map[string]*Workspace{
		"workspace": workspace,
	})
	statusCode, body, err := c.doRequest("POST", fmt.Sprintf("/workspaces?default=%t", isDefault), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("Unauthorized")
		return
	case 201:
		return
	default:
		err = fmt.Errorf("Unknown error: %d - %s", statusCode, body)
		return
	}
}

// UpdateWorkspace updates a workspace
func (c *Client) UpdateWorkspace(name string, workspace *Workspace) (err error) {
	payload, _ := json.Marshal(map[string]*Workspace{
		"workspace": workspace,
	})

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/workspaces/%s", name), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("Unauthorized")
		return
	case 404:
		err = fmt.Errorf("Not Found")
		return
	case 405:
		err = fmt.Errorf("Forbidden")
		return
	case 200:
		return
	default:
		err = fmt.Errorf("Unknown error: %d - %s", statusCode, body)
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
		err = fmt.Errorf("Unauthorized")
		return
	case 403:
		err = fmt.Errorf("Workspace is not empty")
		return
	case 404:
		err = fmt.Errorf("Not Found")
		return
	case 405:
		err = fmt.Errorf("Forbidden")
		return
	case 200:
		return
	default:
		err = fmt.Errorf("Unknown error: %d - %s", statusCode, body)
		return
	}
}
