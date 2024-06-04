package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// LayerMetadata is a metadata for a Layer
type LayerMetadata struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",innerxml"`
}

// LayerStyle is a reference to a style
type LayerStyle struct {
	Name string `xml:"style>string"`
}

// LayerStyles is a list of style references
type LayerStyles struct {
	Class string        `xml:"class,attr,omitempty"`
	List  []*LayerStyle `xml:"style"`
}

// Resource
type Resource struct {
	Class string `xml:"class,attr,omitempty"`
	Name  string `xml:"name"`
}

// Attribution
type Attribution struct {
	Title           string `xml:"title,omitempty"`
	DataProviderUrl string `xml:"href,omitempty"`
	LogoUrl         string `xml:"logoURL,omitempty"`
	LogoWidth       int    `xml:"logoWidth"`
	LogoHeight      int    `xml:"logoHeight"`
	LogoType        string `xml:"logoType,omitempty"`
}

// Layer is a GeoServer object
type Layer struct {
	XMLName             xml.Name         `xml:"layer"`
	Name                string           `xml:"name"`
	Path                string           `xml:"path,omitempty"`
	Type                string           `xml:"type"`
	DefaultStyle        string           `xml:"defaultStyle>name"`
	Styles              LayerStyles      `xml:"styles,omitempty"`
	LayerResource       Resource         `xml:"resource"`
	IsOpaque            bool             `xml:"opaque,omitempty"`
	Metadata            []*LayerMetadata `xml:"metadata>entry,omitempty"`
	ProviderAttribution Attribution      `xml:"attribution"`
}

// Layers is a list of Layer
type Layers struct {
	XMLName xml.Name `xml:"layers"`
	List    []*Layer `xml:"layer"`
}

// GetLayers returns all the layers
func (c *Client) GetLayers(workspace string) (layers []*Layer, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = "/layers"
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layers", workspace)
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

	var data Layers

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return layers, err
	}

	for _, layerRef := range data.List {
		layer, err := c.GetLayer(workspace, layerRef.Name)
		if err != nil {
			return layers, err
		}

		layers = append(layers, layer)
	}

	return
}

// GetLayer return a single layer based on its name
func (c *Client) GetLayer(workspace, name string) (layer *Layer, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/layers/%s", name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layers/%s", workspace, name)
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

	var data Layer
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return layer, err
	}

	layer = &data

	return
}

// UpdateLayer updates a layer
func (c *Client) UpdateLayer(workspace, layerName string, layer *Layer) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/layers/%s", layerName)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layers/%s", workspace, layerName)
	}

	layer.XMLName = xml.Name{
		Local: "layer",
	}
	payload, _ := xml.Marshal(layer)

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

// DeleteLayer deletes a layer
func (c *Client) DeleteLayer(workspace, layerName string, recurse bool) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/layers/%s?recurse=%t", layerName, recurse)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layers/%s?recurse=%t", workspace, layerName, recurse)
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
