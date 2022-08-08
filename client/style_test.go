package client

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStylesNoWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/styles", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<styles>
			<style>
			<name>line</name>
			</style>
		</styles>
		`))
	})
	mux.HandleFunc("/styles/line", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<style>
		<name>line</name>
		<format>sld</format>
		<languageVersion>
		  <version>1.0.0</version>
		</languageVersion>
		<filename>default_line.sld</filename>
	  </style>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*Style{
		&Style{
			XMLName: xml.Name{
				Space: "",
				Local: "style",
			},
			Name:     "line",
			Format:   "sld",
			Version:  &LanguageVersion{Version: "1.0.0"},
			FileName: "default_line.sld",
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	styles, err := cli.GetStyles("")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, styles)
}

func TestGetStylesWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/styles", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<styles>
			<style>
			<name>line</name>
			</style>
		</styles>
		`))
	})
	mux.HandleFunc("/workspaces/foo/styles/line", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<style>
		<name>line</name>
		<format>sld</format>
		<workspace>
    <name>foo</name>
  </workspace>
		<languageVersion>
		  <version>1.0.0</version>
		</languageVersion>
		<filename>default_line.sld</filename>
	  </style>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*Style{
		&Style{
			XMLName: xml.Name{
				Space: "",
				Local: "style",
			},
			Workspace: &WorkspaceRef{Name: "foo"},
			Name:      "line",
			Format:    "sld",
			Version:   &LanguageVersion{Version: "1.0.0"},
			FileName:  "default_line.sld",
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	styles, err := cli.GetStyles("foo")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, styles)
}

func TestGetStyleNoWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/styles/line", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<style>
		<name>line</name>
		<format>sld</format>
		<languageVersion>
		  <version>1.0.0</version>
		</languageVersion>
		<filename>default_line.sld</filename>
	  </style>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &Style{
		XMLName: xml.Name{
			Space: "",
			Local: "style",
		},
		Name:     "line",
		Format:   "sld",
		Version:  &LanguageVersion{Version: "1.0.0"},
		FileName: "default_line.sld",
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	style, err := cli.GetStyle("", "line")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, style)
}

func TestGetStyleWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/workspaces/foo/styles/line", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<style>
		<name>line</name>
		<format>sld</format>
		<workspace>
    <name>foo</name>
  </workspace>
		<languageVersion>
		  <version>1.0.0</version>
		</languageVersion>
		<filename>default_line.sld</filename>
	  </style>
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := &Style{
		XMLName: xml.Name{
			Space: "",
			Local: "style",
		},
		Workspace: &WorkspaceRef{Name: "foo"},
		Name:      "line",
		Format:    "sld",
		Version:   &LanguageVersion{Version: "1.0.0"},
		FileName:  "default_line.sld",
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	style, err := cli.GetStyle("foo", "line")

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, style)
}

func TestGetStyleUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/styles/toto")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	style, err := cli.GetStyle("", "toto")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, style)
}

func TestGetStyleNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/styles/toto")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	style, err := cli.GetStyle("", "toto")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, style)
}

func TestGetStyleUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/styles/toto")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	style, err := cli.GetStyle("", "toto")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, style)
}

func TestGetStyleSLDNoWorkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/styles/point", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<?xml version="1.0" encoding="ISO-8859-1"?>
		<StyledLayerDescriptor version="1.0.0"
				xsi:schemaLocation="http://www.opengis.net/sld StyledLayerDescriptor.xsd"
				xmlns="http://www.opengis.net/sld"
				xmlns:ogc="http://www.opengis.net/ogc"
				xmlns:xlink="http://www.w3.org/1999/xlink"
				xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
				<!-- a named layer is the basic building block of an sld document -->
		
			<NamedLayer>
				<Name>Default Point</Name>
				<UserStyle>
					<!-- they have names, titles and abstracts -->
		
					<Title>Red Square point</Title>
					<Abstract>A sample style that just prints out a red square</Abstract>
					<!-- FeatureTypeStyles describe how to render different features -->
					<!-- a feature type for points -->
		
					<FeatureTypeStyle>
						<!--FeatureTypeName>Feature</FeatureTypeName-->
						<Rule>
							<Name>Rule 1</Name>
							<Title>Red Square point</Title>
							<Abstract>A red fill with 6 pixels size</Abstract>
		
							<!-- like a linesymbolizer but with a fill too -->
							<PointSymbolizer>
								<Graphic>
									<Mark>
										<WellKnownName>square</WellKnownName>
										<Fill>
											<CssParameter name="fill">#FF0000</CssParameter>
										</Fill>
									</Mark>
									<Size>6</Size>
								</Graphic>
							</PointSymbolizer>
						</Rule>
		
					</FeatureTypeStyle>
				</UserStyle>
			</NamedLayer>
		</StyledLayerDescriptor>
		
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	styleFile, err := cli.GetStyleFile("", "point", "sld", "1.0.0")

	assert.Nil(t, err)
	assert.NotEmpty(t, styleFile)
}

func TestGetStyleCssNoWOrkspaceSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/styles/point", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		/* @title teal point */
		* {
			mark: symbol(square);
			mark-size: 6px;
			:mark {
				fill: #00cc33;
			}
		}		
	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	styleFile, err := cli.GetStyleFile("", "point", "css", "1.0.0")

	assert.Nil(t, err)
	assert.NotEmpty(t, styleFile)
}

func TestCreateStyleNoWorkspaceSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/styles")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *Style
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &Style{
			XMLName: xml.Name{
				Local: "style",
			},
			Name:     "test_style",
			FileName: "test_style.sld",
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateStyle("", &Style{
		XMLName: xml.Name{
			Local: "style",
		},
		Name:     "test_style",
		FileName: "test_style.sld",
	})

	assert.Nil(t, err)
}

func TestUpdateStyleContentSldSuccess(t *testing.T) {
	const styleDefinition = `
	<StyledLayerDescriptor version="1.0.0"
	  xsi:schemaLocation="http://www.opengis.net/sld http://schemas.opengis.net/sld/1.0.0/StyledLayerDescriptor.xsd"
	  xmlns="http://www.opengis.net/sld" xmlns:ogc="http://www.opengis.net/ogc"
	  xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	
	  <NamedLayer>
		<Name>go_test</Name>
		<UserStyle>
		  <Title>A teal polygon style</Title>
		  <FeatureTypeStyle>
			<Rule>
			  <Title>teal polygon</Title>
			  <PolygonSymbolizer>
				<Fill>
				  <CssParameter name="fill">#00cc33
				  </CssParameter>
				</Fill>
				<Stroke>
				  <CssParameter name="stroke">#000000</CssParameter>
				  <CssParameter name="stroke-width">0.5</CssParameter>
				</Stroke>
			  </PolygonSymbolizer>
	
			</Rule>
	
		  </FeatureTypeStyle>
		</UserStyle>
	  </NamedLayer>
	</StyledLayerDescriptor>
	
	`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/styles/toto")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/vnd.ogc.sld+xml")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(rawBody), styleDefinition)

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	styleToCreate := &Style{
		XMLName: xml.Name{
			Local: "style",
		},
		Name:     "toto",
		FileName: "test_style.sld",
		Format:   "sld",
		Version:  &LanguageVersion{Version: "1.0.0"},
	}

	err := cli.UpdateStyleContent("", styleToCreate, styleDefinition)

	assert.Nil(t, err)
}

func TestUpdateStyleContentCssSuccess(t *testing.T) {
	const styleDefinition = `
	/* @title cyan line */
	* {
		stroke: #0099cc;
	}		
	`
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/styles/toto")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/vnd.geoserver.geocss+css")

		rawBody, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(rawBody), styleDefinition)

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	styleToCreate := &Style{
		XMLName: xml.Name{
			Local: "style",
		},
		Name:     "toto",
		FileName: "test_style.css",
		Format:   "css",
		Version:  &LanguageVersion{Version: "1.0.0"},
	}

	err := cli.UpdateStyleContent("", styleToCreate, styleDefinition)

	assert.Nil(t, err)
}

func TestDeleteStyleNoWorkspace(t *testing.T) {
	cli := &Client{
		URL:        "http://localhost:8080/geoserver/rest",
		Username:   "admin",
		Password:   "geoserver",
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteStyle("", "test_style", true, false)
	assert.Nil(t, err)
}

func TestDeleteStyleSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/styles/toto")
		keys, ok := r.URL.Query()["recurse"]
		assert.True(t, ok)
		assert.Equal(t, keys[0], "true")
		keys2, ok2 := r.URL.Query()["purge"]
		assert.True(t, ok2)
		assert.Equal(t, keys2[0], "true")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteStyle("", "toto", true, true)

	assert.Nil(t, err)
}
