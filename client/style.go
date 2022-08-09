package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

// Styles is a list of Style
type Styles struct {
	XMLName xml.Name `xml:"styles"`
	List    []*Style `xml:"style"`
}

// Style is GeoServer Resource
type Style struct {
	XMLName   xml.Name         `xml:"style"`
	Name      string           `xml:"name"`
	Workspace *WorkspaceRef    `xml:"workspace,omitempty"`
	Format    string           `xml:"format,omitempty"`
	Version   *LanguageVersion `xml:"languageVersion,omitempty"`
	FileName  string           `xml:"filename"`
}

// WorkspaceRef is a reference to a GeoServer workspace
type WorkspaceRef struct {
	Name string `xml:"name,omitempty"`
}

// LanguageVersion is the version of the language used to described the style
type LanguageVersion struct {
	Version string `xml:"version,omitempty"`
}

// GetHTTPContentTypeFor computes the content type of a http request for the required format and version
func (c *Client) GetHTTPContentTypeFor(format string, version string) (contentType string) {
	switch format {
	case "sld":
		if version == "1.0.0" {
			return "application/vnd.ogc.sld+xml"
		}
		return "application/vnd.ogc.se+xml "
	case "css":
		return "application/vnd.geoserver.geocss+css"
	case "yaml":
		return "application/vnd.geoserver.ysld+yaml"
	case "json":
		return "application/vnd.geoserver.mbstyle+json "
	default:
		return "application/vnd.ogc.sld+xml"
	}
}

// GetStyles returns all the styles
func (c *Client) GetStyles(workspace string) (styles []*Style, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles")
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles", workspace)
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

	var data Styles

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return styles, err
	}

	for _, styleRef := range data.List {
		style, err := c.GetStyle(workspace, styleRef.Name)
		if err != nil {
			return styles, err
		}

		styles = append(styles, style)
	}

	return
}

// GetStyle return a single style based on its name
func (c *Client) GetStyle(workspace, name string) (style *Style, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles/%s", name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles/%s", workspace, name)
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

	var data Style
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return style, err
	}

	style = &data

	return
}

// GetStyleFile retrieves the style definition of a given format
func (c *Client) GetStyleFile(workspace, name string, styleFormat string, formatVersion string) (styleFile string, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles/%s", name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles/%s", workspace, name)
	}

	// Try to retrieve the style file based on the style format
	contentType := c.GetHTTPContentTypeFor(styleFormat, formatVersion)

	statusCode, styleFile, err := c.doTypedRequest("GET", endpoint, nil, contentType)
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
		err = fmt.Errorf("unknown error: %d - %s", statusCode, styleFile)
		return
	}

	return styleFile, err
}

// CreateStyle creates a style
func (c *Client) CreateStyle(workspace string, style *Style) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles")
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles", workspace)
	}

	style.XMLName = xml.Name{
		Local: "style",
	}
	payload, err := xml.Marshal(style)
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

// UpdateStyle creates a style
func (c *Client) UpdateStyle(workspace string, style *Style, styleDefinition string) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles")
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles", workspace)
	}

	contentType := c.GetHTTPContentTypeFor(style.Format, style.Version.Version)

	statusCode, body, err := c.doFullyTypedRequest("POST", endpoint, strings.NewReader(styleDefinition), contentType, "")
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
		err = fmt.Errorf("unknown error: %d - %s - %s", statusCode, body, styleDefinition)
		return
	}
}

// UpdateStyleContent changes the style definition
func (c *Client) UpdateStyleContent(workspace string, style *Style, styleDefinition string) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles/%s", style.Name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles/%s", workspace, style.Name)
	}

	contentType := c.GetHTTPContentTypeFor(style.Format, style.Version.Version)

	statusCode, body, err := c.doFullyTypedRequest("PUT", endpoint, strings.NewReader(styleDefinition), contentType, "")
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 200:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s - %s", statusCode, body, styleDefinition)
		return
	}
}

// DeleteStyle deletes style from GeoServer
func (c *Client) DeleteStyle(workspace string, style string, purge bool, recurse bool) (err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles/%s?purge=%t&recurse=%t", style, purge, recurse)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles/%s?purge=%t&recurse=%t", workspace, style, purge, recurse)
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
