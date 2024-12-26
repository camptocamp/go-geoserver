package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// GridsetReference is a reference to a Gridset
type GridsetReference struct {
	Name string `xml:"name"`
}

// Gridsets is a list of gridset reference
type Gridsets struct {
	XMLName xml.Name            `xml:"gridSets"`
	List    []*GridsetReference `xml:"gridSet"`
}

// Gridset is a Geoserver object
type Gridset struct {
	XMLName           xml.Name          `xml:"gridSet"`
	Name              string            `xml:"name"`
	Description       string            `xml:"description"`
	AlignTopLeft      bool              `xml:"alignTopLeft"`
	MetersPerUnit     float64           `xml:"metersPerUnit"`
	PixelSize         float64           `xml:"pixelSize"`
	TileHeight        int               `xml:"tileHeight"`
	TileWidth         int               `xml:"tileWidth"`
	YCoordinateFirst  bool              `xml:"yCoordinateFirst"`
	Extent            []float64         `xml:"extent>coords>double"`
	ScaleNames        ScaleNames        `xml:"scaleNames"`
	ScaleDenominators ScaleDenominators `xml:"scaleDenominators"`
	Srs               SRS               `xml:"srs"`
}

// SRS is a XML object for SRS
type SRS struct {
	SrsNumber int `xml:"number"`
}

// ScaleNames is a XML object for scale names
type ScaleNames struct {
	ScaleName []string `xml:"string"`
}

// ScaleNames is a XML object for scale names
type ScaleDenominators struct {
	ScaleDenominator []float64 `xml:"double"`
}

// GetGridsets returns all the gridsets
func (c *Client) GetGridsets() (gridsets []*Gridset, err error) {
	var endpoint string = "/gridsets"

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

	var data Gridsets

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return gridsets, err
	}

	for _, gridsetRef := range data.List {
		gridset, err := c.GetGridset(gridsetRef.Name)
		if err != nil {
			return gridsets, err
		}

		gridsets = append(gridsets, gridset)
	}

	return
}

// GetGridset return a single Gridset based on its name
func (c *Client) GetGridset(name string) (gridset *Gridset, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/gridsets/%s", name), nil)
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

	var data Gridset
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return gridset, err
	}

	gridset = &data

	return
}

// CreateGridset creates a Gridset
func (c *Client) CreateGridset(gridsetName string, gridset *Gridset) (err error) {
	payload, _ := xml.Marshal(&gridset)
	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/gridsets/%s", gridsetName), bytes.NewBuffer(payload))
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

// UpdateGridset updates a gridset
func (c *Client) UpdateGridset(gridsetName string, gridset *Gridset) (err error) {
	payload, _ := xml.Marshal(&gridset)

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/gridsets/%s", gridsetName), bytes.NewBuffer(payload))
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
func (c *Client) DeleteGridset(gridsetName string) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/gridsets/%s", gridsetName), nil)
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
