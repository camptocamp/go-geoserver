package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// WmsLayerMetadata is a metadata for a wms layer
type WmsLayerMetadata struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",innerxml"`
}

// WmsLayerCRS describes CRS information
type WmsLayerCRS struct {
	Class string `xml:"class,attr,omitempty"`
	Value string `xml:",chardata"`
}

// WmsLayers is a list of WmsLayer
type WmsLayers struct {
	XMLName xml.Name    `xml:"wmsLayers"`
	List    []*WmsLayer `xml:"wmsLayer"`
}

// WmsLayer is a Geoserver object
type WmsLayer struct {
	XMLName           xml.Name            `xml:"wmsLayer"`
	Name              string              `xml:"name"`
	NativeName        string              `xml:"nativeName"`
	Title             string              `xml:"title"`
	Abstract          string              `xml:"abstract"`
	NativeCRS         WmsLayerCRS         `xml:"nativeCRS,omitempty"`
	SRS               string              `xml:"srs"`
	NativeBoundingBox WmsLayerBoundingBox `xml:"nativeBoundingBox"`
	LatLonBoundingBox WmsLayerBoundingBox `xml:"latLonBoundingBox"`
	ProjectionPolicy  string              `xml:"projectionPolicy"`
	Enabled           bool                `xml:"enabled"`
	Metadata          []*WmsLayerMetadata `xml:"metadata>entry,omitempty"`
}

// WmsLayerBoundingBox contains information regarding a wms layer
type WmsLayerBoundingBox struct {
	MinX float64        `xml:"minx"`
	MaxX float64        `xml:"maxx"`
	MinY float64        `xml:"miny"`
	MaxY float64        `xml:"maxy"`
	CRS  FeatureTypeCRS `xml:"crs"`
}

// GetWmsLayers returns all the wms layers
func (c *Client) GetWmsLayers(workspace, wmsstore string) (wmsLayers []*WmsLayer, err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/wmslayers", workspace)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmsstores/%s/wmslayers", workspace, wmsstore)
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

	var data WmsLayers

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmsLayers, err
	}

	for _, wmsLayerRef := range data.List {
		wmsLayer, err := c.GetWmsLayer(workspace, wmsstore, wmsLayerRef.Name)
		if err != nil {
			return wmsLayers, err
		}

		wmsLayers = append(wmsLayers, wmsLayer)
	}

	return
}

// GetWmsLayer return a single wms layer based on its name
func (c *Client) GetWmsLayer(workspace, wmsstore, name string) (wmsLayer *WmsLayer, err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/wmslayers/%s", workspace, name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmsstores/%s/wmslayers/%s", workspace, wmsstore, name)
	}

	statusCode, body, err := c.doRequest("GET", endpoint, nil)
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

	var data WmsLayer
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmsLayer, err
	}

	wmsLayer = &data

	return
}

// CreateWmsLayer creates a WMS Layer
func (c *Client) CreateWmsLayer(workspace string, wmsstore string, wmsLayer *WmsLayer) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/wmslayers", workspace)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmsstores/%s/wmslayers", workspace, wmsstore)
	}

	wmsLayer.XMLName = xml.Name{
		Local: "wmsLayer",
	}
	payload, err := xml.Marshal(wmsLayer)
	if err != nil {
		return
	}
	statusCode, body, err := c.doRequest("POST", endpoint, bytes.NewBuffer(payload))
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

// UpdateWmsLayer updates a wms layer
func (c *Client) UpdateWmsLayer(workspace, wmsstore, wmsLayerName string, wmsLayer *WmsLayer) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/wmslayers/%s", workspace, wmsLayerName)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmsstores/%s/wmslayers/%s", workspace, wmsstore, wmsLayerName)
	}

	wmsLayer.XMLName = xml.Name{
		Local: "wmsLayer",
	}
	payload, _ := xml.Marshal(wmsLayer)

	statusCode, body, err := c.doRequest("PUT", endpoint, bytes.NewBuffer(payload))
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

// DeleteWmsLayer deletes a WMS layer
func (c *Client) DeleteWmsLayer(workspace, wmsstore, wmsLayerName string, recurse bool) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/wmslayers/%s?recurse=%t", workspace, wmsLayerName, recurse)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmsstores/%s/wmslayers/%s?recurse=%t", workspace, wmsstore, wmsLayerName, recurse)
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
