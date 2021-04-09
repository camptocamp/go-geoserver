package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// DatastoreRef is a reference to a Datastore
type DatastoreReference struct {
	Name string `xml:"name"`
}

// Datastores is a list of datastore reference
type Datastores struct {
	XMLName xml.Name              `xml:"dataStores"`
	List    []*DatastoreReference `xml:"dataStore"`
}

// DatastoreConnectionParameter is a datastore connection parameter
type DatastoreConnectionParameter struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

// Datastore is a Geoserver object
type Datastore struct {
	XMLName              xml.Name                        `xml:"dataStore"`
	Name                 string                          `xml:"name"`
	Description          string                          `xml:"description"`
	Enabled              bool                            `xml:"enabled"`
	Workspace            *WorkspaceReference             `xml:"workspace"`
	ConnectionParameters []*DatastoreConnectionParameter `xml:"connectionParameters>entry"`
	Default              bool                            `xml:"__default"`
}

// GetDatastores returns the list of the datastores
func (c *Client) GetDatastores(workspace string) (datastores []*Datastore, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s/datastores", workspace), nil)
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

	var data Datastores
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return datastores, nil
	}

	for _, datastoreRef := range data.List {
		datastore, err := c.GetDatastore(workspace, datastoreRef.Name)
		if err != nil {
			return datastores, err
		}

		datastores = append(datastores, datastore)
	}

	return
}

// GetDatastore return a single datastore based on its name
func (c *Client) GetDatastore(workspace, name string) (datastore *Datastore, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s/datastores/%s", workspace, name), nil)
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

	var data Datastore
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return datastore, err
	}

	datastore = &data

	return
}

// CreateDatastore creates a datastore
func (c *Client) CreateDatastore(workspace string, datastore *Datastore) (err error) {
	payload, _ := xml.Marshal(&datastore)
	statusCode, body, err := c.doRequest("POST", fmt.Sprintf("/workspaces/%s/datastores", workspace), bytes.NewBuffer(payload))
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

// UpdateDatastore updates a datastore
func (c *Client) UpdateDatastore(workspaceName, datastoreName string, datastore *Datastore) (err error) {
	payload, _ := xml.Marshal(&datastore)

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/workspaces/%s/datastores/%s", workspaceName, datastoreName), bytes.NewBuffer(payload))
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

// DeleteDatastore deletes a datastore
func (c *Client) DeleteDatastore(workspaceName, datastoreName string, recurse bool) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/workspaces/%s/datastores/%s?recurse=%t", workspaceName, datastoreName, recurse), nil)
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
