package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FeatureTypeAttributeWrapper []*FeatureTypeAttribute

func (w *FeatureTypeAttributeWrapper) UnmarshalJSON(data []byte) (err error) {
	x := bytes.TrimLeft(data, " \t\r\n")
	isArray := len(x) > 0 && x[0] == '['
	isObject := len(x) > 0 && x[0] == '{'

	if isArray {
		var a []*FeatureTypeAttribute
		err = json.Unmarshal(data, &a)
		if err != nil {
			return err
		}
		*w = a
	}

	if isObject {
		var a *FeatureTypeAttribute
		err = json.Unmarshal(data, &a)
		if err != nil {
			return err
		}
		*w = []*FeatureTypeAttribute{a}
	}
	return
}

type CRSWrapper CRS

func (w *CRSWrapper) UnmarshalJSON(data []byte) (err error) {
	x := bytes.TrimLeft(data, " \t\r\n")
	isString := len(x) > 0 && x[0] == '"'
	isObject := len(x) > 0 && x[0] == '{'

	if isString {
		var a string
		err = json.Unmarshal(data, &a)
		if err != nil {
			return err
		}
		w.Class = ""
		w.Value = a
	}

	if isObject {
		var a CRS
		err = json.Unmarshal(data, &a)
		if err != nil {
			return err
		}
		w.Class = a.Class
		w.Value = a.Value
	}
	return
}

func (w *CRSWrapper) MarshalJSON() (payload []byte, err error) {
	var a CRS
	if w.Class == "" {
		a.Value = w.Value
		payload, err = json.Marshal(a.Value)
	} else {
		a.Class = w.Class
		a.Value = a.Value
		payload, err = json.Marshal(a)
	}
	return
}

// FeatureTypeRef is a reference to a Layer
type FeatureTypeRef struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

// FeatureType is a Geoserver object
type FeatureType struct {
	Name              string                    `json:"name"`
	NativeName        string                    `json:"nativeName"`
	Title             string                    `json:"title"`
	Abstract          string                    `json:"abstract"`
	Keywords          map[string][]string       `json:"keywords"`
	NativeCRS         CRSWrapper                `json:"nativeCRS"`
	SRS               string                    `json:"srs"`
	NativeBoundingBox BoundingBox               `json:"nativeBoundingBox"`
	LatLonBoundingBox BoundingBox               `json:"latLonBoundingBox"`
	ProjectionPolicy  string                    `json:"projectionPolicy"`
	Enabled           bool                      `json:"enabled"`
	Attributes        FeatureTypeAttributesList `json:"attributes"`
}

// CRS is a Feature Type object
type CRS struct {
	Class string `json:"@class"`
	Value string `json:"$"`
}

// BoundingBox contains information regarding a featuretype
type BoundingBox struct {
	MinX float64    `json:"minx"`
	MaxX float64    `json:"maxx"`
	MinY float64    `json:"miny"`
	MaxY float64    `json:"maxy"`
	CRS  CRSWrapper `json:"crs"`
}

// FeatureTypeAttributesList contains list of feature type attributes
type FeatureTypeAttributesList struct {
	Attribute FeatureTypeAttributeWrapper `json:"attribute"`
}

// FeatureTypeAttribute is a feature type attribute
type FeatureTypeAttribute struct {
	Name      string `json:"name"`
	MinOccurs int    `json:"minOccurs"`
	MaxOccurs int    `json:"maxOccurs"`
	Nillable  bool   `json:"nillable"`
	Binding   string `json:"binding"`
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
		err = fmt.Errorf("Unauthorized")
		return
	case 200:
		break
	default:
		err = fmt.Errorf("Unknown error: %d - %s", statusCode, body)
		return
	}

	var data map[string]map[string][]*FeatureType
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return featureTypes, nil
	}

	for _, featureTypeRef := range data["featureTypes"]["featureType"] {
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

	var data map[string]*FeatureType
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return featureType, err
	}

	featureType = data["featureType"]

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

	payload, err := json.Marshal(map[string]*FeatureType{
		"featureType": featureType,
	})
	if err != nil {
		fmt.Printf("########%+v", err)
		return
	}
	statusCode, body, err := c.doRequest("POST", endpoint, bytes.NewBuffer(payload))
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
		err = fmt.Errorf("Unknown error: %d - %s - %s", statusCode, body, string(payload))
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

	payload, _ := json.Marshal(map[string]*FeatureType{
		"featureType": featureType,
	})

	fmt.Printf("############## %s", payload)
	statusCode, body, err := c.doRequest("PUT", endpoint, bytes.NewBuffer(payload))
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
