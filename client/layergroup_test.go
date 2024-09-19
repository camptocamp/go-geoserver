package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGroupsNoWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/layergroups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layerGroups>
			<layerGroup>
			<name>osm_group</name>
			</layerGroup>
		</layerGroups>
		`))
	})
	mux.HandleFunc("/layergroups/osm_group", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layerGroup>
		<name>osm_group</name>
		<mode>SINGLE</mode>
		<title>toto</title>
		<abstractTxt>ABSTRACT</abstractTxt>
		<publishables>
		  <published type="layer">
			<name>osm:simplified_water_polygons</name>
			<atom:link 
			  xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/layers/simplified_water_polygons.xml" type="application/xml"/>
			</published>
			<published type="layer">
			  <name>osm:water_polygons</name>
			  <atom:link 
				xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/layers/water_polygons.xml" type="application/xml"/>
			  </published>
			</publishables>
			<styles>
			  <style>
				<name>osm:simplified_water</name>
				<atom:link 
				  xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/styles/simplified_water.xml" type="application/xml"/>
				</style>
				<style>
				  <name>osm:water</name>
				  <atom:link 
					xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/styles/water.xml" type="application/xml"/>
				  </style>
				</styles>
				<metadataLinks>
				  <metadataLink>
					<type>text/plain</type>
					<metadataType>ISO19115:2003</metadataType>
					<content>https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa</content>
				  </metadataLink>
				  <metadataLink>
					<type>text/plain</type>
					<metadataType>FGDC</metadataType>
					<content>https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa</content>
				  </metadataLink>
				</metadataLinks>
				<bounds>
				  <minx>-2</minx>
				  <maxx>2</maxx>
				  <miny>-2</miny>
				  <maxy>2</maxy>
				  <crs class="projected">EPSG:3857</crs>
				</bounds>
				<keywords>
				  <string>TestKEwyow</string>
				</keywords>
			  </layerGroup>	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*LayerGroup{
		{
			XMLName: xml.Name{
				Space: "",
				Local: "layerGroup",
			},
			Name:     "osm_group",
			Mode:     "SINGLE",
			Title:    "toto",
			Abstract: "ABSTRACT",
			Publishables: []*LayerRef{
				{
					Type: "layer",
					Name: "osm:simplified_water_polygons",
				},
				{
					Type: "layer",
					Name: "osm:water_polygons",
				},
			},
			Styles: []*StyleRef{
				{
					Name: "osm:simplified_water",
				},
				{
					Name: "osm:water",
				},
			},
			Bounds: &BoundingBox{
				MinX: -2,
				MaxX: 2,
				MinY: -2,
				MaxY: 2,
				CRS: FeatureTypeCRS{
					Class: "projected",
					Value: "EPSG:3857",
				},
			},
			MetadataLinks: []*MetadataLink{
				{
					Type:         "text/plain",
					MetadataType: "ISO19115:2003",
					Content:      "https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa",
				},
				{
					Type:         "text/plain",
					MetadataType: "FGDC",
					Content:      "https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa",
				},
			},
			Keywords: LayerGroupKeywords{
				Keywords: []string{"TestKEwyow"},
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerGroups, err := cli.GetGroups("")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layerGroups)
}

func TestGetGroupsWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/osm/layergroups", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layerGroups>
			<layerGroup>
			<name>osm_group</name>
			</layerGroup>
		</layerGroups>
		`))
	})
	mux.HandleFunc("/workspaces/osm/layergroups/osm_group", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<layerGroup>
		<name>osm_group</name>
		<mode>SINGLE</mode>
		<title>toto</title>
		<abstractTxt>ABSTRACT</abstractTxt>
		<workspace>
    <name>osm</name>
  </workspace>
		<publishables>
		  <published type="layer">
			<name>osm:simplified_water_polygons</name>
			<atom:link 
			  xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/layers/simplified_water_polygons.xml" type="application/xml"/>
			</published>
			<published type="layer">
			  <name>osm:water_polygons</name>
			  <atom:link 
				xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/layers/water_polygons.xml" type="application/xml"/>
			  </published>
			</publishables>
			<styles>
			  <style>
				<name>osm:simplified_water</name>
				<atom:link 
				  xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/styles/simplified_water.xml" type="application/xml"/>
				</style>
				<style>
				  <name>osm:water</name>
				  <atom:link 
					xmlns:atom="http://www.w3.org/2005/Atom" rel="alternate" href="http://localhost:8080/geoserver/rest/workspaces/osm/styles/water.xml" type="application/xml"/>
				  </style>
				</styles>
				<metadataLinks>
				  <metadataLink>
					<type>text/plain</type>
					<metadataType>ISO19115:2003</metadataType>
					<content>https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa</content>
				  </metadataLink>
				  <metadataLink>
					<type>text/plain</type>
					<metadataType>FGDC</metadataType>
					<content>https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa</content>
				  </metadataLink>
				</metadataLinks>
				<bounds>
				  <minx>-2</minx>
				  <maxx>2</maxx>
				  <miny>-2</miny>
				  <maxy>2</maxy>
				  <crs class="projected">EPSG:3857</crs>
				</bounds>
				<keywords>
				  <string>TestKEwyow</string>
				</keywords>
			  </layerGroup>	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*LayerGroup{
		{
			XMLName: xml.Name{
				Space: "",
				Local: "layerGroup",
			},
			Name: "osm_group",
			Workspace: &WorkspaceRef{
				Name: "osm",
			},
			Mode:     "SINGLE",
			Title:    "toto",
			Abstract: "ABSTRACT",
			Publishables: []*LayerRef{
				{
					Type: "layer",
					Name: "osm:simplified_water_polygons",
				},
				{
					Type: "layer",
					Name: "osm:water_polygons",
				},
			},
			Styles: []*StyleRef{
				{
					Name: "osm:simplified_water",
				},
				{
					Name: "osm:water",
				},
			},
			Bounds: &BoundingBox{
				MinX: -2,
				MaxX: 2,
				MinY: -2,
				MaxY: 2,
				CRS: FeatureTypeCRS{
					Class: "projected",
					Value: "EPSG:3857",
				},
			},
			MetadataLinks: []*MetadataLink{
				{
					Type:         "text/plain",
					MetadataType: "ISO19115:2003",
					Content:      "https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa",
				},
				{
					Type:         "text/plain",
					MetadataType: "FGDC",
					Content:      "https://jira.nexsis18-112.fr/jira/secure/Dashboard.jspa",
				},
			},
			Keywords: LayerGroupKeywords{
				Keywords: []string{"TestKEwyow"},
			},
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerGroups, err := cli.GetGroups("osm")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, layerGroups)
}

func TestGetGroupUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layergroups/toto")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerGroup, err := cli.GetGroup("", "toto")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, layerGroup)
}

func TestGetGroupNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layergroups/toto")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerGroup, err := cli.GetGroup("", "toto")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, layerGroup)
}

func TestGetGroupUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/layergroups/toto")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	layerGroup, err := cli.GetGroup("", "toto")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, layerGroup)
}

func TestCreateLayerGroupNoWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/layergroups")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *LayerGroup
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &LayerGroup{
			XMLName: xml.Name{Local: "layerGroup"},
			Name:    "test_style",
			Publishables: []*LayerRef{
				{
					Type: "layer",
					Name: "layer1",
				},
				{
					Type: "layer",
					Name: "layer2",
				},
			},
			Styles: []*StyleRef{
				{
					Name: "style1",
				},
				{
					Name: "",
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

	err := cli.CreateGroup("", &LayerGroup{
		XMLName: xml.Name{Local: "layerGroup"},
		Name:    "test_style",
		Publishables: []*LayerRef{
			{
				Type: "layer",
				Name: "layer1",
			},
			{
				Type: "layer",
				Name: "layer2",
			},
		},
		Styles: []*StyleRef{
			{
				Name: "style1",
			},
			{
				Name: "",
			},
		},
	})

	assert.Nil(t, err)
}

func TestUpdateLayerGroupNoWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/layergroups/test_style")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *LayerGroup
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &LayerGroup{
			XMLName: xml.Name{Local: "layerGroup"},
			Name:    "test_style",
			Publishables: []*LayerRef{
				{
					Type: "layer",
					Name: "layer1",
				},
				{
					Type: "layer",
					Name: "layer2",
				},
			},
			Styles: []*StyleRef{
				{
					Name: "style1",
				},
				{
					Name: "style2",
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

	err := cli.UpdateGroup("", &LayerGroup{
		XMLName: xml.Name{Local: "layerGroup"},
		Name:    "test_style",
		Publishables: []*LayerRef{
			{
				Type: "layer",
				Name: "layer1",
			},
			{
				Type: "layer",
				Name: "layer2",
			},
		},
		Styles: []*StyleRef{
			{
				Name: "style1",
			},
			{
				Name: "style2",
			},
		},
	})

	assert.Nil(t, err)
}

func TestDeleteLayerGroupSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/workspaces/foo/layergroups/toto")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteGroup("foo", "toto")

	assert.Nil(t, err)
}
