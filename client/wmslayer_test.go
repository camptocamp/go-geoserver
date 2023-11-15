package client

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWmsLayersNoWmsStoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/wmslayers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsLayers>
			<wmsLayer>
				<name>toto</name>
			</wmsLayer>
		</wmsLayers>
		`))
	})
	mux.HandleFunc("/workspaces/foo/wmslayers/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsLayer>
			<name>toto</name>
		</wmsLayer>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*WmsLayer{
		{
			XMLName: xml.Name{
				Space: "",
				Local: "wmsLayer",
			},
			Name: "toto",
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsLayers, err := cli.GetWmsLayers("foo", "")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, wmsLayers)
}

func TestGetWmsLayersInDatastoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/wmsstores/bar/wmslayers", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsLayers>
			<wmsLayer>
				<name>toto</name>
			</wmsLayer>
		</wmsLayers>
		`))
	})
	mux.HandleFunc("/workspaces/foo/wmsstores/bar/wmslayers/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsLayer>
			<name>toto</name>
		</wmsLayer>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*WmsLayer{
		{
			XMLName: xml.Name{
				Local: "wmsLayer",
			},
			Name: "toto",
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmslayers, err := cli.GetWmsLayers("foo", "bar")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, wmslayers)
}

func TestGetWmsLayerNoDatastoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/wmslayers/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsLayer>
		<name>ORTHOIMAGERY.ORTHOPHOTOS</name>
		<nativeName>ORTHOIMAGERY.ORTHOPHOTOS</nativeName>
		<namespace>
		  <name>ign</name>
		  <atom:link 
			xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="https://master.dev.scw.ansc.fr/geoserver/rest/namespaces/ign.xml" type="application/xml"/>
		  </namespace>
		  <title>Photographies aériennes</title>
		  <description>Photographies aériennes</description>
		  <abstract>Photographies aériennes</abstract>
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
		  <projectionPolicy>FORCE_DECLARED</projectionPolicy>
		  <enabled>true</enabled>
		  <metadata>
		  		<entry key="time"><dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo></entry>
		  		<entry key="cachingEnabled">true</entry>
	  		</metadata>
			</wmsLayer>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &WmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
		NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
		Title:      "Photographies aériennes",
		Abstract:   "Photographies aériennes",
		NativeCRS: WmsLayerCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "FORCE_DECLARED",
		Enabled:          true,
		Metadata: []*WmsLayerMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsLayer, err := cli.GetWmsLayer("foo", "", "toto")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, wmsLayer)
}

func TestGetWmsLayerInWmsStoreSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/wmsstores/bar/wmslayers/toto", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<wmsLayer>
		<name>ORTHOIMAGERY.ORTHOPHOTOS</name>
		<nativeName>ORTHOIMAGERY.ORTHOPHOTOS</nativeName>
		<namespace>
		  <name>ign</name>
		  <atom:link 
			xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="https://master.dev.scw.ansc.fr/geoserver/rest/namespaces/ign.xml" type="application/xml"/>
		  </namespace>
		  <title>Photographies aériennes</title>
		  <description>Photographies aériennes</description>
		  <abstract>Photographies aériennes</abstract>
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
		  <projectionPolicy>FORCE_DECLARED</projectionPolicy>
		  <enabled>true</enabled>
		  <metadata>
		  		<entry key="time"><dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo></entry>
		  		<entry key="cachingEnabled">true</entry>
	  		</metadata>
			</wmsLayer>
		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &WmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
		NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
		Title:      "Photographies aériennes",
		Abstract:   "Photographies aériennes",
		NativeCRS: WmsLayerCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "FORCE_DECLARED",
		Enabled:          true,
		Metadata: []*WmsLayerMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsLayer, err := cli.GetWmsLayer("foo", "bar", "toto")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, wmsLayer)
}

func TestGetWmsLayerUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmslayers/toto")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsLayer, err := cli.GetWmsLayer("foo", "", "toto")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, wmsLayer)
}

func TestGetWmsLayerNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmslayers/toto")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsLayer, err := cli.GetWmsLayer("foo", "", "toto")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, wmsLayer)
}

func TestGetWmsLayerUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmslayers/toto")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	wmsLayer, err := cli.GetWmsLayer("foo", "", "toto")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, wmsLayer)
}

func TestCreateWmsLayerNoDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmslayers")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *WmsLayer
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &WmsLayer{
			XMLName: xml.Name{
				Local: "wmsLayer",
			},
			Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
			NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
			Title:      "Photographies aériennes",
			Abstract:   "Photographies aériennes",
			NativeCRS: WmsLayerCRS{
				Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
			},
			SRS: "EPSG:4326",
			NativeBoundingBox: WmsLayerBoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00153046439813,
				MinY: 40.70754683896324,
				MaxY: 40.719885123828675,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			LatLonBoundingBox: WmsLayerBoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00857344353275,
				MinY: 40.70754683896324,
				MaxY: 40.711945649065406,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			ProjectionPolicy: "FORCE_DECLARED",
			Enabled:          true,
			Metadata: []*WmsLayerMetadata{
				{
					Key:   "time",
					Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
				},
				{
					Key:   "cachingEnabled",
					Value: "true",
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

	err := cli.CreateWmsLayer("foo", "", &WmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
		NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
		Title:      "Photographies aériennes",
		Abstract:   "Photographies aériennes",
		NativeCRS: WmsLayerCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "FORCE_DECLARED",
		Enabled:          true,
		Metadata: []*WmsLayerMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
	})

	assert.Nil(t, err)
}

func TestCreateWmsLayerInDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmsstores/bar/wmslayers")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *WmsLayer
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &WmsLayer{
			XMLName: xml.Name{
				Local: "wmsLayer",
			},
			Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
			NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
			Title:      "Photographies aériennes",
			Abstract:   "Photographies aériennes",
			NativeCRS: WmsLayerCRS{
				Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
			},
			SRS: "EPSG:4326",
			NativeBoundingBox: WmsLayerBoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00153046439813,
				MinY: 40.70754683896324,
				MaxY: 40.719885123828675,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			LatLonBoundingBox: WmsLayerBoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00857344353275,
				MinY: 40.70754683896324,
				MaxY: 40.711945649065406,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			ProjectionPolicy: "FORCE_DECLARED",
			Enabled:          true,
			Metadata: []*WmsLayerMetadata{
				{
					Key:   "time",
					Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
				},
				{
					Key:   "cachingEnabled",
					Value: "true",
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

	err := cli.CreateWmsLayer("foo", "bar", &WmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
		NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
		Title:      "Photographies aériennes",
		Abstract:   "Photographies aériennes",
		NativeCRS: WmsLayerCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "FORCE_DECLARED",
		Enabled:          true,
		Metadata: []*WmsLayerMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
	})

	assert.Nil(t, err)
}

func TestUpdateWmsLayerNoDatastoreSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmslayers/toto")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *WmsLayer
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &WmsLayer{
			XMLName: xml.Name{
				Local: "wmsLayer",
			},
			Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
			NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
			Title:      "Photographies aériennes",
			Abstract:   "Photographies aériennes",
			NativeCRS: WmsLayerCRS{
				Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
			},
			SRS: "EPSG:4326",
			NativeBoundingBox: WmsLayerBoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00153046439813,
				MinY: 40.70754683896324,
				MaxY: 40.719885123828675,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			LatLonBoundingBox: WmsLayerBoundingBox{
				MinX: -74.0118315772888,
				MaxX: -74.00857344353275,
				MinY: 40.70754683896324,
				MaxY: 40.711945649065406,
				CRS: FeatureTypeCRS{
					Value: "EPSG:4326",
				},
			},
			ProjectionPolicy: "FORCE_DECLARED",
			Enabled:          true,
			Metadata: []*WmsLayerMetadata{
				{
					Key:   "time",
					Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
				},
				{
					Key:   "cachingEnabled",
					Value: "true",
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

	err := cli.UpdateWmsLayer("foo", "", "toto", &WmsLayer{
		XMLName: xml.Name{
			Local: "wmsLayer",
		},
		Name:       "ORTHOIMAGERY.ORTHOPHOTOS",
		NativeName: "ORTHOIMAGERY.ORTHOPHOTOS",
		Title:      "Photographies aériennes",
		Abstract:   "Photographies aériennes",
		NativeCRS: WmsLayerCRS{
			Value: "GEOGCS[\"WGS 84\", DATUM[\"World Geodetic System 1984\", SPHEROID[\"WGS 84\", 6378137.0, 298.257223563, AUTHORITY[\"EPSG\",\"7030\"]], AUTHORITY[\"EPSG\",\"6326\"]], PRIMEM[\"Greenwich\", 0.0, AUTHORITY[\"EPSG\",\"8901\"]], UNIT[\"degree\", 0.017453292519943295], AXIS[\"Geodetic longitude\", EAST], AXIS[\"Geodetic latitude\", NORTH], AUTHORITY[\"EPSG\",\"4326\"]]",
		},
		SRS: "EPSG:4326",
		NativeBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00153046439813,
			MinY: 40.70754683896324,
			MaxY: 40.719885123828675,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		LatLonBoundingBox: WmsLayerBoundingBox{
			MinX: -74.0118315772888,
			MaxX: -74.00857344353275,
			MinY: 40.70754683896324,
			MaxY: 40.711945649065406,
			CRS: FeatureTypeCRS{
				Value: "EPSG:4326",
			},
		},
		ProjectionPolicy: "FORCE_DECLARED",
		Enabled:          true,
		Metadata: []*WmsLayerMetadata{
			{
				Key:   "time",
				Value: "<dimensionInfo><enabled>false</enabled><defaultValue/></dimensionInfo>",
			},
			{
				Key:   "cachingEnabled",
				Value: "true",
			},
		},
	})

	assert.Nil(t, err)
}

func TestDeleteWmsLayerSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/wmslayers/toto")
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

	err := cli.DeleteWmsLayer("foo", "", "toto", true)

	assert.Nil(t, err)
}
