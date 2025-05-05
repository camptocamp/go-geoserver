package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// WmtsLayerMetadata is a metadata for a wmts layer
type WmtsLayerMetadata struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",innerxml"`
}

// WmtsLayerCRS describes CRS information
type WmtsLayerCRS struct {
	Class string `xml:"class,attr,omitempty"`
	Value string `xml:",chardata"`
}

// WmtsLayers is a list of WmtsLayer
type WmtsLayers struct {
	XMLName xml.Name     `xml:"wmtsLayers"`
	List    []*WmtsLayer `xml:"wmtsLayer"`
}

// WmtsLayer is a Geoserver object
type WmtsLayer struct {
	XMLName           xml.Name             `xml:"wmtsLayer"`
	Name              string               `xml:"name"`
	NativeName        string               `xml:"nativeName"`
	Title             string               `xml:"title"`
	Abstract          string               `xml:"abstract"`
	NativeCRS         WmtsLayerCRS         `xml:"nativeCRS,omitempty"`
	SRS               string               `xml:"srs"`
	NativeBoundingBox WmtsLayerBoundingBox `xml:"nativeBoundingBox"`
	LatLonBoundingBox WmtsLayerBoundingBox `xml:"latLonBoundingBox"`
	ProjectionPolicy  string               `xml:"projectionPolicy"`
	Enabled           bool                 `xml:"enabled"`
	Metadata          []*WmtsLayerMetadata `xml:"metadata>entry,omitempty"`
}

// WmtsLayerBoundingBox contains information regarding a wmts layer
type WmtsLayerBoundingBox struct {
	MinX float64        `xml:"minx"`
	MaxX float64        `xml:"maxx"`
	MinY float64        `xml:"miny"`
	MaxY float64        `xml:"maxy"`
	CRS  FeatureTypeCRS `xml:"crs"`
}

// GetWmtsLayers returns all the wmts layers
func (c *Client) GetWmtsLayers(workspace, wmtsstore string) (wmtsLayers []*WmtsLayer, err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmtsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/layers", workspace)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmtsstores/%s/layers", workspace, wmtsstore)
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

	var data WmtsLayers

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmtsLayers, err
	}

	for _, wmtsLayerRef := range data.List {
		wmtsLayer, err := c.GetWmtsLayer(workspace, wmtsstore, wmtsLayerRef.Name)
		if err != nil {
			return wmtsLayers, err
		}

		wmtsLayers = append(wmtsLayers, wmtsLayer)
	}

	return
}

// GetWmtsLayer return a single wmts layer based on its name
func (c *Client) GetWmtsLayer(workspace, wmtsstore, name string) (wmtsLayer *WmtsLayer, err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmtsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/layers/%s", workspace, name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmtsstores/%s/layers/%s", workspace, wmtsstore, name)
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

	var data WmtsLayer
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return wmtsLayer, err
	}

	wmtsLayer = &data

	return
}

// CreateWmtsLayer creates a WMTS Layer
func (c *Client) CreateWmtsLayer(workspace string, wmtsstore string, wmtsLayer *WmtsLayer) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmtsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/layers", workspace)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmtsstores/%s/layers", workspace, wmtsstore)
	}

	wmtsLayer.XMLName = xml.Name{
		Local: "wmtsLayer",
	}
	payload, err := xml.Marshal(wmtsLayer)
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

// UpdateWmtsLayer updates a wmts layer
func (c *Client) UpdateWmtsLayer(workspace, wmtsstore, wmtsLayerName string, wmtsLayer *WmtsLayer) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmtsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/layers/%s", workspace, wmtsLayerName)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmtsstores/%s/layers/%s", workspace, wmtsstore, wmtsLayerName)
	}

	wmtsLayer.XMLName = xml.Name{
		Local: "wmtsLayer",
	}
	payload, _ := xml.Marshal(wmtsLayer)

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

// DeleteWmtsLayer deletes a WMS layer
func (c *Client) DeleteWmtsLayer(workspace, wmtsstore, wmtsLayerName string, recurse bool) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if wmtsstore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/layers/%s?recurse=%t", workspace, wmtsLayerName, recurse)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/wmtsstores/%s/layers/%s?recurse=%t", workspace, wmtsstore, wmtsLayerName, recurse)
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
