package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// FeatureTypeMetadata is a metadata for a Feature Type
type FeatureTypeMetadata struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",innerxml"`
}

// FeatureTypeKeywords is a XML object for Feature Type Keywords
type FeatureTypeKeywords struct {
	Keywords []string `xml:"string"`
}

type FeatureTypeCRS struct {
	Class string `xml:"class,attr,omitempty"`
	Value string `xml:",chardata"`
}

type FeatureTypes struct {
	XMLName xml.Name       `xml:"featureTypes"`
	List    []*FeatureType `xml:"featureType"`
}

// FeatureType is a Geoserver object
type FeatureType struct {
	XMLName           xml.Name                `xml:"featureType"`
	Name              string                  `xml:"name"`
	NativeName        string                  `xml:"nativeName"`
	Title             string                  `xml:"title"`
	Abstract          string                  `xml:"abstract"`
	Keywords          FeatureTypeKeywords     `xml:"keywords"`
	NativeCRS         FeatureTypeCRS          `xml:"nativeCRS,omitempty"`
	SRS               string                  `xml:"srs"`
	NativeBoundingBox BoundingBox             `xml:"nativeBoundingBox"`
	LatLonBoundingBox BoundingBox             `xml:"latLonBoundingBox"`
	ProjectionPolicy  string                  `xml:"projectionPolicy"`
	Enabled           bool                    `xml:"enabled"`
	Attributes        []*FeatureTypeAttribute `xml:"attributes>attribute"`
	Metadata          []*FeatureTypeMetadata  `xml:"metadata>entry,omitempty"`
}

// BoundingBox contains information regarding a featuretype
type BoundingBox struct {
	MinX float64        `xml:"minx"`
	MaxX float64        `xml:"maxx"`
	MinY float64        `xml:"miny"`
	MaxY float64        `xml:"maxy"`
	CRS  FeatureTypeCRS `xml:"crs"`
}

// FeatureTypeAttribute is a feature type attribute
type FeatureTypeAttribute struct {
	Name      string `xml:"name"`
	MinOccurs int    `xml:"minOccurs"`
	MaxOccurs int    `xml:"maxOccurs"`
	Nillable  bool   `xml:"nillable"`
	Binding   string `xml:"binding"`
}

// GetFeatureTypes returns all the layers
func (c *Client) GetFeatureTypes(workspace, datastore string) (featureTypes []*FeatureType, err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if datastore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/featuretypes", workspace)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/datastores/%s/featuretypes", workspace, datastore)
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

	var data FeatureTypes

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return featureTypes, err
	}

	for _, featureTypeRef := range data.List {
		featureType, err := c.GetFeatureType(workspace, datastore, featureTypeRef.Name)
		if err != nil {
			return featureTypes, err
		}

		featureTypes = append(featureTypes, featureType)
	}

	return
}

// GetFeatureType return a single featuretype based on its name
func (c *Client) GetFeatureType(workspace, datastore, name string) (featureType *FeatureType, err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if datastore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/featuretypes/%s", workspace, name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/datastores/%s/featuretypes/%s", workspace, datastore, name)
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

	var data FeatureType
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return featureType, err
	}

	featureType = &data

	return
}

// CreateFeatureType creates a Feature Type
func (c *Client) CreateFeatureType(workspace string, datastore string, featureType *FeatureType) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if datastore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/featuretypes", workspace)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/datastores/%s/featuretypes", workspace, datastore)
	}

	featureType.XMLName = xml.Name{
		Local: "featureType",
	}
	payload, err := xml.Marshal(featureType)
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

// UpdateFeatureType updates a featuretype
func (c *Client) UpdateFeatureType(workspace, datastore, featureTypeName string, featureType *FeatureType) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if datastore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/featuretypes/%s", workspace, featureTypeName)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/datastores/%s/featuretypes/%s", workspace, datastore, featureTypeName)
	}

	featureType.XMLName = xml.Name{
		Local: "featureType",
	}
	payload, _ := xml.Marshal(featureType)

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

// DeleteFeatureType deletes a datastore
func (c *Client) DeleteFeatureType(workspace, datastore, featureType string, recurse bool) (err error) {
	var endpoint string
	if workspace == "" {
		err = fmt.Errorf("workspace cannot be null")
		return
	}

	if datastore == "" {
		endpoint = fmt.Sprintf("/workspaces/%s/featuretypes/%s?recurse=%t", workspace, featureType, recurse)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/datastores/%s/featuretypes/%s?recurse=%t", workspace, datastore, featureType, recurse)
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
