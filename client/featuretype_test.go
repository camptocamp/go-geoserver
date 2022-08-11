package client

import (
	"encoding/xml"
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
		<featureTypes>
			<featureType>
				<name>toto</name>
			</featureType>
		</featureTypes>
		`))
	})
	mux.HandleFunc("/workspaces/foo/featuretypes/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<featureType>
			<name>toto</name>
		</featureType>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*FeatureType{
		{
			XMLName: xml.Name{
				Space: "",
				Local: "featureType",
			},
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
		<featureTypes>
			<featureType>
				<name>toto</name>
			</featureType>
		</featureTypes>
		`))
	})
	mux.HandleFunc("/workspaces/foo/datastores/bar/featuretypes/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<featureType>
			<name>toto</name>
		</featureType>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*FeatureType{
		{
			XMLName: xml.Name{
				Local: "featureType",
			},
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
		<featureType>
			<name>poi</name>
			<nativeName>poi</nativeName>
			<namespace>
				<name>tiger</name>
				<atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/namespaces/tiger.xml" type="application/xml"/>
			</namespace>
			<title>Manhattan (NY) points of interest</title>
			<abstract>Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.</abstract>
			<keywords>
				<string>poi</string>
				<string>Manhattan</string>
				<string>DS_poi</string>
				<string>points_of_interest</string>
				<string>fred\@language=ab\;</string>
				<string>area of effect\@language=bg\;\@vocabulary=Technical\;</string>
				<string>Привет\@language=ru\;\@vocabulary=Friendly\;</string>
			</keywords>
			<metadataLinks>
				<metadataLink>
					<type>text/plain</type>
					<metadataType>FGDC</metadataType>
					<content>http://www.google.com</content>
				</metadataLink>
			</metadataLinks>
			<dataLinks>
				<org.geoserver.catalog.impl.DataLinkInfoImpl>
					<type>text/plain</type>
					<content>http://www.google.com</content>
				</org.geoserver.catalog.impl.DataLinkInfoImpl>
			</dataLinks>
			<nativeCRS>GEOGCS["WGS 84", DATUM["World Geodetic System 1984", SPHEROID["WGS 84", 6378137.0, 298.257223563, AUTHORITY["EPSG","7030"]], AUTHORITY["EPSG","6326"]], PRIMEM["Greenwich", 0.0, AUTHORITY["EPSG","8901"]], UNIT["degree", 0.017453292519943295], AXIS["Geodetic longitude", EAST], AXIS["Geodetic latitude", NORTH], AUTHORITY["EPSG","4326"]]</nativeCRS>
			<srs>EPSG:4326</srs>
			<nativeBoundingBox>
				<minx>-74.0118315772888</minx>
				<maxx>-74.00153046439813</maxx>
				<miny>40.70754683896324</miny>
				<maxy>40.719885123828675</maxy>
				<crs>EPSG:4326</crs>
			</nativeBoundingBox>
			<latLonBoundingBox>
				<minx>-74.0118315772888</minx>
				<maxx>-74.00857344353275</maxx>
				<miny>40.70754683896324</miny>
				<maxy>40.711945649065406</maxy>
				<crs>EPSG:4326</crs>
			</latLonBoundingBox>
			<projectionPolicy>NONE</projectionPolicy>
			<enabled>true</enabled>
			<metadata>
				<entry key="time"><dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo></entry>
				<entry key="cachingEnabled">true</entry>
			</metadata>
			<store class="dataStore">
				<name>tiger:nyc</name>
				<atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/tiger/datastores/nyc.xml" type="application/xml"/>
			</store>
			<cqlFilter>INCLUDE</cqlFilter>
			<maxFeatures>100</maxFeatures>
			<numDecimals>6</numDecimals>
			<responseSRS>
				<string>4326</string>
			</responseSRS>
			<overridingServiceSRS>true</overridingServiceSRS>
			<skipNumberMatched>true</skipNumberMatched>
			<circularArcPresent>true</circularArcPresent>
			<linearizationTolerance>10</linearizationTolerance>
			<attributes>
				<attribute>
					<name>the_geom</name>
					<minOccurs>0</minOccurs>
					<maxOccurs>1</maxOccurs>
					<nillable>true</nillable>
					<binding>org.locationtech.jts.geom.Point</binding>
				</attribute>
			</attributes>
		</featureType>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &FeatureType{
		XMLName: xml.Name{
			Local: "featureType",
		},
		Name:       "poi",
		NativeName: "poi",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: FeatureTypeKeywords{
			Keywords: []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"fred\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
			},
		},
		NativeCRS: FeatureTypeCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "NONE",
		Enabled:          true,
		Metadata: []*FeatureTypeMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
		Attributes: []*FeatureTypeAttribute{
			{
				Name:      "the_geom",
				MinOccurs: 0,
				MaxOccurs: 1,
				Nillable:  true,
				Binding:   "org.locationtech.jts.geom.Point",
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
		<featureType>
			<name>poi</name>
			<nativeName>poi</nativeName>
			<namespace>
				<name>tiger</name>
				<atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/namespaces/tiger.xml" type="application/xml"/>
			</namespace>
			<title>Manhattan (NY) points of interest</title>
			<abstract>Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.</abstract>
			<keywords>
				<string>poi</string>
				<string>Manhattan</string>
				<string>DS_poi</string>
				<string>points_of_interest</string>
				<string>fred\@language=ab\;</string>
				<string>area of effect\@language=bg\;\@vocabulary=Technical\;</string>
				<string>Привет\@language=ru\;\@vocabulary=Friendly\;</string>
			</keywords>
			<metadataLinks>
				<metadataLink>
					<type>text/plain</type>
					<metadataType>FGDC</metadataType>
					<content>http://www.google.com</content>
				</metadataLink>
			</metadataLinks>
			<dataLinks>
				<org.geoserver.catalog.impl.DataLinkInfoImpl>
					<type>text/plain</type>
					<content>http://www.google.com</content>
				</org.geoserver.catalog.impl.DataLinkInfoImpl>
			</dataLinks>
			<nativeCRS>GEOGCS["WGS 84", DATUM["World Geodetic System 1984", SPHEROID["WGS 84", 6378137.0, 298.257223563, AUTHORITY["EPSG","7030"]], AUTHORITY["EPSG","6326"]], PRIMEM["Greenwich", 0.0, AUTHORITY["EPSG","8901"]], UNIT["degree", 0.017453292519943295], AXIS["Geodetic longitude", EAST], AXIS["Geodetic latitude", NORTH], AUTHORITY["EPSG","4326"]]</nativeCRS>
			<srs>EPSG:4326</srs>
			<nativeBoundingBox>
				<minx>-74.0118315772888</minx>
				<maxx>-74.00153046439813</maxx>
				<miny>40.70754683896324</miny>
				<maxy>40.719885123828675</maxy>
				<crs>EPSG:4326</crs>
			</nativeBoundingBox>
			<latLonBoundingBox>
				<minx>-74.0118315772888</minx>
				<maxx>-74.00857344353275</maxx>
				<miny>40.70754683896324</miny>
				<maxy>40.711945649065406</maxy>
				<crs>EPSG:4326</crs>
			</latLonBoundingBox>
			<projectionPolicy>NONE</projectionPolicy>
			<enabled>true</enabled>
			<metadata>
				<entry key="time"><dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo></entry>
				<entry key="cachingEnabled">true</entry>
			</metadata>
			<store class="dataStore">
				<name>tiger:nyc</name>
				<atom:link xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/tiger/datastores/nyc.xml" type="application/xml"/>
			</store>
			<cqlFilter>INCLUDE</cqlFilter>
			<maxFeatures>100</maxFeatures>
			<numDecimals>6</numDecimals>
			<responseSRS>
				<string>4326</string>
			</responseSRS>
			<overridingServiceSRS>true</overridingServiceSRS>
			<skipNumberMatched>true</skipNumberMatched>
			<circularArcPresent>true</circularArcPresent>
			<linearizationTolerance>10</linearizationTolerance>
			<attributes>
				<attribute>
					<name>the_geom</name>
					<minOccurs>0</minOccurs>
					<maxOccurs>1</maxOccurs>
					<nillable>true</nillable>
					<binding>org.locationtech.jts.geom.Point</binding>
				</attribute>
			</attributes>
		</featureType>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &FeatureType{
		XMLName: xml.Name{
			Local: "featureType",
		},
		Name:       "poi",
		NativeName: "poi",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: FeatureTypeKeywords{
			Keywords: []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"fred\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
			},
		},
		NativeCRS: FeatureTypeCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "NONE",
		Enabled:          true,
		Metadata: []*FeatureTypeMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
		Attributes: []*FeatureTypeAttribute{
			{
				Name:      "the_geom",
				MinOccurs: 0,
				MaxOccurs: 1,
				Nillable:  true,
				Binding:   "org.locationtech.jts.geom.Point",
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
		var payload *FeatureType
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &FeatureType{
			XMLName: xml.Name{
				Local: "featureType",
			},
			Name:       "poi",
			NativeName: "poi",
			Title:      "Manhattan (NY) points of interest",
			Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
			Keywords: FeatureTypeKeywords{
				Keywords: []string{
					"poi",
					"Manhattan",
					"DS_poi",
					"points_of_interest",
					"fred\\@language=ab\\;",
					"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
					"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
				},
			},
			NativeCRS: FeatureTypeCRS{
				Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
			},
			SRS: "EPSG:4326",
			NativeBoundingBox: BoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00153046439813,
				MinY: 40.70754683896324,
				MaxY: 40.719885123828675,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			LatLonBoundingBox: BoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00857344353275,
				MinY: 40.70754683896324,
				MaxY: 40.711945649065406,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			ProjectionPolicy: "NONE",
			Enabled:          true,
			Metadata: []*FeatureTypeMetadata{
				{
					Key:   "time",
					Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
				},
				{
					Key:   "cachingEnabled",
					Value: "true",
				},
			},
			Attributes: []*FeatureTypeAttribute{
				{
					Name:      "the_geom",
					MinOccurs: 0,
					MaxOccurs: 1,
					Nillable:  true,
					Binding:   "org.locationtech.jts.geom.Point",
				},
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
		XMLName: xml.Name{
			Local: "featureType",
		},
		Name:       "poi",
		NativeName: "poi",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: FeatureTypeKeywords{
			Keywords: []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"fred\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
			},
		},
		NativeCRS: FeatureTypeCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "NONE",
		Enabled:          true,
		Metadata: []*FeatureTypeMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
		Attributes: []*FeatureTypeAttribute{
			{
				Name:      "the_geom",
				MinOccurs: 0,
				MaxOccurs: 1,
				Nillable:  true,
				Binding:   "org.locationtech.jts.geom.Point",
			},
		},
	})

	assert.Nil(t, err)
}

func TestCreateFeatureTypeInDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/datastores/bar/featuretypes")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *FeatureType
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &FeatureType{
			XMLName: xml.Name{
				Local: "featureType",
			},
			Name:       "poi",
			NativeName: "poi",
			Title:      "Manhattan (NY) points of interest",
			Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
			Keywords: FeatureTypeKeywords{
				Keywords: []string{
					"poi",
					"Manhattan",
					"DS_poi",
					"points_of_interest",
					"fred\\@language=ab\\;",
					"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
					"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
				},
			},
			NativeCRS: FeatureTypeCRS{
				Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
			},
			SRS: "EPSG:4326",
			NativeBoundingBox: BoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00153046439813,
				MinY: 40.70754683896324,
				MaxY: 40.719885123828675,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			LatLonBoundingBox: BoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00857344353275,
				MinY: 40.70754683896324,
				MaxY: 40.711945649065406,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			ProjectionPolicy: "NONE",
			Enabled:          true,
			Metadata: []*FeatureTypeMetadata{
				{
					Key:   "time",
					Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
				},
				{
					Key:   "cachingEnabled",
					Value: "true",
				},
			},
			Attributes: []*FeatureTypeAttribute{
				{
					Name:      "the_geom",
					MinOccurs: 0,
					MaxOccurs: 1,
					Nillable:  true,
					Binding:   "org.locationtech.jts.geom.Point",
				},
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
		XMLName: xml.Name{
			Local: "featureType",
		},
		Name:       "poi",
		NativeName: "poi",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: FeatureTypeKeywords{
			Keywords: []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"fred\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
			},
		},
		NativeCRS: FeatureTypeCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "NONE",
		Enabled:          true,
		Metadata: []*FeatureTypeMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
		Attributes: []*FeatureTypeAttribute{
			{
				Name:      "the_geom",
				MinOccurs: 0,
				MaxOccurs: 1,
				Nillable:  true,
				Binding:   "org.locationtech.jts.geom.Point",
			},
		},
	})

	assert.Nil(t, err)
}

func TestUpdateFeatureTypeNoDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/featuretypes/toto")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *FeatureType
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &FeatureType{
			XMLName: xml.Name{
				Local: "featureType",
			},
			Name:       "poi",
			NativeName: "poi",
			Title:      "Manhattan (NY) points of interest",
			Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
			Keywords: FeatureTypeKeywords{
				Keywords: []string{
					"poi",
					"Manhattan",
					"DS_poi",
					"points_of_interest",
					"fred\\@language=ab\\;",
					"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
					"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
				},
			},
			NativeCRS: FeatureTypeCRS{
				Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
			},
			SRS: "EPSG:4326",
			NativeBoundingBox: BoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00153046439813,
				MinY: 40.70754683896324,
				MaxY: 40.719885123828675,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			LatLonBoundingBox: BoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00857344353275,
				MinY: 40.70754683896324,
				MaxY: 40.711945649065406,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			ProjectionPolicy: "NONE",
			Enabled:          true,
			Metadata: []*FeatureTypeMetadata{
				{
					Key:   "time",
					Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
				},
				{
					Key:   "cachingEnabled",
					Value: "true",
				},
			},
			Attributes: []*FeatureTypeAttribute{
				{
					Name:      "the_geom",
					MinOccurs: 0,
					MaxOccurs: 1,
					Nillable:  true,
					Binding:   "org.locationtech.jts.geom.Point",
				},
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
		XMLName: xml.Name{
			Local: "featureType",
		},
		Name:       "poi",
		NativeName: "poi",
		Title:      "Manhattan (NY) points of interest",
		Abstract:   "Points of interest in New York, New York (on Manhattan). One of the attributes contains the name of a file with a picture of the point of interest.",
		Keywords: FeatureTypeKeywords{
			Keywords: []string{
				"poi",
				"Manhattan",
				"DS_poi",
				"points_of_interest",
				"fred\\@language=ab\\;",
				"area of effect\\@language=bg\\;\\@vocabulary=Technical\\;",
				"Привет\\@language=ru\\;\\@vocabulary=Friendly\\;",
			},
		},
		NativeCRS: FeatureTypeCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: BoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "NONE",
		Enabled:          true,
		Metadata: []*FeatureTypeMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
		Attributes: []*FeatureTypeAttribute{
			{
				Name:      "the_geom",
				MinOccurs: 0,
				MaxOccurs: 1,
				Nillable:  true,
				Binding:   "org.locationtech.jts.geom.Point",
			},
		},
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
