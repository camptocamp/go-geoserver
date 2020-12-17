package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DatastoreRef is a reference to a Datastore
type DatastoreRef struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

// Datastore is a Geoserver object
type Datastore struct {
	Name                 string                         `json:"name"`
	Description          string                         `json:"description"`
	Enabled              bool                           `json:"enabled"`
	Workspace            *WorkspaceRef                  `json:"workspace"`
	ConnectionParameters *DatastoreConnectionParameters `json:"connectionParameters"`
	Default              bool                           `json:"__default"`
	FeatureTypes         string                         `json:"featureTypes"`
}

// DatastoreConnectionParameters contains the list of parameters of a connection to a datasource
type DatastoreConnectionParameters struct {
	Entries []*DatastoreEntry `json:"entry"`
}

// Entry is Datastore object
type DatastoreEntry struct {
	Key   string `json:"@key"`
	Value string `json:"$"`
}

// GetDatastores returns the list of the datastores
func (c *Client) GetDatastores(workspace string) (datastores []*Datastore, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s/datastores", workspace), nil)
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

	var data map[string]map[string][]*DatastoreRef
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return datastores, nil
	}

	for _, datastoreRef := range data["dataStores"]["dataStore"] {
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

	var data map[string]*Datastore
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return datastore, err
	}

	datastore = data["dataStore"]

	return
}

// CreateDatastore creates a datastore
func (c *Client) CreateDatastore(workspace string, datastore *Datastore) (err error) {
	payload, _ := json.Marshal(map[string]*Datastore{
		"dataStore": datastore,
	})
	statusCode, body, err := c.doRequest("POST", fmt.Sprintf("/workspaces/%s/datastores", workspace), bytes.NewBuffer(payload))
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

// UpdateDatastore updates a datastore
func (c *Client) UpdateDatastore(workspaceName, datastoreName string, datastore *Datastore) (err error) {
	payload, _ := json.Marshal(map[string]*Datastore{
		"dataStore": datastore,
	})

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/workspaces/%s/datastores/%s", workspaceName, datastoreName), bytes.NewBuffer(payload))
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

// DeleteDatastore deletes a datastore
func (c *Client) DeleteDatastore(workspaceName, datastoreName string, recurse bool) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/workspaces/%s/datastores/%s?recurse=%t", workspaceName, datastoreName, recurse), nil)
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
