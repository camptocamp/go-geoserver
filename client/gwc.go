package client

type GridSubset struct {
	Name          string `xml:"gridSetName"`
	MinCacheLevel int    `xml:"minCachedLevel,omitempty"`
	MaxCacheLevel int    `xml:"maxCachedLevel,omitempty"`
}

// ScaleNames is a XML object for scale names
type MimeFormats struct {
	Formats []string `xml:"string"`
}
