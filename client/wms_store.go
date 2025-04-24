package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// WmsStoreReference is a reference to a WmsStore
type WmsStoreReference struct {
	Name string `xml:"name"`
}

// WmsStores is a list of WmsStoreReference reference
type WmsStores struct {
	XMLName xml.Name             `xml:"wmsStores"`
	List    []*WmsStoreReference `xml:"wmsStore"`
}

// WmsStore is a Geoserver object
type WmsStore struct {
	XMLName                    xml.Name            `xml:"wmsStore"`
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
	Type             	   string              `xml:"type"`
}

// NewWmsStore creates a new WmsStore with default values
func NewWmsStore() *WmsStore {
	return &WmsStore{
		Type: "WMS",
	}
}

// GetWmsStores returns the list of the wms stores
func (c *Client) GetWmsStores(workspace string) (wmsStores []*WmsStore, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s/wmsstores", workspace), nil)
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

	var data WmsStores
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmsStores, nil
	}

	for _, wmsStoreRef := range data.List {
		wmsStore, err := c.GetWmsStore(workspace, wmsStoreRef.Name)
		if err != nil {
			return wmsStores, err
		}

		wmsStores = append(wmsStores, wmsStore)
	}

	return
}

// GetWmsStore return a single wms store based on its name
func (c *Client) GetWmsStore(workspace, name string) (wmsStore *WmsStore, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/workspaces/%s/wmsstores/%s", workspace, name), nil)
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

	var data WmsStore
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmsStore, err
	}

	wmsStore = &data

	return
}

// CreateWmStore creates a wms store
func (c *Client) CreateWmStore(workspace string, wmsStore *WmsStore) (err error) {
	payload, _ := xml.Marshal(&wmsStore)
	statusCode, body, err := c.doRequest("POST", fmt.Sprintf("/workspaces/%s/wmsstores", workspace), bytes.NewBuffer(payload))
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
func (c *Client) UpdateWmsStore(workspaceName, wmsStoreName string, wmsStore *WmsStore) (err error) {
	payload, _ := xml.Marshal(&wmsStore)

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/workspaces/%s/wmsstores/%s", workspaceName, wmsStoreName), bytes.NewBuffer(payload))
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

// DeleteWmsStore deletes a wms store
func (c *Client) DeleteWmsStore(workspaceName, wmsStoreName string, recurse bool) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/workspaces/%s/wmsstores/%s?recurse=%t", workspaceName, wmsStoreName, recurse), nil)
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
