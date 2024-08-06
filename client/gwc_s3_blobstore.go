package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// BlobstoreS3 is a Geoserver object
type BlobstoreS3 struct {
	XMLName        xml.Name `xml:"S3BlobStore"`
	Id             string   `xml:"id"`
	Bucket         string   `xml:"bucket"`
	Prefix         string   `xml:"prefix"`
	AwsAccessKey   string   `xml:"awsAccessKey"`
	AwsSecretKey   string   `xml:"awsSecretKey"`
	Access         string   `xml:"access"`
	Endpoint       string   `xml:"endpoint"`
	MaxConnections int      `xml:"maxConnections"`
	UseHTTPS       bool     `xml:"useHTTPS"`
	UseGzip        bool     `xml:"useGzip"`
	Enabled        bool     `xml:"enabled"`
	Default        bool     `xml:"__default"`
}

// GetBlobstoreS3 return a single S3 datastore based on its name
func (c *Client) GetBlobstoreS3(name string) (blobstore *BlobstoreS3, err error) {
	statusCode, body, err := c.doRequest("GET", fmt.Sprintf("/blobstores/%s", name), nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 404:
		err = fmt.Errorf("not found")
		return
	case 200:
		break
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}

	var data BlobstoreS3
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return blobstore, err
	}

	blobstore = &data

	return
}

// CreateBlobstoreS3 creates a blobstore on S3
func (c *Client) CreateBlobstoreS3(blobstoreName string, blobstore *BlobstoreS3) (err error) {
	payload, _ := xml.Marshal(&blobstore)
	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/blobstores/%s", blobstoreName), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 201:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}

// UpdateBlobstoreS3 updates a blobstore
func (c *Client) UpdateBlobstoreS3(blobstoreName string, blobstore *BlobstoreS3) (err error) {
	payload, _ := xml.Marshal(&blobstore)

	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/blobstores/%s", blobstoreName), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 404:
		err = fmt.Errorf("not found")
		return
	case 405:
		err = fmt.Errorf("forbidden")
		return
	case 200:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}

// DeleteDatastore deletes a datastore
func (c *Client) DeleteBlobstoreS3(blobstoreName string) (err error) {
	statusCode, body, err := c.doRequest("DELETE", fmt.Sprintf("/blobstores/%s", blobstoreName), nil)
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 403:
		err = fmt.Errorf("workspace is not empty")
		return
	case 404:
		err = fmt.Errorf("not found")
		return
	case 405:
		err = fmt.Errorf("forbidden")
		return
	case 200:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}
