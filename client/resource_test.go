package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResourceSuccess(t *testing.T) {
	const resourceContent = `
	<styles>
		<style>
		<name>line</name>
		</style>
	</styles>
	`
	mux := http.NewServeMux()
	mux.HandleFunc("/resource/test_dir/style.xml", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(resourceContent))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}
	resourceDef, err := cli.GetResource("test_dir/style", "xml")

	assert.Nil(t, err)
	assert.Equal(t, resourceContent, resourceDef)
}

func TestGetResourceUnknown(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/resource/style.sld")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	resource, err := cli.GetResource("style", "sld")

	assert.Error(t, err, "Resource not found")
	assert.Equal(t, string(resource), "")
}

func TestCreateResource(t *testing.T) {
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
		assert.Equal(t, r.URL.Path, "/resource/test_dir/test.sld")

		rawBody, err := io.ReadAll(r.Body)
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

	err := cli.CreateResource("test_dir/test", "sld", styleDefinition)

	assert.Nil(t, err)
}

func TestUpdateResourceSuccess(t *testing.T) {
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
		assert.Equal(t, r.URL.Path, "/resource/test_dir/test.sld")

		rawBody, err := io.ReadAll(r.Body)
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

	err := cli.UpdateResource("test_dir/test", "sld", styleDefinition)

	assert.Nil(t, err)
}

func TestDeleteResource(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/resource/test_dir/test.sld")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteResource("test_dir/test.sld")

	assert.Nil(t, err)
}
