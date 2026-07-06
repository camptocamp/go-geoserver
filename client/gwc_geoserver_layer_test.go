package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGwcGsLayerSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/layers/osm:fdp_normal", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<GeoServerLayer>
  <id>LayerGroupInfo-01KW94WAF9DH2DEHSVN4QCE7ZT</id>
  <enabled>true</enabled>
  <name>osm:fdp_normal</name>
  <mimeFormats>
    <string>image/png</string>
    <string>image/jpeg</string>
  </mimeFormats>
  <gridSubsets>
    <gridSubset>
      <gridSetName>EPSG:4326</gridSetName>
    </gridSubset>
    <gridSubset>
      <gridSetName>EPSG:900913</gridSetName>
    </gridSubset>
    <gridSubset>
      <gridSetName>EPSG:3857</gridSetName>
    </gridSubset>
  </gridSubsets>
  <metaWidthHeight>
    <int>4</int>
    <int>4</int>
  </metaWidthHeight>
  <expireCache>0</expireCache>
  <expireClients>0</expireClients>
  <parameterFilters/>
  <gutter>0</gutter>
  <cacheWarningSkips/>
</GeoServerLayer>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &GwcGsLayer{
		XMLName: xml.Name{
			Local: "GeoServerLayer",
		},
		Name:        "osm:fdp_normal",
		Enabled:     true,
		Id:          "LayerGroupInfo-01KW94WAF9DH2DEHSVN4QCE7ZT",
		MimeFormats: MimeFormats{Formats: []string{"image/png", "image/jpeg"}},
		GridSubsets: []*GridSubset{
			{
				Name: "EPSG:4326",
			},
			{
				Name: "EPSG:900913",
			},
			{
				Name: "EPSG:3857",
			},
		},
		MetaTileDimensions:   []int{4, 4},
		ExpireCacheDuration:  0,
		ExpireClientDuration: 0,
		GutterSize:           0,
		CacheBypassAllowed:   false,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetGwcGsLayer("osm:fdp_normal")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetGwcGsLayerUnauthorized(t *testing.T) {
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

	datastore, err := cli.GetGwcGsLayer("sf")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, datastore)
}

func TestGetGwcGsLayerNotFound(t *testing.T) {
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

	datastore, err := cli.GetGwcGsLayer("sf")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, datastore)
}

func TestGetGwcGsLayerUnknownError(t *testing.T) {
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

	datastore, err := cli.GetGwcGsLayer("sf")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, datastore)
}

func TestCreateGwcGsLayerSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/layers/osm:fdp_normal")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *GwcGsLayer
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &GwcGsLayer{
			XMLName: xml.Name{
				Local: "GeoServerLayer",
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
			MetaTileDimensions:   []int{4, 4},
			ExpireCacheDuration:  0,
			ExpireClientDuration: 0,
			GutterSize:           0,
			CacheBypassAllowed:   false,
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateGwcGsLayer("osm:fdp_normal", &GwcGsLayer{
		XMLName: xml.Name{
			Local: "GeoServerLayer",
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
		MetaTileDimensions:   []int{4, 4},
		ExpireCacheDuration:  0,
		ExpireClientDuration: 0,
		GutterSize:           0,
		CacheBypassAllowed:   false,
	})

	assert.Nil(t, err)
}

func TestUpdateGwcGsLayerSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/layers/osm:fdp_normal")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *GwcGsLayer
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &GwcGsLayer{
			XMLName: xml.Name{
				Local: "GeoServerLayer",
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
			MetaTileDimensions:   []int{4, 4},
			ExpireCacheDuration:  0,
			ExpireClientDuration: 0,
			GutterSize:           0,
			CacheBypassAllowed:   false,
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateGwcGsLayer("osm:fdp_normal", &GwcGsLayer{
		XMLName: xml.Name{
			Local: "GeoServerLayer",
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
		MetaTileDimensions:   []int{4, 4},
		ExpireCacheDuration:  0,
		ExpireClientDuration: 0,
		GutterSize:           0,
		CacheBypassAllowed:   false,
	})

	assert.Nil(t, err)
}

func TestDeleteGwcGssLayerSuccess(t *testing.T) {
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

	err := cli.DeleteGwcGsLayer("osm:fdp_normal")

	assert.Nil(t, err)
}
