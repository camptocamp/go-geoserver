package client

import (
	"fmt"
	"strings"
)

// GetResource returns a file stored in the resource store
func (c *Client) GetResource(pathToResource string, resourceExtension string) (resourceContent string, err error) {
	if resourceExtension == "" {
		err = fmt.Errorf("retrieving content of resource is only possible for files")
		return
	}

	var endpoint string = fmt.Sprintf("/resource/%s.%s", pathToResource, resourceExtension)

	statusCode, resourceContent, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 404:
		err = fmt.Errorf("resource not found")
		return
	case 200:
		break
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, resourceContent)
		return
	}

	return resourceContent, err
}

// CreateResource creates a resource on GeoServer
func (c *Client) CreateResource(pathToResource string, resourceExtension string, resourceContent string) (err error) {
	if resourceExtension == "" {
		err = fmt.Errorf("creation of resource is only possible for files")
		return
	}

	if resourceContent == "" {
		err = fmt.Errorf("resource content must be defined")
		return
	}

	var endpoint string = fmt.Sprintf("/resource/%s.%s", pathToResource, resourceExtension)
	statusCode, body, err := c.doRequest("PUT", endpoint, strings.NewReader(resourceContent))

	if err != nil {
		return
	}

	switch statusCode {
	case 404:
		err = fmt.Errorf("source path that doesn’t exist")
		return
	case 405:
		err = fmt.Errorf("PUT to directory or copy where source path is directory")
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

// UpdateResource updates an existing resource
func (c *Client) UpdateResource(pathToResource string, resourceExtension string, resourceContent string) (err error) {
	if resourceExtension == "" {
		err = fmt.Errorf("creation of resource is only possible for files")
		return
	}

	if resourceContent == "" {
		err = fmt.Errorf("resource content must be defined")
		return
	}

	var endpoint string = fmt.Sprintf("/resource/%s.%s", pathToResource, resourceExtension)
	statusCode, body, err := c.doRequest("PUT", endpoint, strings.NewReader(resourceContent))

	if err != nil {
		return
	}

	switch statusCode {
	case 404:
		err = fmt.Errorf("source path that doesn’t exist")
		return
	case 405:
		err = fmt.Errorf("PUT to directory or copy where source path is directory")
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

// DeleteResource deletes a resource from GeoServer
func (c *Client) DeleteResource(resource string) (err error) {
	var endpoint string = fmt.Sprintf("/resource/%s", resource)

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
