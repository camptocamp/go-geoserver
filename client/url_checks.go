package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// UrlChecks is a list of url check
type UrlChecks struct {
	XMLName xml.Name    `xml:"urlChecks"`
	List    []*UrlCheck `xml:"urlCheck"`
}

// UrlCheck is a reference to a url check definition
type UrlCheck struct {
	XMLName xml.Name `xml:"urlCheck"`
	Name    string   `xml:"name"`
}

// RegexUrlCheck defines a Url Check based on a regular expression
type RegexUrlCheck struct {
	XMLName     xml.Name `xml:"regexUrlCheck"`
	Name        string   `xml:"name"`
	Description string   `xml:"description,omitempty"`
	Regex       string   `xml:"regex"`
	IsEnabled   bool     `xml:"enabled,omitempty"`
}

// GetUrlChecks returns all the URL checks
func (c *Client) GetUrlChecks() (urlChecks []*RegexUrlCheck, err error) {
	var endpoint string = "/urlchecks"

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

	var data UrlChecks

	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return urlChecks, err
	}

	for _, urlCheckRef := range data.List {
		urlCheck, err := c.GetRegExUrlCheck(urlCheckRef.Name)
		if err != nil {
			return urlChecks, err
		}

		urlChecks = append(urlChecks, urlCheck)
	}

	return
}

// GetRegExUrlCheck returns the definition of a RegEx based URL check
func (c *Client) GetRegExUrlCheck(urlCheckName string) (regExUrlCheck *RegexUrlCheck, err error) {

	var endpoint string = fmt.Sprintf("/urlchecks/%s", urlCheckName)

	statusCode, body, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 404:
		err = fmt.Errorf("URL checks not found")
		return
	case 200:
		break
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}

	var data RegexUrlCheck
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return regExUrlCheck, err
	}

	regExUrlCheck = &data

	return
}

// CreateRegExUrlCheck creates a new URl checks on GeoServer
func (c *Client) CreateRegExUrlCheck(checkName string, checkDefinition *RegexUrlCheck) (err error) {
	var endpoint string = "/urlchecks"

	checkDefinition.XMLName = xml.Name{
		Local: "regexUrlCheck",
	}
	payload, err := xml.Marshal(checkDefinition)
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

// UpdateRegExUrlCheck updates an existing URL Check based on a regexp
func (c *Client) UpdateRegExUrlCheck(checkName string, checkDefinition *RegexUrlCheck) (err error) {
	var endpoint string = fmt.Sprintf("/urlchecks/%s", checkName)

	checkDefinition.XMLName = xml.Name{
		Local: "featureType",
	}
	payload, _ := xml.Marshal(checkDefinition)

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

// DeleteUrlCheck deletes a url check from GeoServer
func (c *Client) DeleteUrlCheck(checkName string) (err error) {
	var endpoint string = fmt.Sprintf("/urlchecks/%s", checkName)

	statusCode, body, err := c.doRequest("DELETE", endpoint, nil)
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
