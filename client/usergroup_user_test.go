package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsersNoServiceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/usergroup/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<users>
			<user>
				<enabled>true</enabled>
				<userName>admin</userName>
			</user>
		</users>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := Users{
		XMLName: xml.Name{
			Local: "users",
		},
		List: []*User{
			{
				XMLName: xml.Name{
					Local: "user",
				},
				Name:    "admin",
				Enabled: true,
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	users, err := cli.GetUsers("")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, users)
}

func TestGetUsersServiceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/usergroup/service/foo/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<users>
			<user>
				<enabled>true</enabled>
				<userName>admin</userName>
			</user>
		</users>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := Users{
		XMLName: xml.Name{
			Local: "users",
		},
		List: []*User{
			{
				XMLName: xml.Name{
					Local: "user",
				},
				Name:    "admin",
				Enabled: true,
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	users, err := cli.GetUsers("foo")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, users)
}

func TestGetUserNoServiceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/usergroup/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<users>
			<user>
				<enabled>true</enabled>
				<userName>admin</userName>
			</user>
			<user>
				<enabled>true</enabled>
				<userName>smart_admin</userName>
			</user>
		</users>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &User{
		XMLName: xml.Name{
			Local: "user",
		},
		Name:    "smart_admin",
		Enabled: true,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	user, err := cli.GetUser("", "smart_admin")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, user)
}

func TestGetUserServiceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/usergroup/service/foo/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<users>
			<user>
				<enabled>true</enabled>
				<userName>admin</userName>
			</user>
			<user>
				<enabled>true</enabled>
				<userName>smart_admin</userName>
			</user>
		</users>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &User{
		XMLName: xml.Name{
			Local: "user",
		},
		Name:    "smart_admin",
		Enabled: true,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	user, err := cli.GetUser("foo", "smart_admin")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, user)
}

func TestGetUserUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/usergroup/users")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	user, err := cli.GetUser("", "toto")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, user)
}

func TestGetUserNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/usergroup/users")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	user, err := cli.GetUser("", "toto")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, user)
}

func TestGetUserUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/usergroup/users")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	user, err := cli.GetUser("", "toto")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, user)
}

func TestCreateUserNoWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/usergroup/users")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *User
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &User{
			XMLName: xml.Name{
				Local: "user",
			},
			Name:     "admin",
			Enabled:  true,
			Password: "nabucho",
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateUser("", &User{
		XMLName: xml.Name{
			Local: "User",
		},
		Name:     "admin",
		Enabled:  true,
		Password: "nabucho",
	})

	assert.Nil(t, err)
}

func TestUpdateUserNoWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/usergroup/user/admin")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *User
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &User{
			XMLName: xml.Name{
				Local: "user",
			},
			Name:     "admin",
			Enabled:  true,
			Password: "nabucho",
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateUser("", "admin", &User{
		XMLName: xml.Name{
			Local: "User",
		},
		Name:     "admin",
		Enabled:  true,
		Password: "nabucho",
	})

	assert.Nil(t, err)
}

func TestDeleteUserSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/usergroup/user/admin")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteUser("", "admin")

	assert.Nil(t, err)
}
