package client

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrlChecksSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/urlchecks", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
<urlChecks>
    <urlCheck>
        <name>icons</name>
        <atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/urlchecks/icons.xml" type="application/atom+xml"/>
    </urlCheck>
</urlChecks>
		`))
	})
	mux.HandleFunc("/urlchecks/icons", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
<regexUrlCheck>
    <name>icons</name>
    <description>External graphic icons</description>
    <enabled>true</enabled>
    <regex>^https://styles.server.net/icons/.*$</regex>
</regexUrlCheck>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*RegexUrlCheck{
		{
			XMLName: xml.Name{
				Space: "",
				Local: "regexUrlCheck",
			},
			Name:        "icons",
			Description: "External graphic icons",
			Regex:       "^https://styles.server.net/icons/.*$",
			IsEnabled:   true,
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	urlChecks, err := cli.GetUrlChecks()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, urlChecks)
}

func TestGetUrlChecksUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/urlchecks")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	urlChecks, err := cli.GetUrlChecks()

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, urlChecks)
}

func TestGetUrlChecksUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/urlchecks")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	urlChecks, err := cli.GetUrlChecks()

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, urlChecks)
}

func TestGetCheckSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(200)
		w.Write([]byte(`
<regexUrlCheck>
    <name>icons</name>
    <description>External graphic icons</description>
    <enabled>true</enabled>
    <regex>^https://styles.server.net/icons/.*$</regex>
</regexUrlCheck>
		`))
	}))
	defer testServer.Close()

	expectedResult := &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
		IsEnabled:   true,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	urlCheck, err := cli.GetRegExUrlCheck("icons")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, urlCheck)
}

func TestGetCheckUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/urlchecks/foo")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	urlCheck, err := cli.GetRegExUrlCheck("foo")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, urlCheck)
}

func TestGetCheckNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/urlchecks/foo")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	urlCheck, err := cli.GetRegExUrlCheck("foo")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, urlCheck)
}

func TestGetCheckUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/urlchecks/foo")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	urlCheck, err := cli.GetRegExUrlCheck("foo")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, urlCheck)
}

func TestCreateRegexCheckSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/urlchecks")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *RegexUrlCheck
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &RegexUrlCheck{
			XMLName: xml.Name{
				Space: "",
				Local: "regexUrlCheck",
			},
			Name:        "icons",
			Description: "External graphic icons",
			Regex:       "^https://styles.server.net/icons/.*$",
			IsEnabled:   true,
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateRegExUrlCheck("icons", &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
		IsEnabled:   true,
	})

	assert.Nil(t, err)
}

func TestCreateRegexUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/urlchecks")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateRegExUrlCheck("icons", &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
	})

	assert.Error(t, err, "Unauthorized")
}

func TestCreateRegexUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/urlchecks")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateRegExUrlCheck("icons", &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
		IsEnabled:   true,
	})

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}

func TestUpdateRegExSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *RegexUrlCheck
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &RegexUrlCheck{
			XMLName: xml.Name{
				Space: "",
				Local: "regexUrlCheck",
			},
			Name:        "icons",
			Description: "External graphic icons",
			Regex:       "^https://styles.server.net/icons/.*$",
			IsEnabled:   true,
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateRegExUrlCheck("icons", &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
		IsEnabled:   true,
	})

	assert.Nil(t, err)
}

func TestUpdateRegExUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateRegExUrlCheck("icons", &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
		IsEnabled:   true,
	})

	assert.Error(t, err, "Unauthorized")
}

func TestUpdateRegExNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateRegExUrlCheck("icons", &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
		IsEnabled:   true,
	})

	assert.Error(t, err, "Not Found")
}

func TestUpdateRegExUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateRegExUrlCheck("icons", &RegexUrlCheck{
		XMLName: xml.Name{
			Space: "",
			Local: "regexUrlCheck",
		},
		Name:        "icons",
		Description: "External graphic icons",
		Regex:       "^https://styles.server.net/icons/.*$",
		IsEnabled:   true,
	})

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}

func TestDeleteURLCheckSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteUrlCheck("icons")

	assert.Nil(t, err)
}

func TestDeleteURLCheckUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteUrlCheck("icons")

	assert.Error(t, err, "Unauthorized")
}

func TestDeleteURLCheckNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteUrlCheck("icons")

	assert.Error(t, err, "Not Found")
}

func TestDeleteURLCheckForbidden(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(405)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteUrlCheck("icons")

	assert.Error(t, err, "forbidden")
}

func TestDeleteURLCheckUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/urlchecks/icons")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteUrlCheck("icons")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
}
