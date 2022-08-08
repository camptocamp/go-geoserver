package client

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Client contains information to connect to a Geoserver instance
type Client struct {
	URL      string
	Username string
	Password string

	HTTPClient *http.Client
}

// NewClient returns a Client that connect to a Geoserver instance
func NewClient(url, username, password string) (client *Client, err error) {
	client = &Client{
		URL:      url,
		Username: username,
		Password: password,

		HTTPClient: &http.Client{},
	}
	return
}

func (c *Client) doRequest(method, path string, data io.Reader) (statusCode int, body string, err error) {
	return c.doTypedRequest(method, path, data, "application/xml")
}

func (c *Client) doTypedRequest(method, path string, data io.Reader, contentType string) (statusCode int, body string, err error) {
	request, err := http.NewRequest(method, c.URL+path, data)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", contentType)
	if c.Username != "" && c.Password != "" {
		request.SetBasicAuth(c.Username, c.Password)
	}
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return
	}
	statusCode = response.StatusCode

	defer response.Body.Close()
	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	body = string(rawBody)

	return
}

func (c *Client) doFullyTypedRequest(method, path string, data io.Reader, contentType string, acceptType string) (statusCode int, body string, err error) {
	request, err := http.NewRequest(method, c.URL+path, data)
	if err != nil {
		return
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if acceptType != "" {
		request.Header.Set("Accept", acceptType)
	}
	if c.Username != "" && c.Password != "" {
		request.SetBasicAuth(c.Username, c.Password)
	}
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return
	}
	statusCode = response.StatusCode

	defer response.Body.Close()
	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	body = string(rawBody)

	return
}
