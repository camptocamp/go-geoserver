package client

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlobstoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/blobstores/sf", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<S3BlobStore default="false">
		  <id>sf</id>
		  <enabled>true</enabled>
		  <bucket>the_bucket</bucket>
		  <prefix>gwc_master</prefix>
		  <awsAccessKey>MyWonderfulAccessKey</awsAccessKey>
		  <awsSecretKey>MyWonderfulSecretKey</awsSecretKey>
		  <access>PRIVATE</access>
		  <maxConnections>50</maxConnections>
		  <useHTTPS>true</useHTTPS>
		  <useGzip>false</useGzip>
		  <endpoint>https://url.to.endpoint</endpoint>
		</S3BlobStore>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &BlobstoreS3{
		XMLName: xml.Name{
			Local: "S3BlobStore",
		},
		Id:             "sf",
		Enabled:        true,
		Bucket:         "the_bucket",
		Prefix:         "gwc_master",
		AwsAccessKey:   "MyWonderfulAccessKey",
		AwsSecretKey:   "MyWonderfulSecretKey",
		Access:         "PRIVATE",
		MaxConnections: 50,
		UseHTTPS:       true,
		UseGzip:        false,
		Endpoint:       "https://url.to.endpoint",
		Default:        false,
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	datastores, err := cli.GetBlobstoreS3("sf")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, datastores)
}

func TestGetS3BlobstoreUnauthorized(t *testing.T) {
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

	datastore, err := cli.GetBlobstoreS3("sf")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, datastore)
}

func TestGetS3BlobstoreNotFound(t *testing.T) {
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

	datastore, err := cli.GetBlobstoreS3("sf")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, datastore)
}

func TestGetS3BlobstoreUnknownError(t *testing.T) {
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

	datastore, err := cli.GetBlobstoreS3("sf")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, datastore)
}

func TestCreateS3BlobstoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *BlobstoreS3
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &BlobstoreS3{
			XMLName: xml.Name{
				Local: "S3BlobStore",
			},
			Enabled: true,
			Id:      "sf",
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateBlobstoreS3("sf", &BlobstoreS3{
		XMLName: xml.Name{
			Local: "S3BlobStore",
		},
		Enabled: true,
		Id:      "sf",
	})

	assert.Nil(t, err)
}

func TestUpdateS3BlobstoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/blobstores/sf")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *BlobstoreS3
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &BlobstoreS3{
			XMLName: xml.Name{
				Local: "S3BlobStore",
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

	err := cli.UpdateBlobstoreS3("sf", &BlobstoreS3{
		XMLName: xml.Name{
			Local: "S3BlobStore",
		},
		Enabled: true,
		Id:      "sf",
	})

	assert.Nil(t, err)
}

func TestDeleteS3BlobstoreSuccess(t *testing.T) {
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

	err := cli.DeleteBlobstoreS3("sf")

	assert.Nil(t, err)
}
