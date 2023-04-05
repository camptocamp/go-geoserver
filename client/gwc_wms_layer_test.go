package client

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGwcWmsLayerSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/layers/osm:fdp_normal", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsLayer>
  <blobStoreId>cache-nexsis-osm</blobStoreId>
  <enabled>true</enabled>
  <name>osm:fdp_normal</name>
  <mimeFormats>
    <string>image/png</string>
    <string>image/jpeg</string>
  </mimeFormats>
  <gridSubsets>
    <gridSubset>
      <gridSetName>EPSG:3857</gridSetName>
    </gridSubset>
    <gridSubset>
      <gridSetName>EPSG:4326</gridSetName>
    </gridSubset>
    <gridSubset>
      <gridSetName>EPSG:900913</gridSetName>
    </gridSubset>
  </gridSubsets>
  <metaWidthHeight>
    <int>4</int>
    <int>4</int>
  </metaWidthHeight>
  <expireCache>0</expireCache>
  <expireCacheList>
    <expirationRule minZoom="0" expiration="0"/>
  </expireCacheList>
  <expireClients>0</expireClients>
  <expireClientsList>
    <expirationRule minZoom="0" expiration="0"/>
  </expireClientsList>
  <backendTimeout>120</backendTimeout>
  <cacheBypassAllowed>false</cacheBypassAllowed>
  <parameterFilters/>
  <wmsUrl>
    <string>https://master.dev.scw.ansc.fr/geoserver/ows?service=WMS</string>
  </wmsUrl>
  <wmsLayers>osm:fdp_normal</wmsLayers>
  <gutter>0</gutter>
  <concurrency>32</concurrency>
</wmsLayer>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &GwcWmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:        "osm:fdp_normal",
		Enabled:     true,
		BlobStoreId: "cache-nexsis-osm",
		MimeFormats: MimeFormats{Formats: []string{"image/png", "image/jpeg"}},
		GridSubsets: []*GridSubset{
			{
				Name: "EPSG:3857",
			},
			{
				Name: "EPSG:4326",
			},
			{
				Name: "EPSG:900913",
			},
		},
		MetaWidthHeight:      MetaTileDim{Width: 4, Height: 4},
		ExpireCacheDuration:  0,
		ExpireClientDuration: 0,
		GutterSize:           0,
		BackendTimeout:       120,
		CacheBypassAllowed:   false,
		WmsUrl:               "https://master.dev.scw.ansc.fr/geoserver/ows?service=WMS",
		WmsLayer:             "osm:fdp_normal",
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetGwcWMSLayer("osm:fdp_normal")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetGwcWmsLayerUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layers/sf")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetGwcWMSLayer("sf")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, datastore)
}

func TestGetGwcWmsLayerNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layers/sf")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetGwcWMSLayer("sf")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, datastore)
}

func TestGetGwcWmsLayerUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layers/sf")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetGwcWMSLayer("sf")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, datastore)
}

func TestCreateGwcWmsLayerSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/layers/osm:fdp_normal")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *GwcWmsLayer
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &GwcWmsLayer{
			XMLName: xml.Name{
				Local: "wmsLayer",
			},
			Name:        "osm:fdp_normal",
			Enabled:     true,
			BlobStoreId: "cache-nexsis-osm",
			MimeFormats: MimeFormats{Formats: []string{"image/png", "image/jpeg"}},
			GridSubsets: []*GridSubset{
				{
					Name: "EPSG:3857",
				},
				{
					Name: "EPSG:4326",
				},
				{
					Name: "EPSG:900913",
				},
			},
			MetaWidthHeight:      MetaTileDim{Width: 4, Height: 4},
			ExpireCacheDuration:  0,
			ExpireClientDuration: 0,
			GutterSize:           0,
			BackendTimeout:       120,
			CacheBypassAllowed:   false,
			WmsUrl:               "https://master.dev.scw.ansc.fr/geoserver/ows?service=WMS",
			WmsLayer:             "osm:fdp_normal",
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateGwcWmsLayer("osm:fdp_normal", &GwcWmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:        "osm:fdp_normal",
		Enabled:     true,
		BlobStoreId: "cache-nexsis-osm",
		MimeFormats: MimeFormats{Formats: []string{"image/png", "image/jpeg"}},
		GridSubsets: []*GridSubset{
			{
				Name: "EPSG:3857",
			},
			{
				Name: "EPSG:4326",
			},
			{
				Name: "EPSG:900913",
			},
		},
		MetaWidthHeight:      MetaTileDim{Width: 4, Height: 4},
		ExpireCacheDuration:  0,
		ExpireClientDuration: 0,
		GutterSize:           0,
		BackendTimeout:       120,
		CacheBypassAllowed:   false,
		WmsUrl:               "https://master.dev.scw.ansc.fr/geoserver/ows?service=WMS",
		WmsLayer:             "osm:fdp_normal",
	})

	assert.Nil(t, err)
}

func TestUpdateGwcWmsLayerSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/layers/osm:fdp_normal")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *GwcWmsLayer
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &GwcWmsLayer{
			XMLName: xml.Name{
				Local: "wmsLayer",
			},
			Name:        "osm:fdp_normal",
			Enabled:     true,
			BlobStoreId: "cache-nexsis-osm",
			MimeFormats: MimeFormats{Formats: []string{"image/png", "image/jpeg"}},
			GridSubsets: []*GridSubset{
				{
					Name: "EPSG:3857",
				},
				{
					Name: "EPSG:4326",
				},
				{
					Name: "EPSG:900913",
				},
			},
			MetaWidthHeight:      MetaTileDim{Width: 4, Height: 4},
			ExpireCacheDuration:  0,
			ExpireClientDuration: 0,
			GutterSize:           0,
			BackendTimeout:       120,
			CacheBypassAllowed:   false,
			WmsUrl:               "https://master.dev.scw.ansc.fr/geoserver/ows?service=WMS",
			WmsLayer:             "osm:fdp_normal",
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateGwcWmsLayer("osm:fdp_normal", &GwcWmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:        "osm:fdp_normal",
		Enabled:     true,
		BlobStoreId: "cache-nexsis-osm",
		MimeFormats: MimeFormats{Formats: []string{"image/png", "image/jpeg"}},
		GridSubsets: []*GridSubset{
			{
				Name: "EPSG:3857",
			},
			{
				Name: "EPSG:4326",
			},
			{
				Name: "EPSG:900913",
			},
		},
		MetaWidthHeight:      MetaTileDim{Width: 4, Height: 4},
		ExpireCacheDuration:  0,
		ExpireClientDuration: 0,
		GutterSize:           0,
		BackendTimeout:       120,
		CacheBypassAllowed:   false,
		WmsUrl:               "https://master.dev.scw.ansc.fr/geoserver/ows?service=WMS",
		WmsLayer:             "osm:fdp_normal",
	})

	assert.Nil(t, err)
}

func TestDeleteGwcWmsLayerSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/layers/osm:fdp_normal")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteGwcWmsLayer("osm:fdp_normal")

	assert.Nil(t, err)
}
