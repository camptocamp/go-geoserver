package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type GwcQuota struct {
	Value int    `xml:"value"`
	Units string `xml:"units"`
}

type GwcLayerQuota struct {
	Layer                string   `xml:"layer"`
	ExpirationPolicyName string   `xml:"expirationPolicyName"`
	Quota                GwcQuota `xml:"quota"`
}

type GwcQuotaConfiguration struct {
	XMLName                    xml.Name         `xml:"gwcQuotaConfiguration"`
	Enabled                    bool             `xml:"enabled"`
	CacheCleanUpFrequency      int              `xml:"cacheCleanUpFrequency"`
	CacheCleanUpUnits          string           `xml:"cacheCleanUpUnits"`
	MaxConcurrentCleanUps      int              `xml:"maxConcurrentCleanUps"`
	GlobalExpirationPolicyName string           `xml:"globalExpirationPolicyName"`
	GlobalQuota                GwcQuota         `xml:"globalQuota"`
	LayersQuotas               []*GwcLayerQuota `xml:"layerQuotas>LayerQuota"`
}

// GetGwcQuotaConfiguration return the GeoWebCache Quota Configuration of the instance
func (c *Client) GetGwcQuotaConfiguration() (gwcQuotCfg *GwcQuotaConfiguration, err error) {
	statusCode, body, err := c.doRequest("GET", "/diskquota.xml", nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 200:
		break
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}

	var data GwcQuotaConfiguration
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return gwcQuotCfg, err
	}

	gwcQuotCfg = &data

	return
}

// UpdateGwcQuotaConfiguration updates the disk quota configuration of GeoWebCache
func (c *Client) UpdateGwcQuotaConfiguration(gwcQuotCfg *GwcQuotaConfiguration) (err error) {
	payload, _ := xml.Marshal(&gwcQuotCfg)

	statusCode, body, err := c.doRequest("PUT", "/diskquota.xml", bytes.NewBuffer(payload))
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
