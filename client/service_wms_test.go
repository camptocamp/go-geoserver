package client

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWmsserviceNoWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/services/wms/settings", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
<wms>
  <enabled>true</enabled>
  <name>WMS</name>
  <title>nexsis</title>
  <abstrct>Abstract de chez abstract</abstrct>
  <citeCompliant>false</citeCompliant>
  <schemaBaseURL>http://schemas.opengis.net</schemaBaseURL>
  <verbose>false</verbose>
  <metadata>
    <entry key="disableDatelineWrappingHeuristic">false</entry>
    <entry key="kmlSuperoverlayMode">auto</entry>
  </metadata>
  <bboxForEachCRS>false</bboxForEachCRS>
  <watermark>
    <enabled>false</enabled>
    <position>BOT_RIGHT</position>
    <transparency>100</transparency>
  </watermark>
  <interpolation>Nearest</interpolation>
  <getFeatureInfoMimeTypeCheckingEnabled>false</getFeatureInfoMimeTypeCheckingEnabled>
  <getMapMimeTypeCheckingEnabled>false</getMapMimeTypeCheckingEnabled>
  <dynamicStylingDisabled>false</dynamicStylingDisabled>
  <featuresReprojectionDisabled>false</featuresReprojectionDisabled>
  <maxBuffer>0</maxBuffer>
  <maxRequestMemory>0</maxRequestMemory>
  <maxRenderingTime>0</maxRenderingTime>
  <maxRenderingErrors>0</maxRenderingErrors>
  <rootLayerTitle>ansc</rootLayerTitle>
  <maxRequestedDimensionValues>100</maxRequestedDimensionValues>
  <cacheConfiguration>
    <enabled>false</enabled>
    <maxEntries>1000</maxEntries>
    <maxEntrySize>51200</maxEntrySize>
  </cacheConfiguration>
  <remoteStyleMaxRequestTime>60000</remoteStyleMaxRequestTime>
  <remoteStyleTimeout>30000</remoteStyleTimeout>
  <defaultGroupStyleEnabled>true</defaultGroupStyleEnabled>
  <transformFeatureInfoDisabled>false</transformFeatureInfoDisabled>
  <autoEscapeTemplateValues>false</autoEscapeTemplateValues>
</wms>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &ServiceWms{
		XMLName: xml.Name{
			Space: "",
			Local: "wms",
		},
		Name:            "WMS",
		Title:           "nexsis",
		IsEnabled:       true,
		Abstract:        "Abstract de chez abstract",
		IsCiteCompliant: false,
		SchemaBaseURL:   "http://schemas.opengis.net",
		IsVerbose:       false,
		Metadata: []*ServiceWmsMetadata{
			{
				Key:   "disableDatelineWrappingHeuristic",
				Value: "false",
			},
			{
				Key:   "kmlSuperoverlayMode",
				Value: "auto",
			},
		},
		UseBBOXForEachCRS:                       false,
		Watermark:                               Watermark{IsEnabled: false, Position: "BOT_RIGHT", Transparency: 100},
		Interpolation:                           "Nearest",
		IsGetFeatureInfoMimeTypeCheckingEnabled: false,
		IsGetMapMimeTypeCheckingEnabled:         false,
		IsDynamicStylingDisabled:                false,
		IsFeaturesReprojectionDisabled:          false,
		MaximumBuffer:                           0,
		MaximumRequestMemory:                    0,
		MaximumRenderingTime:                    0,
		MaximumRenderingErrors:                  0,
		RootLayerTitle:                          "ansc",
		MaximumRequestedDimensionValues:         100,
		CacheConfiguration:                      CacheConfiguration{IsEnabled: false, MaxEntries: 1000, MaxEntrySize: 51200},
		RemoteStyleMaxRequestTime:               60000,
		RemoteStyleTimeout:                      30000,
		IsDefaultGroupStyleEnabled:              true,
		IsTransformFeatureInfoDisabled:          false,
		IsAutoEscapeTemplateValuesEnabled:       false,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	serviceWms, err := cli.GetServiceWMS("")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, serviceWms)
}

func TestGetWmsServiceWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/services/wms/workspaces/foo/settings", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
<wms>
  <enabled>true</enabled>
  <name>WMS</name>
  <workspace>
    <name>foo</name>
  </workspace>
  <title>nexsis</title>
  <abstrct>Abstract de chez abstract</abstrct>
  <citeCompliant>false</citeCompliant>
  <schemaBaseURL>http://schemas.opengis.net</schemaBaseURL>
  <verbose>false</verbose>
  <metadata>
    <entry key="disableDatelineWrappingHeuristic">false</entry>
    <entry key="kmlSuperoverlayMode">auto</entry>
  </metadata>
  <bboxForEachCRS>false</bboxForEachCRS>
  <watermark>
    <enabled>false</enabled>
    <position>BOT_RIGHT</position>
    <transparency>100</transparency>
  </watermark>
  <interpolation>Nearest</interpolation>
  <getFeatureInfoMimeTypeCheckingEnabled>false</getFeatureInfoMimeTypeCheckingEnabled>
  <getMapMimeTypeCheckingEnabled>false</getMapMimeTypeCheckingEnabled>
  <dynamicStylingDisabled>false</dynamicStylingDisabled>
  <featuresReprojectionDisabled>false</featuresReprojectionDisabled>
  <maxBuffer>0</maxBuffer>
  <maxRequestMemory>0</maxRequestMemory>
  <maxRenderingTime>0</maxRenderingTime>
  <maxRenderingErrors>0</maxRenderingErrors>
  <rootLayerTitle>ansc</rootLayerTitle>
  <maxRequestedDimensionValues>100</maxRequestedDimensionValues>
  <cacheConfiguration>
    <enabled>false</enabled>
    <maxEntries>1000</maxEntries>
    <maxEntrySize>51200</maxEntrySize>
  </cacheConfiguration>
  <remoteStyleMaxRequestTime>60000</remoteStyleMaxRequestTime>
  <remoteStyleTimeout>30000</remoteStyleTimeout>
  <defaultGroupStyleEnabled>true</defaultGroupStyleEnabled>
  <transformFeatureInfoDisabled>false</transformFeatureInfoDisabled>
  <autoEscapeTemplateValues>false</autoEscapeTemplateValues>
</wms>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &ServiceWms{
		XMLName: xml.Name{
			Space: "",
			Local: "wms",
		},
		Name:            "WMS",
		Workspace:       &WorkspaceRef{Name: "foo"},
		Title:           "nexsis",
		IsEnabled:       true,
		Abstract:        "Abstract de chez abstract",
		IsCiteCompliant: false,
		SchemaBaseURL:   "http://schemas.opengis.net",
		IsVerbose:       false,
		Metadata: []*ServiceWmsMetadata{
			{
				Key:   "disableDatelineWrappingHeuristic",
				Value: "false",
			},
			{
				Key:   "kmlSuperoverlayMode",
				Value: "auto",
			},
		},
		UseBBOXForEachCRS:                       false,
		Watermark:                               Watermark{IsEnabled: false, Position: "BOT_RIGHT", Transparency: 100},
		Interpolation:                           "Nearest",
		IsGetFeatureInfoMimeTypeCheckingEnabled: false,
		IsGetMapMimeTypeCheckingEnabled:         false,
		IsDynamicStylingDisabled:                false,
		IsFeaturesReprojectionDisabled:          false,
		MaximumBuffer:                           0,
		MaximumRequestMemory:                    0,
		MaximumRenderingTime:                    0,
		MaximumRenderingErrors:                  0,
		RootLayerTitle:                          "ansc",
		MaximumRequestedDimensionValues:         100,
		CacheConfiguration:                      CacheConfiguration{IsEnabled: false, MaxEntries: 1000, MaxEntrySize: 51200},
		RemoteStyleMaxRequestTime:               60000,
		RemoteStyleTimeout:                      30000,
		IsDefaultGroupStyleEnabled:              true,
		IsTransformFeatureInfoDisabled:          false,
		IsAutoEscapeTemplateValuesEnabled:       false,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	serviceWms, err := cli.GetServiceWMS("foo")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, serviceWms)
}

func TestGetServiceWmsUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/services/wms/settings")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	serviceWms, err := cli.GetServiceWMS("")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, serviceWms)
}

func TestGetServiceWmsNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/services/wms/workspaces/foo/settings")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	serviceWms, err := cli.GetServiceWMS("foo")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, serviceWms)
}

func TestGetServiceWmsUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/services/wms/settings")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	serviceWms, err := cli.GetServiceWMS("")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, serviceWms)
}

func TestDeleteServiceWmsSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/services/wms/workspaces/foo/settings")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteWorkspaceServiceWms("foo")

	assert.Nil(t, err)
}
