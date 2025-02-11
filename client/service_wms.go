package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type ServiceWmsMetadata struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",innerxml"`
}

type ServiceWmsKeywords struct {
	Keywords []string `xml:"string"`
}

type ServiceVersion struct {
	XMLName xml.Name `xml:"org.geotools.util.Version"`
	Version string   `xml:"version"`
}

// UrlChecks is a list of url check
type ServiceWms struct {
	XMLName                                 xml.Name              `xml:"wms"`
	Name                                    string                `xml:"name"`
	IsEnabled                               bool                  `xml:"enabled"`
	Title                                   string                `xml:"title,omitempty"`
	Maintainer                              string                `xml:"maintainer,omitempty"`
	Abstract                                string                `xml:"abstrct,omitempty"`
	AccessConstraints                       string                `xml:"accessConstraints,omitempty"`
	OnlineResource                          string                `xml:"onlineResource,omitempty"`
	IsVerbose                               bool                  `xml:"verbose"`
	Watermark                               Watermark             `xml:"watermark"`
	Interpolation                           string                `xml:"interpolation"`
	IsCiteCompliant                         bool                  `xml:"citeCompliant"`
	MaximumBuffer                           int                   `xml:"maxBuffer"`
	IsDynamicStylingDisabled                bool                  `xml:"dynamicStylingDisabled"`
	Metadata                                []*ServiceWmsMetadata `xml:"metadata>entry,omitempty"`
	Keywords                                ServiceWmsKeywords    `xml:"keywords,omitempty"`
	IsGetFeatureInfoMimeTypeCheckingEnabled bool                  `xml:"getFeatureInfoMimeTypeCheckingEnabled"`
	MaximumRequestMemory                    int                   `xml:"maxRequestMemory"`
	Fees                                    string                `xml:"fees,omitempty"`
	MaximumRenderingErrors                  int                   `xml:"maxRenderingErrors"`
	MaximumRenderingTime                    int                   `xml:"maxRenderingTime"`
	Workspace                               *WorkspaceRef         `xml:"workspace,omitempty"`
	SupportedVersions                       []*ServiceVersion     `xml:"versions,omitempty"`
	SchemaBaseURL                           string                `xml:"schemaBaseURL"`
	UseBBOXForEachCRS                       bool                  `xml:"bboxForEachCRS"`                // Undocumented in Swagger
	IsGetMapMimeTypeCheckingEnabled         bool                  `xml:"getMapMimeTypeCheckingEnabled"` // Undocumented in Swagger
	IsFeaturesReprojectionDisabled          bool                  `xml:"featuresReprojectionDisabled"`  // Undocumented in Swagger
	MaximumRequestedDimensionValues         int                   `xml:"maxRequestedDimensionValues"`   // Undocumented in Swagger
	CacheConfiguration                      CacheConfiguration    `xml:"cacheConfiguration"`            // Undocumented in Swagger
	RemoteStyleMaxRequestTime               int                   `xml:"remoteStyleMaxRequestTime"`     // Undocumented in Swagger
	RemoteStyleTimeout                      int                   `xml:"remoteStyleTimeout"`            // Undocumented in Swagger
	IsDefaultGroupStyleEnabled              bool                  `xml:"defaultGroupStyleEnabled"`      // Undocumented in Swagger
	IsTransformFeatureInfoDisabled          bool                  `xml:"transformFeatureInfoDisabled"`  // Undocumented in Swagger
	IsAutoEscapeTemplateValuesEnabled       bool                  `xml:"autoEscapeTemplateValues"`      // Undocumented in Swagger
	RootLayerTitle                          string                `xml:"rootLayerTitle,omitempty"`      // Undocumented in Swagger
}

type Watermark struct {
	IsEnabled    bool   `xml:"enabled"`
	Position     string `xml:"position"`
	Transparency int    `xml:"transparency"`
}

type CacheConfiguration struct {
	IsEnabled    bool `xml:"enabled"`
	MaxEntrySize int  `xml:"maxEntrySize"`
	MaxEntries   int  `xml:"maxEntries"`
}

// GetServiceWMS return WMS Service Configuration
func (c *Client) GetServiceWMS(workspace string) (serviceWms *ServiceWms, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = "/services/wms/settings"
	} else {
		endpoint = fmt.Sprintf("/services/wms/workspaces/%s/settings", workspace)
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

	var data ServiceWms
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return serviceWms, err
	}

	serviceWms = &data

	return
}

// UpdateServiceWMS update the configuration of a WMS service
func (c *Client) UpdateServiceWMS(workspace string, serviceWms *ServiceWms) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = "/services/wms/settings"
	} else {
		endpoint = fmt.Sprintf("/services/wms/workspaces/%s/settings", workspace)
	}

	serviceWms.XMLName = xml.Name{
		Local: "wms",
	}
	serviceWms.Name = "WMS"
	payload, _ := xml.Marshal(serviceWms)

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
	case 201:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}

func (c *Client) DeleteWorkspaceServiceWms(workspace string) (err error) {
	var endpoint string

	if workspace == "" {
		err = fmt.Errorf("Workspace MUST be defined")
		return
	} else {
		endpoint = fmt.Sprintf("/services/wms/workspaces/%s/settings", workspace)
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
