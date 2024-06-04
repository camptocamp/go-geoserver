package client

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLayersNoWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/layers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layers>
			<layer>
			<name>signalement</name>
			</layer>
		</layers>
		`))
	})
	mux.HandleFunc("/layers/signalement", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layer>
		<name>signalement</name>
		<type>VECTOR</type>
		<defaultStyle>
		  <name>generic</name>
		</defaultStyle>
		<resource class="featureType">
		  <name>nexsis:signalement</name>
		</resource>
		<attribution>
		  <logoWidth>0</logoWidth>
		  <logoHeight>0</logoHeight>
		</attribution>
	  </layer>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*Layer{
		{
			XMLName:             xml.Name{Space: "", Local: "layer"},
			Name:                "signalement",
			Type:                "VECTOR",
			DefaultStyle:        "generic",
			LayerResource:       Resource{Class: "featureType", Name: "nexsis:signalement"},
			ProviderAttribution: Attribution{LogoWidth: 0, LogoHeight: 0},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layers, err := cli.GetLayers("")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layers)
}

func TestGetLayersWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/layers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layers>
			<layer>
			<name>signalement</name>
			</layer>
		</layers>
		`))
	})
	mux.HandleFunc("/workspaces/foo/layers/signalement", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layer>
		<name>signalement</name>
		<type>VECTOR</type>
		<defaultStyle>
		  <name>generic</name>
		</defaultStyle>
		<resource class="featureType">
		  <name>nexsis:signalement</name>
		</resource>
		<attribution>
		  <logoWidth>0</logoWidth>
		  <logoHeight>0</logoHeight>
		</attribution>
	  </layer>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*Layer{
		{
			XMLName:             xml.Name{Space: "", Local: "layer"},
			Name:                "signalement",
			Type:                "VECTOR",
			DefaultStyle:        "generic",
			LayerResource:       Resource{Class: "featureType", Name: "nexsis:signalement"},
			ProviderAttribution: Attribution{LogoWidth: 0, LogoHeight: 0},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layers, err := cli.GetLayers("foo")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layers)
}

func TestGetLayerNoWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/layers/signalement", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layer>
		<name>signalement</name>
		<type>VECTOR</type>
		<defaultStyle>
		  <name>generic</name>
		</defaultStyle>
		<resource class="featureType">
		  <name>nexsis:signalement</name>
		</resource>
		<attribution>
		  <logoWidth>0</logoWidth>
		  <logoHeight>0</logoHeight>
		</attribution>
	  </layer>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &Layer{
		XMLName:             xml.Name{Space: "", Local: "layer"},
		Name:                "signalement",
		Type:                "VECTOR",
		DefaultStyle:        "generic",
		LayerResource:       Resource{Class: "featureType", Name: "nexsis:signalement"},
		ProviderAttribution: Attribution{LogoWidth: 0, LogoHeight: 0},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layer, err := cli.GetLayer("", "signalement")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layer)
}

func TestGetLayerWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/layers/signalement", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layer>
		<name>signalement</name>
		<type>VECTOR</type>
		<defaultStyle>
		  <name>generic</name>
		</defaultStyle>
		<resource class="featureType">
		  <name>nexsis:signalement</name>
		</resource>
		<attribution>
		  <logoWidth>0</logoWidth>
		  <logoHeight>0</logoHeight>
		</attribution>
	  </layer>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &Layer{
		XMLName:             xml.Name{Space: "", Local: "layer"},
		Name:                "signalement",
		Type:                "VECTOR",
		DefaultStyle:        "generic",
		LayerResource:       Resource{Class: "featureType", Name: "nexsis:signalement"},
		ProviderAttribution: Attribution{LogoWidth: 0, LogoHeight: 0},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layer, err := cli.GetLayer("foo", "signalement")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layer)
}

func TestGetLayerUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layers/toto")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layer, err := cli.GetLayer("", "toto")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, layer)
}

func TestGetLayerNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layers/toto")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layer, err := cli.GetLayer("", "toto")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, layer)
}

func TestGetLayerUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layers/toto")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layer, err := cli.GetLayer("", "toto")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, layer)
}

func TestDeleteLayerSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/layers/toto")
		keys, ok := r.URL.Query()["recurse"]
		assert.True(t, ok)
		assert.Equal(t, keys[0], "true")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteLayer("", "toto", true)

	assert.Nil(t, err)
}
