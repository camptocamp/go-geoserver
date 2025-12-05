package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLayeRulesSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(200)
		w.Write([]byte(`
		<rules>
			<rule resource="*.*.r">ROLE_AUTHENTICATED</rule>
			<rule resource="*.*.w">GROUP_ADMIN,ADMIN</rule>
		</rules>
		`))
	}))
	defer testServer.Close()

	expectedResult := LayerRules{
		XMLName: xml.Name{
			Local: "rules",
		},
		List: []*LayerRule{
			{
				XMLName: xml.Name{
					Local: "rule",
				},
				Resource: "*.*.r",
				Rule:     "ROLE_AUTHENTICATED",
			},
			{
				XMLName: xml.Name{
					Local: "rule",
				},
				Resource: "*.*.w",
				Rule:     "GROUP_ADMIN,ADMIN",
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerRules, err := cli.GetLayerRules()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layerRules)
}

func TestGetLayerRuleSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(200)
		w.Write([]byte(`
		<rules>
			<rule resource="*.*.r">ROLE_AUTHENTICATED</rule>
			<rule resource="*.*.w">GROUP_ADMIN,ADMIN</rule>
		</rules>
		`))
	}))
	defer testServer.Close()

	expectedResult := &LayerRule{
		XMLName: xml.Name{
			Local: "rule",
		},
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerRule, err := cli.GetLayerRule("*.*.w")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layerRule)
}

func TestGetLayerRuleUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerRule, err := cli.GetLayerRule("foo")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, layerRule)
}

func TestGetLayerRuleNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerRule, err := cli.GetLayerRule("foo")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, layerRule)
}

func TestGetLayerRuleUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerRule, err := cli.GetLayerRule("foo")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, layerRule)
}

func TestCreateLayerRuleSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *LayerRules
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &LayerRules{
			XMLName: xml.Name{
				Local: "rules",
			},
			List: []*LayerRule{
				{
					XMLName: xml.Name{
						Local: "rule",
					},
					Resource: "*.*.w",
					Rule:     "GROUP_ADMIN,ADMIN",
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

	err := cli.CreateLayerRule(&LayerRule{
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	})

	assert.Nil(t, err)
}

func TestCreateLayeRuleUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateLayerRule(&LayerRule{
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	})

	assert.Error(t, err, "Unauthorized")
}

func TestCreateLayerRuleUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateLayerRule(&LayerRule{
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	})

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}

func TestUpdateLayerRuleSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *LayerRules
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &LayerRules{
			XMLName: xml.Name{
				Local: "rules",
			},
			List: []*LayerRule{
				{
					XMLName: xml.Name{
						Local: "rule",
					},
					Resource: "*.*.w",
					Rule:     "GROUP_ADMIN,ADMIN",
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

	err := cli.UpdateLayerRule(&LayerRule{
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	})

	assert.Nil(t, err)
}

func TestUpdateLayerRuleUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateLayerRule(&LayerRule{
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	})

	assert.Error(t, err, "Unauthorized")
}

func TestUpdateLayerRuleNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(409)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateLayerRule(&LayerRule{
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	})

	assert.Error(t, err, "Not Found")
}

func TestUpdateLayerRuleUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/security/acl/layers")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateLayerRule(&LayerRule{
		Resource: "*.*.w",
		Rule:     "GROUP_ADMIN,ADMIN",
	})

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}

func TestDeleteLayerRuleSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/security/acl/layers/*.*.w")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteLayerRule("*.*.w")

	assert.Nil(t, err)
}

func TestDeleteLayerRuleUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/security/acl/layers/*.*.w")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteLayerRule("*.*.w")

	assert.Error(t, err, "Unauthorized")
}

func TestDeleteLayerRuleNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/security/acl/layers/*.*.w")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteLayerRule("*.*.w")

	assert.Error(t, err, "Not Found")
}

func TestDeleteLayerRuleForbidden(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/security/acl/layers/*.*.w")

		w.WriteHeader(405)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteLayerRule("*.*.w")

	assert.Error(t, err, "forbidden")
}

func TestDeleteLayerRuleUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/security/acl/layers/*.*.w")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteLayerRule("*.*.w")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}
