package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFeatureTypesNoDatastoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/featuretypes", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		{
			"featureTypes": {
				"featureType": [
					{
						"name": "toto",
						"href": "http:\/\/localhost:8080"
					}
				] 
			}
		}
		`))
	})
	mux.HandleFunc("/workspaces/foo/featuretypes/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		{
			"featureType": {
				"name": "toto"
			}
		}
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*FeatureType{
		&FeatureType{
			Name: "toto",
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	featureTypes, err := cli.GetFeatureTypes("foo", "")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, featureTypes)
}

func TestGetFeatureTypesInDatastoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/datastores/bar/featuretypes", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		{
			"featureTypes": {
				"featureType": [
					{
						"name": "toto",
						"href": "http:\/\/localhost:8080"
					}
				] 
			}
		}
		`))
	})
	mux.HandleFunc("/workspaces/foo/datastores/bar/featuretypes/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		{
			"featureType": {
				"name": "toto"
			}
		}
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*FeatureType{
		&FeatureType{
			Name: "toto",
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	featuretypes, err := cli.GetFeatureTypes("foo", "bar")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, featuretypes)
}

func TestGetFeatureTypeNoDatastoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/featuretypes/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		{
			"featureType": {
			  "name": "toto",
			  "nativeName": "toto",
			  "namespace": {
				"name": "tiger",
				"href": "http://localhost:8080/geoserver/rest/namespaces/tiger.json"
			  },
			  "title": "Manhattan (NY) points of interest",
			  "abstract": "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
			  "keywords": {
				"string": [
				  "poi",
				  "Manhattan",
				  "DS_poi",
				  "points_of_interest",
				  "sampleKeyword\\@language=ab\\;",
				  "area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
				  "Привет\\@language=ru\\;\\@vocabulary=friendly\\;"
				]
			  },
			  "metadataLinks": {
				"metadataLink": [
				  {
					"type": "text/plain",
					"metadataType": "FGDC",
					"content": "www.google.com"
				  }
				]
			  },
			  "dataLinks": {
				"org.geoserver.catalog.impl.DataLinkInfoImpl": [
				  {
					"type": "text/plain",
					"content": "http://www.google.com"
				  }
				]
			  },
			  "nativeCRS": "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
			  "srs": "EPSG:4326",
			  "nativeBoundingBox": {
				"minx": -74.0118315772888,
				"maxx": -74.00153046439813,
				"miny": 40.70754683896324,
				"maxy": 40.719885123828675,
				"crs": "EPSG:4326"
			  },
			  "latLonBoundingBox": {
				"minx": -74.0118315772888,
				"maxx": -74.00857344353275,
				"miny": 40.70754683896324,
				"maxy": 40.711945649065406,
				"crs": "EPSG:4326"
			  },
			  "projectionPolicy": "REPROJECT_TO_DECLARED",
			  "enabled": true,
			  "metadata": {
				"entry": [
				  {
					"@key": "kml.regionateStrategy",
					"$": "external-sorting"
				  },
				  {
					"@key": "kml.regionateFeatureLimit",
					"$": "15"
				  },
				  {
					"@key": "cacheAgeMax",
					"$": "3000"
				  },
				  {
					"@key": "cachingEnabled",
					"$": "true"
				  },
				  {
					"@key": "kml.regionateAttribute",
					"$": "NAME"
				  },
				  {
					"@key": "indexingEnabled",
					"$": "false"
				  },
				  {
					"@key": "dirName",
					"$": "DS_poi_poi"
				  }
				]
			  },
			  "store": {
				"@class": "dataStore",
				"name": "tiger:nyc",
				"href": "http://localhost:8080/geoserver/rest/workspaces/tiger/datastores/nyc.json"
			  },
			  "cqlFilter": "INCLUDE",
			  "maxFeatures": 100,
			  "numDecimals": 6,
			  "responseSRS": {
				"string": [
				  4326
				]
			  },
			  "overridingServiceSRS": true,
			  "skipNumberMatched": true,
			  "circularArcPresent": true,
			  "linearizationTolerance": 10,
			  "attributes": {
				"attribute": [
				  {
					"name": "the_geom",
					"minOccurs": 0,
					"maxOccurs": 1,
					"nillable": true,
					"binding": "org.locationtech.jts.geom.Point"
				  }
				]
			  }
			  }
		}
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &FeatureType{
		Name:       "toto",
		NativeName: "toto",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: map[string][]string{
			"string": []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"sampleKeyword\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
			},
		},
		NativeCRS: CRSWrapper{
			Class: "",
			Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "REPROJECT_TO_DECLARED",
		Enabled:          true,
		Attributes: FeatureTypeAttributesList{
			Attribute: []*FeatureTypeAttribute{
				&FeatureTypeAttribute{
					Name:      "the_geom",
					MinOccurs: 0,
					MaxOccurs: 1,
					Nillable:  true,
					Binding:   "org.locationtech.jts.geom.Point",
				},
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	featureType, err := cli.GetFeatureType("foo", "", "toto")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, featureType)
}

func TestGetFeatureTypeInDatastoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/datastores/bar/featuretypes/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		{
			"featureType": {
			  "name": "toto",
			  "nativeName": "toto",
			  "namespace": {
				"name": "tiger",
				"href": "http://localhost:8080/geoserver/rest/namespaces/tiger.json"
			  },
			  "title": "Manhattan (NY) points of interest",
			  "abstract": "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
			  "keywords": {
				"string": [
				  "poi",
				  "Manhattan",
				  "DS_poi",
				  "points_of_interest",
				  "sampleKeyword\\@language=ab\\;",
				  "area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
				  "Привет\\@language=ru\\;\\@vocabulary=friendly\\;"
				]
			  },
			  "metadataLinks": {
				"metadataLink": [
				  {
					"type": "text/plain",
					"metadataType": "FGDC",
					"content": "www.google.com"
				  }
				]
			  },
			  "dataLinks": {
				"org.geoserver.catalog.impl.DataLinkInfoImpl": [
				  {
					"type": "text/plain",
					"content": "http://www.google.com"
				  }
				]
			  },
			  "nativeCRS": "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
			  "srs": "EPSG:4326",
			  "nativeBoundingBox": {
				"minx": -74.0118315772888,
				"maxx": -74.00153046439813,
				"miny": 40.70754683896324,
				"maxy": 40.719885123828675,
				"crs": "EPSG:4326"
			  },
			  "latLonBoundingBox": {
				"minx": -74.0118315772888,
				"maxx": -74.00857344353275,
				"miny": 40.70754683896324,
				"maxy": 40.711945649065406,
				"crs": "EPSG:4326"
			  },
			  "projectionPolicy": "REPROJECT_TO_DECLARED",
			  "enabled": true,
			  "metadata": {
				"entry": [
				  {
					"@key": "kml.regionateStrategy",
					"$": "external-sorting"
				  },
				  {
					"@key": "kml.regionateFeatureLimit",
					"$": "15"
				  },
				  {
					"@key": "cacheAgeMax",
					"$": "3000"
				  },
				  {
					"@key": "cachingEnabled",
					"$": "true"
				  },
				  {
					"@key": "kml.regionateAttribute",
					"$": "NAME"
				  },
				  {
					"@key": "indexingEnabled",
					"$": "false"
				  },
				  {
					"@key": "dirName",
					"$": "DS_poi_poi"
				  }
				]
			  },
			  "store": {
				"@class": "dataStore",
				"name": "tiger:nyc",
				"href": "http://localhost:8080/geoserver/rest/workspaces/tiger/datastores/nyc.json"
			  },
			  "cqlFilter": "INCLUDE",
			  "maxFeatures": 100,
			  "numDecimals": 6,
			  "responseSRS": {
				"string": [
				  4326
				]
			  },
			  "overridingServiceSRS": true,
			  "skipNumberMatched": true,
			  "circularArcPresent": true,
			  "linearizationTolerance": 10,
			  "attributes": {
				"attribute": {
					"name": "the_geom",
					"minOccurs": 0,
					"maxOccurs": 1,
					"nillable": true,
					"binding": "org.locationtech.jts.geom.Point"
				  }
			  }
		  }
		}
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &FeatureType{
		Name:       "toto",
		NativeName: "toto",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: map[string][]string{
			"string": []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"sampleKeyword\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
			},
		},
		NativeCRS: CRSWrapper{
			Class: "",
			Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "REPROJECT_TO_DECLARED",
		Enabled:          true,
		Attributes: FeatureTypeAttributesList{
			Attribute: []*FeatureTypeAttribute{
				&FeatureTypeAttribute{
					Name:      "the_geom",
					MinOccurs: 0,
					MaxOccurs: 1,
					Nillable:  true,
					Binding:   "org.locationtech.jts.geom.Point",
				},
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	featureType, err := cli.GetFeatureType("foo", "bar", "toto")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, featureType)
}

func TestGetFeatureTypeUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/featuretypes/toto")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	featureType, err := cli.GetFeatureType("foo", "", "toto")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, featureType)
}

func TestGetFeatureTypeNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/featuretypes/toto")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	featureType, err := cli.GetFeatureType("foo", "", "toto")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, featureType)
}

func TestGetFeatureTypeUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/featuretypes/toto")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	featureType, err := cli.GetFeatureType("foo", "", "toto")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, featureType)
}

func TestCreateFeatureTypeNoDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/featuretypes")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload map[string]*FeatureType
		err = json.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, map[string]*FeatureType{
			"featureType": &FeatureType{
				Name:       "toto",
				NativeName: "toto",
				Title:      "Manhattan (NY) points of interest",
				Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
				Keywords: map[string][]string{
					"string": []string{
						"poi",
						"Manhattan",
						"DS_poi",
						"points_of_interest",
						"sampleKeyword\\@language=ab\\;",
						"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
						"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
					},
				},
				NativeCRS: CRSWrapper{
					Class: "",
					Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
				},
				SRS: "EPSG:4326",
				NativeBoundingBox: BoundingBox{
					MinX: -74.0118315772888,
					MaxX: -74.00153046439813,
					MinY: 40.70754683896324,
					MaxY: 40.719885123828675,
					CRS: CRSWrapper{
						Class: "",
						Value: "EPSG:4326",
					},
				},
				LatLonBoundingBox: BoundingBox{
					MinX: -74.0118315772888,
					MaxX: -74.00857344353275,
					MinY: 40.70754683896324,
					MaxY: 40.711945649065406,
					CRS: CRSWrapper{
						Class: "",
						Value: "EPSG:4326",
					},
				},
				ProjectionPolicy: "REPROJECT_TO_DECLARED",
				Enabled:          true,
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

	err := cli.CreateFeatureType("foo", "", &FeatureType{
		Name:       "toto",
		NativeName: "toto",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: map[string][]string{
			"string": []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"sampleKeyword\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
			},
		},
		NativeCRS: CRSWrapper{
			Class: "",
			Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "REPROJECT_TO_DECLARED",
		Enabled:          true,
	})

	assert.Nil(t, err)
}

func TestCreateFeatureTypeInDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores/bar/featuretypes")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload map[string]*FeatureType
		err = json.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, map[string]*FeatureType{
			"featureType": &FeatureType{
				Name:       "toto",
				NativeName: "toto",
				Title:      "Manhattan (NY) points of interest",
				Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
				Keywords: map[string][]string{
					"string": []string{
						"poi",
						"Manhattan",
						"DS_poi",
						"points_of_interest",
						"sampleKeyword\\@language=ab\\;",
						"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
						"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
					},
				},
				NativeCRS: CRSWrapper{
					Class: "",
					Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
				},
				SRS: "EPSG:4326",
				NativeBoundingBox: BoundingBox{
					MinX: -74.0118315772888,
					MaxX: -74.00153046439813,
					MinY: 40.70754683896324,
					MaxY: 40.719885123828675,
					CRS: CRSWrapper{
						Class: "",
						Value: "EPSG:4326",
					},
				},
				LatLonBoundingBox: BoundingBox{
					MinX: -74.0118315772888,
					MaxX: -74.00857344353275,
					MinY: 40.70754683896324,
					MaxY: 40.711945649065406,
					CRS: CRSWrapper{
						Class: "",
						Value: "EPSG:4326",
					},
				},
				ProjectionPolicy: "REPROJECT_TO_DECLARED",
				Enabled:          true,
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

	err := cli.CreateFeatureType("foo", "bar", &FeatureType{
		Name:       "toto",
		NativeName: "toto",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: map[string][]string{
			"string": []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"sampleKeyword\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
			},
		},
		NativeCRS: CRSWrapper{
			Class: "",
			Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "REPROJECT_TO_DECLARED",
		Enabled:          true,
	})

	assert.Nil(t, err)
}

func TestUpdateFeatureTypeNoDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/featuretypes/toto")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload map[string]*FeatureType
		err = json.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, map[string]*FeatureType{
			"featureType": &FeatureType{
				Name:       "toto",
				NativeName: "toto",
				Title:      "Manhattan (NY) points of interest",
				Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
				Keywords: map[string][]string{
					"string": []string{
						"poi",
						"Manhattan",
						"DS_poi",
						"points_of_interest",
						"sampleKeyword\\@language=ab\\;",
						"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
						"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
					},
				},
				NativeCRS: CRSWrapper{
					Class: "",
					Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
				},
				SRS: "EPSG:4326",
				NativeBoundingBox: BoundingBox{
					MinX: -74.0118315772888,
					MaxX: -74.00153046439813,
					MinY: 40.70754683896324,
					MaxY: 40.719885123828675,
					CRS: CRSWrapper{
						Class: "",
						Value: "EPSG:4326",
					},
				},
				LatLonBoundingBox: BoundingBox{
					MinX: -74.0118315772888,
					MaxX: -74.00857344353275,
					MinY: 40.70754683896324,
					MaxY: 40.711945649065406,
					CRS: CRSWrapper{
						Class: "",
						Value: "EPSG:4326",
					},
				},
				ProjectionPolicy: "REPROJECT_TO_DECLARED",
				Enabled:          true,
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

	err := cli.UpdateFeatureType("foo", "", "toto", &FeatureType{
		Name:       "toto",
		NativeName: "toto",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: map[string][]string{
			"string": []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"sampleKeyword\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=friendly\\;",
			},
		},
		NativeCRS: CRSWrapper{
			Class: "",
			Value: "GEOGCS[\"WGS 84\", \n  DATUM[\"World Geodetic System 1984\", \n    SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], \n    AUTHORITY[\"EPSG\",\"6326\"]], \n  PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], \n  UNIT[\"degree\", 0.017453292519943295], \n  AXIS[\"Geodetic longitude\", EAST], \n  AXIS[\"Geodetic latitude\", NORTH], \n  AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: CRSWrapper{
				Class: "",
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "REPROJECT_TO_DECLARED",
		Enabled:          true,
	})

	assert.Nil(t, err)
}

func TestDeleteFeatureTypeSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/featuretypes/toto")
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

	err := cli.DeleteFeatureType("foo", "", "toto", true)

	assert.Nil(t, err)
}
