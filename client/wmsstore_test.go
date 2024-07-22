package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWmsStoresSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/wmsstores", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsStores>
			<wmsStore>
				<name>sf</name>
				<atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/sf/wmsstores/sf.xml" type="application/xml"/>
			</wmsStore>
		</wmsStores>
		`))
	})
	mux.HandleFunc("/workspaces/foo/wmsstores/sf", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsStore>
		<name>sf</name>
		<description>Services de la GeoPlateforme IGN</description>
		<enabled>true</enabled>
		<workspace>
		  <name>foo</name>
		  <atom:link 
			xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="https://master.dev.scw.ansc.fr/geoserver/rest/workspaces/ign.xml" type="application/xml"/>
		  </workspace>
		  <__default>false</__default>
		  <disableOnConnFailure>false</disableOnConnFailure>
		  <capabilitiesURL>https://data.geopf.fr/wms-r/wms?SERVICE=WMS&amp;Version=1.3.0&amp;Request=GetCapabilities</capabilitiesURL>
		  <maxConnections>6</maxConnections>
		  <readTimeout>60</readTimeout>
		  <connectTimeout>30</connectTimeout>
		  <wmslayers>
			<atom:link 
			  xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="https://master.dev.scw.ansc.fr/geoserver/rest/workspaces/ign/wmsstores/GeoPlatforme/wmslayers.xml" type="application/xml"/>
			</wmslayers>
		  </wmsStore>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*WmsStore{
		{
			XMLName: xml.Name{
				Local: "wmsStore",
			},
			Name:    "sf",
			Enabled: true,
			Workspace: &WorkspaceReference{
				Name: "foo",
			},
			Default:                    false,
			Description:                "Services de la GeoPlateforme IGN",
			DisableConnectionOnFailure: false,
			CapabilitiesUrl:            "https://data.geopf.fr/wms-r/wms?SERVICE=WMS&Version=1.3.0&Request=GetCapabilities",
			MaxConnections:             6,
			ReadTimeOut:                60,
			ConnectTimeOut:             30,
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetWmsStores("foo")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetWmsStoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/wmsstores/sf", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsStore>
		<name>sf</name>
		<description>Services de la GeoPlateforme IGN</description>
		<enabled>true</enabled>
		<workspace>
		  <name>foo</name>
		  <atom:link 
			xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="https://master.dev.scw.ansc.fr/geoserver/rest/workspaces/ign.xml" type="application/xml"/>
		  </workspace>
		  <__default>false</__default>
		  <disableOnConnFailure>false</disableOnConnFailure>
		  <capabilitiesURL>https://data.geopf.fr/wms-r/wms?SERVICE=WMS&amp;Version=1.3.0&amp;Request=GetCapabilities</capabilitiesURL>
		  <maxConnections>6</maxConnections>
		  <readTimeout>60</readTimeout>
		  <connectTimeout>30</connectTimeout>
		  <wmslayers>
			<atom:link 
			  xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="https://master.dev.scw.ansc.fr/geoserver/rest/workspaces/ign/wmsstores/GeoPlatforme/wmslayers.xml" type="application/xml"/>
			</wmslayers>
		  </wmsStore>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &WmsStore{
		XMLName: xml.Name{
			Local: "wmsStore",
		},
		Name:    "sf",
		Enabled: true,
		Workspace: &WorkspaceReference{
			Name: "foo",
		},
		Default:                    false,
		Description:                "Services de la GeoPlateforme IGN",
		DisableConnectionOnFailure: false,
		CapabilitiesUrl:            "https://data.geopf.fr/wms-r/wms?SERVICE=WMS&Version=1.3.0&Request=GetCapabilities",
		MaxConnections:             6,
		ReadTimeOut:                60,
		ConnectTimeOut:             30,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsStore, err := cli.GetWmsStore("foo", "sf")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, wmsStore)
}

func TestGetWmsStoreUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmsstores/sf")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsStore, err := cli.GetWmsStore("foo", "sf")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, wmsStore)
}

func TestGetWmsStoreNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmsstores/sf")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsStore, err := cli.GetWmsStore("foo", "sf")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, wmsStore)
}

func TestGetWmsStoreUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmsstores/sf")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsStore, err := cli.GetWmsStore("foo", "sf")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, wmsStore)
}

func TestCreateWmsStoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmsstores")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *WmsStore
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &WmsStore{
			XMLName: xml.Name{
				Local: "wmsStore",
			},
			Name:    "sf",
			Enabled: true,
			Workspace: &WorkspaceReference{
				Name: "foo",
			},
			Default:                    false,
			Description:                "Services de la GeoPlateforme IGN",
			DisableConnectionOnFailure: false,
			CapabilitiesUrl:            "https://data.geopf.fr/wms-r/wms?SERVICE=WMS&Version=1.3.0&Request=GetCapabilities",
			MaxConnections:             6,
			ReadTimeOut:                60,
			ConnectTimeOut:             30,
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateWmStore("foo", &WmsStore{
		XMLName: xml.Name{
			Local: "wmsStore",
		},
		Name:    "sf",
		Enabled: true,
		Workspace: &WorkspaceReference{
			Name: "foo",
		},
		Default:                    false,
		Description:                "Services de la GeoPlateforme IGN",
		DisableConnectionOnFailure: false,
		CapabilitiesUrl:            "https://data.geopf.fr/wms-r/wms?SERVICE=WMS&Version=1.3.0&Request=GetCapabilities",
		MaxConnections:             6,
		ReadTimeOut:                60,
		ConnectTimeOut:             30,
	})

	assert.Nil(t, err)
}

func TestUpdateWmsStoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmsstores/sf")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *WmsStore
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &WmsStore{
			XMLName: xml.Name{
				Local: "wmsStore",
			},
			Name:    "sf",
			Enabled: true,
			Workspace: &WorkspaceReference{
				Name: "foo",
			},
			Default:                    false,
			Description:                "Services de la GeoPlateforme IGN",
			DisableConnectionOnFailure: false,
			CapabilitiesUrl:            "https://data.geopf.fr/wms-r/wms?SERVICE=WMS&Version=1.3.0&Request=GetCapabilities",
			MaxConnections:             6,
			ReadTimeOut:                60,
			ConnectTimeOut:             30,
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateWmsStore("foo", "sf", &WmsStore{
		XMLName: xml.Name{
			Local: "wmsStore",
		},
		Name:    "sf",
		Enabled: true,
		Workspace: &WorkspaceReference{
			Name: "foo",
		},
		Default:                    false,
		Description:                "Services de la GeoPlateforme IGN",
		DisableConnectionOnFailure: false,
		CapabilitiesUrl:            "https://data.geopf.fr/wms-r/wms?SERVICE=WMS&Version=1.3.0&Request=GetCapabilities",
		MaxConnections:             6,
		ReadTimeOut:                60,
		ConnectTimeOut:             30,
	})

	assert.Nil(t, err)
}

func TestDeleteWmsStoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmsstores/sf")
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

	err := cli.DeleteWmsStore("foo", "sf", true)

	assert.Nil(t, err)
}
