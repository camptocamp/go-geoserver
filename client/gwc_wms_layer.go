package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// GwcLayerReference is a reference to a Gridset
type GwcLayerReference struct {
	Name string `xml:"name"`
}

// GwcLayers is a list of gwc layer references
type GwcLayers struct {
	XMLName xml.Name             `xml:"layers"`
	List    []*GwcLayerReference `xml:"layer"`
}

// GwcWmsLayer is a Geoserver object
type GwcWmsLayer struct {
	XMLName              xml.Name      `xml:"wmsLayer"`
	Name                 string        `xml:"name"`
	Enabled              bool          `xml:"enabled"`
	BlobStoreId          string        `xml:"blobStoreId"`
	MimeFormats          MimeFormats   `xml:"mimeFormats"`
	GridSubsets          []*GridSubset `xml:"gridSubsets>gridSubset"`
	MetaTileDimensions   []int         `xml:"metaWidthHeight>int"`
	ExpireCacheDuration  int           `xml:"expireCache"`
	ExpireClientDuration int           `xml:"expireClients"`
	GutterSize           int           `xml:"gutter"`
	BackendTimeout       int           `xml:"backendTimeout"`
	CacheBypassAllowed   bool          `xml:"cacheBypassAllowed"`
	WmsUrl               string        `xml:"wmsUrl>string"`
	WmsLayer             string        `xml:"wmsLayers"`
	WmsVersion           string        `xml:"wmsVersion,omitempty"`
	VendorParameters     string        `xml:"vendorParameters,omitempty"`
	Transparent          bool          `xml:"transparent"`
	BgColor              string        `xml:"bgColor,omitempty"`
}

type GridSubset struct {
	Name          string `xml:"gridSetName"`
	MinCacheLevel int    `xml:"minCachedLevel,omitempty"`
	MaxCacheLevel int    `xml:"maxCachedLevel,omitempty"`
}

// ScaleNames is a XML object for scale names
type MimeFormats struct {
	Formats []string `xml:"string"`
}

// GetGridset return a single Gridset based on its name
func (c *Client) GetGwcWMSLayer(name string) (layer *GwcWmsLayer, err error) {
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

	var data GwcWmsLayer
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return layer, err
	}

	layer = &data

	return
}

// CreateGwcWmsLayer creates a GWC WMS Layer
func (c *Client) CreateGwcWmsLayer(layerName string, layer *GwcWmsLayer) (err error) {
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

// UpdateGwcWmsLayer updates a GWC wms layer
func (c *Client) UpdateGwcWmsLayer(layerName string, layer *GwcWmsLayer) (err error) {
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
func (c *Client) DeleteGwcWmsLayer(layerName string) (err error) {
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
