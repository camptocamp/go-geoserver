package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"slices"
)

// LayerRules is a list of ACL rule
type LayerRules struct {
	XMLName xml.Name     `xml:"rules"`
	List    []*LayerRule `xml:"rule"`
}

// LayerRule is a Geoserver object
type LayerRule struct {
	XMLName  xml.Name `xml:"rule"`
	Resource string   `xml:"resource,attr"`
	Rule     string   `xml:",chardata"`
}

// GetLayerRules returns the list of the layer rules
func (c *Client) GetLayerRules() (rules LayerRules, err error) {
	statusCode, body, err := c.doRequest("GET", "/security/acl/layers", nil)
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

	var data LayerRules
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return rules, nil
	}

	return data, err
}

// GetLayerRule return a rule based on its definition
func (c *Client) GetLayerRule(ruleDef string) (rule *LayerRule, err error) {
	rules, rulesErr := c.GetLayerRules()

	if rulesErr != nil {
		return rule, rulesErr
	}

	ruleIdx := slices.IndexFunc(rules.List, func(rule *LayerRule) bool { return rule.Resource == ruleDef })

	if ruleIdx == -1 {
		return rule, fmt.Errorf("Rule not found")
	}

	rule = rules.List[ruleIdx]
	return
}

// CreateRule creates a ACL rule
func (c *Client) CreateLayerRule(rule *LayerRule) (err error) {
	rule.XMLName = xml.Name{
		Local: "rule",
	}

	layerRules := LayerRules{
		XMLName: xml.Name{
			Local: "rules",
		},
		List: []*LayerRule{
			rule,
		},
	}

	payload, err := xml.Marshal(layerRules)
	if err != nil {
		return
	}
	statusCode, body, err := c.doFullyTypedRequest("POST", "/security/acl/layers", bytes.NewBuffer(payload), "application/xml", "")

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
		err = fmt.Errorf("unknown error: %d - %s - %s", statusCode, body, string(payload))
		return
	}
}

// UpdateUser creates a user
func (c *Client) UpdateLayerRule(rule *LayerRule) (err error) {
	rule.XMLName = xml.Name{
		Local: "rule",
	}

	layerRules := LayerRules{
		XMLName: xml.Name{
			Local: "rules",
		},
		List: []*LayerRule{
			rule,
		},
	}

	payload, err := xml.Marshal(layerRules)
	if err != nil {
		return
	}
	statusCode, body, err := c.doFullyTypedRequest("PUT", "/security/acl/layers", bytes.NewBuffer(payload), "application/xml", "")

	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 409:
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

// DeleteLayerRule deletes ACL Rule from GeoServer
func (c *Client) DeleteLayerRule(ruleDefinition string) (err error) {
	var endpoint string = fmt.Sprintf("/security/acl/layers/%s", ruleDefinition)

	statusCode, body, err := c.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 403:
		err = fmt.Errorf("service name is not empty")
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
