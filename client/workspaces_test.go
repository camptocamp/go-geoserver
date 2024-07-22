package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWorkspacesSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces")

		w.WriteHeader(200)
		w.Write([]byte(`
		<workspaces>
			<workspace>
				<name>topp</name>
			</workspace>
			<workspace>
				<name>it.geosolutions</name>
			</workspace>
		</workspaces>
		`))
	}))
	defer testServer.Close()

	expectedResult := []*Workspace{
		{
			Name: "topp",
		},
		{
			Name: "it.geosolutions",
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	workspaces, err := cli.GetWorkspaces()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, workspaces)
}

func TestGetWorkspacesUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	workspaces, err := cli.GetWorkspaces()

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, workspaces)
}

func TestGetWorkspacesUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	workspaces, err := cli.GetWorkspaces()

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, workspaces)
}

func TestGetWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/topp")

		w.WriteHeader(200)
		w.Write([]byte(`
		<workspace>
  			<name>topp</name>
  			<isolated>false</isolated>
		</workspace>
		`))
	}))
	defer testServer.Close()

	expectedResult := &Workspace{
		XMLName: xml.Name{
			Local: "workspace",
		},
		Name: "topp",
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	workspace, err := cli.GetWorkspace("topp")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, workspace)
}

func TestGetWorkspaceUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	workspaces, err := cli.GetWorkspace("foo")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, workspaces)
}

func TestGetWorkspaceNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	workspaces, err := cli.GetWorkspace("foo")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, workspaces)
}

func TestGetWorkspaceUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	workspaces, err := cli.GetWorkspace("foo")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, workspaces)
}

func TestCreateWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces")
		keys, ok := r.URL.Query()["default"]
		assert.True(t, ok)
		assert.Equal(t, keys[0], "true")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *Workspace
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &Workspace{
			XMLName: xml.Name{
				Local: "workspace",
			},
			Name: "foo",
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateWorkspace(&Workspace{
		Name:     "foo",
		Isolated: false,
	}, true)

	assert.Nil(t, err)
}

func TestCreateWorkspaceNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateWorkspace(&Workspace{
		Name: "foo",
	}, true)

	assert.Error(t, err, "Unauthorized")
}

func TestCreateWorkspaceUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateWorkspace(&Workspace{
		Name: "foo",
	}, true)

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}

func TestUpdateWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *Workspace
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &Workspace{
			XMLName: xml.Name{
				Local: "workspace",
			},
			Name:     "foo",
			Isolated: true,
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateWorkspace("foo", &Workspace{
		Name:     "foo",
		Isolated: true,
	})

	assert.Nil(t, err)
}

func TestUpdateWorkspaceUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateWorkspace("foo", &Workspace{})

	assert.Error(t, err, "Unauthorized")
}

func TestUpdateWorkspaceNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateWorkspace("foo", &Workspace{})

	assert.Error(t, err, "Not Found")
}

func TestUpdateWorkspaceUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateWorkspace("foo", &Workspace{})

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}

func TestDeleteWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")
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

	err := cli.DeleteWorkspace("foo", true)

	assert.Nil(t, err)
}

func TestDeleteWorkspaceUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteWorkspace("foo", false)

	assert.Error(t, err, "Unauthorized")
}

func TestDeleteWorkspaceNotEmpty(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(403)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteWorkspace("foo", false)

	assert.Error(t, err, "Workspace is not empty")
}

func TestDeleteWorkspaceNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteWorkspace("foo", false)

	assert.Error(t, err, "Not Found")
}

func TestDeleteWorkspaceForbidden(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(405)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteWorkspace("foo", false)

	assert.Error(t, err, "Cannot delete default workspace")
}

func TestDeleteWorkspaceUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteWorkspace("foo", false)

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}
