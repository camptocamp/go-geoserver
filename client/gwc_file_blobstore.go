package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// BlobstoreFile is a Geoserver object
type BlobstoreFile struct {
	XMLName             xml.Name `xml:"FileBlobStore"`
	Id                  string   `xml:"id"`
	Enabled             bool     `xml:"enabled"`
	BaseDirectory       string   `xml:"baseDirectory"`
	FileSystemBlockSize int      `xml:"fileSystemBlockSize"`
}

// GetBlobstoreFile return a single File datastore based on its name
func (c *Client) GetBlobstoreFile(name string) (blobstore *BlobstoreFile, err error) {
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

	var data BlobstoreFile
	if err := xml.Unmarshal([]byte(body), &data); err != nil {
		return blobstore, err
	}

	blobstore = &data

	return
}

// CreateBlobstoreFile creates a blobstore on disk
func (c *Client) CreateBlobstoreFile(blobstoreName string, blobstore *BlobstoreFile) (err error) {
	payload, _ := xml.Marshal(&blobstore)
	statusCode, body, err := c.doRequest("PUT", fmt.Sprintf("/blobstores/%s", blobstoreName), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	switch statusCode {
	case 401:
		err = fmt.Errorf("unauthorized")
		return
	case 200:
		return
	case 201:
		return
	default:
		err = fmt.Errorf("unknown error: %d - %s", statusCode, body)
		return
	}
}

// UpdateBlobstoreFile updates a blobstore
func (c *Client) UpdateBlobstoreFile(blobstoreName string, blobstore *BlobstoreFile) (err error) {
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
func (c *Client) DeleteBlobstoreFile(blobstoreName string) (err error) {
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
