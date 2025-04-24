package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// WmtsStoreReference is a reference to a WmtsStore
type WmtsStoreReference struct {
	Name string `xml:"name"`
}

// WmtsStores is a list of WmtsStoreReference reference
type WmtsStores struct {
	XMLName xml.Name              `xml:"wmtsStores"`
	List    []*WmtsStoreReference `xml:"wmtsStore"`
}

// WmtsStore is a Geoserver object
type WmtsStore struct {
	XMLName                    xml.Name            `xml:"wmtsStore"`
	Name                       string              `xml:"name"`
	Description                string              `xml:"description"`
	Enabled                    bool                `xml:"enabled"`
	Workspace                  *WorkspaceReference `xml:"workspace"`
	Default                    bool                `xml:"__default"`
	DisableConnectionOnFailure bool                `xml:"disableOnConnFailure"`
	CapabilitiesUrl            string              `xml:"capabilitiesURL"`
	MaxConnections             int                 `xml:"maxConnections"`
	ReadTimeOut                int                 `xml:"readTimeout"`
	ConnectTimeOut             int                 `xml:"connectTimeout"`
	Type                       string              `xml:"type"`
}

// NewWmtsStore creates a new WmtsStore with default values
func NewWmtsStore() *WmtsStore {
	return &WmtsStore{
		Type: "WMTS",
	}
}

// GetWmtsStores returns the list of the wmts stores
func (c *Client) GetWmtsStores(workspace string) (wmtsStores []*WmtsStore, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s/wmtsstores", workspace), nil)
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

	var data WmtsStores
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmtsStores, nil
	}

	for _, wmtsStoreRef := range data.List {
		wmtsStore, err := c.GetWmtsStore(workspace, wmtsStoreRef.Name)
		if err != nil {
			return wmtsStores, err
		}

		wmtsStores = append(wmtsStores, wmtsStore)
	}

	return
}

// GetWmtsStore return a single wms store based on its name
func (c *Client) GetWmtsStore(workspace, name string) (wmtsStore *WmtsStore, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s/wmtsstores/%s", workspace, name), nil)
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

	var data WmtsStore
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmtsStore, err
	}

	wmtsStore = &data

	return
}

// CreateWmStore creates a wmts store
func (c *Client) CreateWmtStore(workspace string, wmtsStore *WmtsStore) (err error) {
	payload, _ := xml.Marshal(&wmtsStore)
	statusCode, body, err := c.doRequest("POST", fmt.Sprintf("/workspaces/%s/wmtsstores", workspace), bytes.NewBuffer(payload))
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

// UpdateWmsStore updates a wms store
func (c *Client) UpdateWmtsStore(workspaceName, wmtsStoreName string, wmtsStore *WmtsStore) (err error) {
	payload, _ := xml.Marshal(&wmtsStore)

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/workspaces/%s/wmtsstores/%s", workspaceName, wmtsStoreName), bytes.NewBuffer(payload))
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

// DeleteWmtsStore deletes a wmts store
func (c *Client) DeleteWmtsStore(workspaceName, wmtsStoreName string, recurse bool) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/workspaces/%s/wmtsstores/%s?recurse=%t", workspaceName, wmtsStoreName, recurse), nil)
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
