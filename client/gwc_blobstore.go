package client

import (
	"encoding/xml"
)

// BlobstoreReference is a reference to a Blobstore
type BlobstoreReference struct {
	Name string `xml:"name"`
}

// Blobstores is a list of blobstore reference
type Blobstores struct {
	XMLName xml.Name              `xml:"blobStores"`
	List    []*BlobstoreReference `xml:"blobStore"`
}
