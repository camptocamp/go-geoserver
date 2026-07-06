package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// GwcGsLayer is a Geoserver object
type GwcGsLayer struct {
	XMLName              xml.Name      `xml:"GeoServerLayer"`
	Id                   string        `xml:"id,omitempty"`
	Name                 string        `xml:"name"`
	Enabled              bool          `xml:"enabled"`
	BlobStoreId          string        `xml:"blobStoreId,omitempty"`
	MimeFormats          MimeFormats   `xml:"mimeFormats"`
	GridSubsets          []*GridSubset `xml:"gridSubsets>gridSubset"`
	MetaTileDimensions   []int         `xml:"metaWidthHeight>int"`
	ExpireCacheDuration  int           `xml:"expireCache"`
	ExpireClientDuration int           `xml:"expireClients"`
	GutterSize           int           `xml:"gutter"`
	CacheBypassAllowed   bool          `xml:"cacheBypassAllowed"`
	AutoCacheStyles      bool          `xml:"autoCacheStyles,omitempty"`
}

// GetGridset return a single Gridset based on its name
func (c *Client) GetGwcGsLayer(name string) (layer *GwcGsLayer, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/layers/%s", name), nil)
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
	case 201:
		break
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}

	var data GwcGsLayer
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return layer, err
	}

	layer = &data

	return
}

// CreateGwcWmsLayer creates a GWC GS Layer
func (c *Client) CreateGwcGsLayer(layerName string, layer *GwcGsLayer) (err error) {
	payload, _ := xml.Marshal(&layer)
	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/layers/%s", layerName), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 200:
		return
	case 201:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}

// UpdateGwcWmsLayer updates a GWC GS layer
func (c *Client) UpdateGwcGsLayer(layerName string, layer *GwcGsLayer) (err error) {
	payload, _ := xml.Marshal(&layer)

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/layers/%s", layerName), bytes.NewBuffer(payload))
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

// DeleteGridset deletes a gridset
func (c *Client) DeleteGwcGsLayer(layerName string) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/layers/%s", layerName), nil)
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
