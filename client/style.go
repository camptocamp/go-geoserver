package client

import (
	"encoding/xml"
	"fmt"
)

// Style is a Geoserver object
type Styles struct {
	XMLName xml.Name `xml:"styles"`
	List    []*Style `xml:"style"`
}

type Style struct {
	XMLName   xml.Name        `xml:"style"`
	Name      string          `xml:"name"`
	Workspace WorkspaceRef    `xml:"workspace"`
	Format    string          `xml:"format"`
	Version   LanguageVersion `xml:"languageVersion"`
	FileName  string          `xml:"filename"`
}

type WorkspaceRef struct {
	Name string `xml:"name"`
}

type LanguageVersion struct {
	Version string `xml:"version"`
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

	// Try to retrieve the style file based on the style format
	switch style.Format {
	case "sld":
		endpoint += ".sld"
		break
	case "css":
		endpoint += ".css"
		break
	case "yaml":
		endpoint += ".yaml"
		break
	case "json":
		endpoint += ".json"
		break
	}

	return
}

// Get the style definition of a given format
func (c *Client) GetStyleFile(workspace, name string, styleFormat string, formatVersion string) (styleFile string, err error) {
	var endpoint string

	if workspace == "" {
		endpoint = fmt.Sprintf("/styles/%s", name)
	} else {
		endpoint = fmt.Sprintf("/workspaces/%s/styles/%s", workspace, name)
	}

	// Try to retrieve the style file based on the style format
	var contentType string
	switch styleFormat {
	case "sld":
		if formatVersion == "1.0.0" {
			contentType = "application/vnd.ogc.sld+xml"
		} else {
			contentType = "application/vnd.ogc.se+xml "
		}
		break
	case "css":
		contentType = "application/vnd.geoserver.geocss+css"
		break
	case "yaml":
		contentType = "application/vnd.geoserver.ysld+yaml"
		break
	case "json":
		contentType = "application/vnd.geoserver.mbstyle+json "
		break
	}

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
