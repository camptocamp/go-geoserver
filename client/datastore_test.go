package client

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatastoresSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/datastores", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<dataStores>
			<dataStore>
				<name>sf</name>
				<atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/sf/datastores/sf.xml" type="application/xml"/>
			</dataStore>
		</dataStores>
		`))
	})
	mux.HandleFunc("/workspaces/foo/datastores/sf", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<dataStore>
			<name>sf</name>
			<enabled>true</enabled>
			<workspace>
			  <name>foo</name>
			  <atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/foo.xml" type="application/xml"/>
			</workspace>
			<connectionParameters>
			  <entry key="url">file:data/sf</entry>
			  <entry key="namespace">http://www.openplans.org/spearfish</entry>
			</connectionParameters>
			<__default>false</__default>
			<featureTypes>
			  <atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/sf/datastores/sf/featuretypes.xml" type="application/xml"/>
			</featureTypes>
		</dataStore>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*Datastore{
		&Datastore{
			XMLName: xml.Name{
				Local: "dataStore",
			},
			Name:    "sf",
			Enabled: true,
			ConnectionParameters: []*DatastoreConnectionParameter{
				&DatastoreConnectionParameter{
					Key:   "url",
					Value: "file:data/sf",
				},
				&DatastoreConnectionParameter{
					Key:   "namespace",
					Value: "http://www.openplans.org/spearfish",
				},
			},
			Workspace: &WorkspaceReference{
				Name: "foo",
			},
			Default: false,
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetDatastores("foo")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetDatastoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/datastores/sf", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<dataStore>
			<name>sf</name>
			<enabled>true</enabled>
			<workspace>
			  <name>foo</name>
			  <atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/foo.xml" type="application/xml"/>
			</workspace>
			<connectionParameters>
			  <entry key="url">file:data/sf</entry>
			  <entry key="namespace">http://www.openplans.org/spearfish</entry>
			</connectionParameters>
			<__default>false</__default>
			<featureTypes>
			  <atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/sf/datastores/sf/featuretypes.xml" type="application/xml"/>
			</featureTypes>
		</dataStore>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &Datastore{
		XMLName: xml.Name{
			Local: "dataStore",
		},
		Name:    "sf",
		Enabled: true,
		ConnectionParameters: []*DatastoreConnectionParameter{
			&DatastoreConnectionParameter{
				Key:   "url",
				Value: "file:data/sf",
			},
			&DatastoreConnectionParameter{
				Key:   "namespace",
				Value: "http://www.openplans.org/spearfish",
			},
		},
		Workspace: &WorkspaceReference{
			Name: "foo",
		},
		Default: false,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetDatastore("foo", "sf")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetDatastoreUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores/sf")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetDatastore("foo", "sf")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, datastore)
}

func TestGetDatastoreNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores/sf")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetDatastore("foo", "sf")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, datastore)
}

func TestGetDatastoreUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores/sf")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetDatastore("foo", "sf")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, datastore)
}

func TestCreateDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *Datastore
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &Datastore{
			XMLName: xml.Name{
				Local: "dataStore",
			},
			Name:    "sf",
			Enabled: true,
			ConnectionParameters: []*DatastoreConnectionParameter{
				&DatastoreConnectionParameter{
					Key:   "url",
					Value: "file:data/sf",
				},
				&DatastoreConnectionParameter{
					Key:   "namespace",
					Value: "http://www.openplans.org/spearfish",
				},
			},
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateDatastore("foo", &Datastore{
		Name:    "sf",
		Enabled: true,
		ConnectionParameters: []*DatastoreConnectionParameter{
			&DatastoreConnectionParameter{
				Key:   "url",
				Value: "file:data/sf",
			},
			&DatastoreConnectionParameter{
				Key:   "namespace",
				Value: "http://www.openplans.org/spearfish",
			},
		},
	})

	assert.Nil(t, err)
}

func TestUpdateDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores/sf")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *Datastore
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &Datastore{
			XMLName: xml.Name{
				Local: "dataStore",
			},
			Name:    "sf",
			Enabled: true,
			ConnectionParameters: []*DatastoreConnectionParameter{
				&DatastoreConnectionParameter{
					Key:   "url",
					Value: "file:data/sf",
				},
				&DatastoreConnectionParameter{
					Key:   "namespace",
					Value: "http://www.openplans.org/spearfish",
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

	err := cli.UpdateDatastore("foo", "sf", &Datastore{
		Name:    "sf",
		Enabled: true,
		ConnectionParameters: []*DatastoreConnectionParameter{
			&DatastoreConnectionParameter{
				Key:   "url",
				Value: "file:data/sf",
			},
			&DatastoreConnectionParameter{
				Key:   "namespace",
				Value: "http://www.openplans.org/spearfish",
			},
		},
	})

	assert.Nil(t, err)
}

func TestDeleteDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores/sf")
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

	err := cli.DeleteDatastore("foo", "sf", true)

	assert.Nil(t, err)
}
