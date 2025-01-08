package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQuotaConfigurationSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/diskquota.xml", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<gwcQuotaConfiguration>
			<enabled>false</enabled>
			<cacheCleanUpFrequency>10</cacheCleanUpFrequency>
			<cacheCleanUpUnits>SECONDS</cacheCleanUpUnits>
			<maxConcurrentCleanUps>2</maxConcurrentCleanUps>
			<globalExpirationPolicyName>LFU</globalExpirationPolicyName>
			<globalQuota>
				<value>512</value>
				<units>GiB</units>
			</globalQuota>
			<layerQuotas> <!-- optional -->
				<LayerQuota>
				<layer>topp:states</layer>
				<expirationPolicyName>LRU</expirationPolicyName>
				<quota>
					<value>100</value>
					<units>GiB</units>
				</quota>
				</LayerQuota>
				<!-- Other layers -->
			</layerQuotas>
			</gwcQuotaConfiguration>		
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &GwcQuotaConfiguration{
		XMLName: xml.Name{
			Local: "gwcQuotaConfiguration",
		},
		Enabled:                    false,
		CacheCleanUpFrequency:      10,
		CacheCleanUpUnits:          "SECONDS",
		MaxConcurrentCleanUps:      2,
		GlobalExpirationPolicyName: "LFU",
		GlobalQuota:                GwcQuota{Value: 512, Units: "GiB"},
		LayersQuotas: []*GwcLayerQuota{
			{
				Layer:                "topp:states",
				ExpirationPolicyName: "LRU",
				Quota:                GwcQuota{Value: 100, Units: "GiB"},
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetGwcQuotaConfiguration()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetQuotaConfigurationUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/diskquota.xml")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetGwcQuotaConfiguration()

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, datastore)
}

func TestUpdateQuotaConfigurationSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/diskquota.xml")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *GwcQuotaConfiguration
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &GwcQuotaConfiguration{
			XMLName: xml.Name{
				Local: "gwcQuotaConfiguration",
			},
			Enabled:                    false,
			CacheCleanUpFrequency:      10,
			CacheCleanUpUnits:          "SECONDS",
			MaxConcurrentCleanUps:      2,
			GlobalExpirationPolicyName: "LFU",
			GlobalQuota:                GwcQuota{Value: 512, Units: "GiB"},
			LayersQuotas: []*GwcLayerQuota{
				{
					Layer:                "topp:states",
					ExpirationPolicyName: "LRU",
					Quota:                GwcQuota{Value: 100, Units: "GiB"},
				},
			},
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateGwcQuotaConfiguration(&GwcQuotaConfiguration{
		XMLName: xml.Name{
			Local: "gwcQuotaConfiguration",
		},
		Enabled:                    false,
		CacheCleanUpFrequency:      10,
		CacheCleanUpUnits:          "SECONDS",
		MaxConcurrentCleanUps:      2,
		GlobalExpirationPolicyName: "LFU",
		GlobalQuota:                GwcQuota{Value: 512, Units: "GiB"},
		LayersQuotas: []*GwcLayerQuota{
			{
				Layer:                "topp:states",
				ExpirationPolicyName: "LRU",
				Quota:                GwcQuota{Value: 100, Units: "GiB"},
			},
		},
	})

	assert.Nil(t, err)
}
