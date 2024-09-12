package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// LayerGroups is a list of layer groups
type LayerGroups struct {
	XMLName xml.Name      `xml:"layerGroups"`
	List    []*LayerGroup `xml:"layerGroup"`
}

// LayerGroup is a GeoServer resource representing a group of layer
type LayerGroup struct {
	XMLName       xml.Name           `xml:"layerGroup"`
	Name          string             `xml:"name"`
	Workspace     *WorkspaceRef      `xml:"workspace,omitempty"`
	Mode          string             `xml:"mode,omitempty"`
	Title         string             `xml:"title,omitempty"`
	Abstract      string             `xml:"abstractTxt,omitempty"`
	Publishables  []*LayerRef        `xml:"publishables>published,omitempty"`
	Styles        []*StyleRef        `xml:"styles>style,omitempty"`
	Bounds        *BoundingBox       `xml:"bounds,omitempty"`
	MetadataLinks []*MetadataLink    `xml:"metadataLinks>metadataLink,omitempty"`
	Keywords      LayerGroupKeywords `xml:"keywords"`
}

// LayerGroupKeywords is a XML object for Layer Group Keywords
type LayerGroupKeywords struct {
	Keywords []string `xml:"string"`
}

// LayerRef is a reference to an existing layer in GeoServer
type LayerRef struct {
	Type string `xml:"type,attr"`
	Name string `xml:"name"`
}

// StyleRef is a reference to an existing style in GeoServer
type StyleRef struct {
	Name string `xml:"name,omitempty"`
}

// MetadataLink gives informations on external metadata
type MetadataLink struct {
	Type         string `xml:"type"`
	MetadataType string `xml:"metadataType"`
	Content      string `xml:"content"`
}

// GetGroups returns all the groups
func (c *Client) GetGroups(workspace string) (layerGroups []*LayerGroup, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = "/layergroups"
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layergroups", workspace)
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

	var data LayerGroups

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return layerGroups, err
	}

	for _, groupRef := range data.List {
		group, err := c.GetGroup(workspace, groupRef.Name)
		if err != nil {
			return layerGroups, err
		}

		layerGroups = append(layerGroups, group)
	}

	return
}

// GetGroup return a single group based on its name
func (c *Client) GetGroup(workspace, name string) (layerGroup *LayerGroup, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/layergroups/%s", name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layergroups/%s", workspace, name)
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

	var data LayerGroup
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return layerGroup, err
	}

	layerGroup = &data

	return
}

// CreateGroup creates a layer group
func (c *Client) CreateGroup(workspace string, layerGroup *LayerGroup) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = "/layergroups"
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layergroups", workspace)
	}

	layerGroup.XMLName = xml.Name{
		Local: "layerGroup",
	}
	payload, err := xml.Marshal(layerGroup)
	if err != nil {
		return
	}
	statusCode, body, err := c.doFullyTypedRequest("POST", endpoint, bytes.NewBuffer(payload), "application/xml", "")

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

// UpdateGroup updates an existing layer group
func (c *Client) UpdateGroup(workspace string, layerGroup *LayerGroup) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/layergroups/%s", layerGroup.Name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layergroups/%s", workspace, layerGroup.Name)
	}

	layerGroup.XMLName = xml.Name{
		Local: "layerGroup",
	}
	payload, _ := xml.Marshal(layerGroup)

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

// DeleteGroup deletes layer group from GeoServer
func (c *Client) DeleteGroup(workspace string, layerGroup string) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/layergroups/%s", layerGroup)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/layergroups/%s", workspace, layerGroup)
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
