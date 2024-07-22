package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeFiletBlobstoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/blobstores/sf", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<FileBlobStore default="false">
			<id>sf</id>
			<enabled>true</enabled>
			<baseDirectory>/diskcache</baseDirectory>
			<fileSystemBlockSize>4096</fileSystemBlockSize>
	  	</FileBlobStore>		
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &BlobstoreFile{
		XMLName: xml.Name{
			Local: "FileBlobStore",
		},
		Id:                  "sf",
		Enabled:             true,
		BaseDirectory:       "/diskcache",
		FileSystemBlockSize: 4096,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetBlobstoreFile("sf")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetFileBlobstoreUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetBlobstoreFile("sf")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, datastore)
}

func TestGetFileBlobstoreNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetBlobstoreFile("sf")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, datastore)
}

func TestGetFileBlobstoreUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastore, err := cli.GetBlobstoreFile("sf")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, datastore)
}

func TestCreateFileBlobstoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *BlobstoreFile
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &BlobstoreFile{
			XMLName: xml.Name{
				Local: "FileBlobStore",
			},
			Enabled: true,
			Id:      "sf",
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateBlobstoreFile("sf", &BlobstoreFile{
		XMLName: xml.Name{
			Local: "FileBlobStore",
		},
		Enabled: true,
		Id:      "sf",
	})

	assert.Nil(t, err)
}

func TestUpdateFileBlobstoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *BlobstoreFile
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &BlobstoreFile{
			XMLName: xml.Name{
				Local: "FileBlobStore",
			},
			Enabled: true,
			Id:      "sf",
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateBlobstoreFile("sf", &BlobstoreFile{
		XMLName: xml.Name{
			Local: "FileBlobStore",
		},
		Enabled: true,
		Id:      "sf",
	})

	assert.Nil(t, err)
}

func TestDeleteFileBlobstoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteBlobstoreFile("sf")

	assert.Nil(t, err)
}
