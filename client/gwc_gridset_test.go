package client

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGridsetsSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/gridsets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<gridSets>
			<gridSet>
			<name>EPSG:3857</name>
			</gridSet>
		</gridSets>
		`))
	})
	mux.HandleFunc("/gridsets/EPSG:3857", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")

		w.WriteHeader(200)
		w.Write([]byte(`
		<gridSet>
  <name>EPSG:3857</name>
  <description></description>
  <srs>
    <number>3857</number>
  </srs>
  <extent>
    <coords>
      <double>-2.003750834E7</double>
      <double>-2.003750834E7</double>
      <double>2.003750834E7</double>
      <double>2.003750834E7</double>
    </coords>
  </extent>
  <alignTopLeft>true</alignTopLeft>
  <scaleDenominators>
    <double>5.59082264028717E8</double>
    <double>2.79541132014358E8</double>
  </scaleDenominators>
  <metersPerUnit>1.0</metersPerUnit>
  <pixelSize>2.8E-4</pixelSize>
  <scaleNames>
    <string>0</string>
    <string>1</string>
  </scaleNames>
  <tileHeight>256</tileHeight>
  <tileWidth>256</tileWidth>
  <yCoordinateFirst>false</yCoordinateFirst>
</gridSet>	  		`))
	})

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	expectedResult := []*Gridset{
		{
			XMLName: xml.Name{
				Space: "",
				Local: "gridSet",
			},
			Name:              "EPSG:3857",
			Srs:               SRS{SrsNumber: 3857},
			Description:       "",
			Extent:            []float64{-2.003750834e7, -2.003750834e7, 2.003750834e7, 2.003750834e7},
			AlignTopLeft:      true,
			ScaleDenominators: ScaleDenominators{ScaleDenominator: []float64{5.59082264028717e8, 2.79541132014358e8}},
			MetersPerUnit:     1.0,
			PixelSize:         2.8e-4,
			ScaleNames:        ScaleNames{ScaleName: []string{"0", "1"}},
			TileHeight:        256,
			TileWidth:         256,
			YCoordinateFirst:  false,
		},
	}

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	gridsets, err := cli.GetGridsets()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, gridsets)
}

func TestGetGridsetUnauthorized(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/gridsets/toto")

		w.WriteHeader(401)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	gridSet, err := cli.GetGridset("toto")

	assert.Error(t, err, "Unauthorized")
	assert.Nil(t, gridSet)
}

func TestGetGridsetNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/gridsets/toto")

		w.WriteHeader(404)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	gridSet, err := cli.GetGridset("toto")

	assert.Error(t, err, "Not Found")
	assert.Nil(t, gridSet)
}

func TestGetGridsetUnknownError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Path, "/gridsets/toto")

		w.WriteHeader(418)
		w.Write([]byte(`I'm a teapot!`))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	gridSet, err := cli.GetGridset("toto")

	assert.Error(t, err, "Unknown error: 418 - I'm a teapot!")
	assert.Nil(t, gridSet)
}

func TestCreateGridsetSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/gridsets/EPSG:3857")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *Gridset
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &Gridset{
			XMLName: xml.Name{
				Space: "",
				Local: "gridSet",
			},
			Name:              "EPSG:3857",
			Srs:               SRS{SrsNumber: 3857},
			Description:       "",
			Extent:            []float64{-2.003750834e7, -2.003750834e7, 2.003750834e7, 2.003750834e7},
			AlignTopLeft:      true,
			ScaleDenominators: ScaleDenominators{ScaleDenominator: []float64{5.59082264028717e8, 2.79541132014358e8}},
			MetersPerUnit:     1.0,
			PixelSize:         2.8e-4,
			ScaleNames:        ScaleNames{ScaleName: []string{"0", "1"}},
			TileHeight:        256,
			TileWidth:         256,
			YCoordinateFirst:  false,
		})

		w.WriteHeader(201)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.CreateGridset("EPSG:3857", &Gridset{
		XMLName: xml.Name{
			Space: "",
			Local: "gridSet",
		},
		Name:              "EPSG:3857",
		Srs:               SRS{SrsNumber: 3857},
		Description:       "",
		Extent:            []float64{-2.003750834e7, -2.003750834e7, 2.003750834e7, 2.003750834e7},
		AlignTopLeft:      true,
		ScaleDenominators: ScaleDenominators{ScaleDenominator: []float64{5.59082264028717e8, 2.79541132014358e8}},
		MetersPerUnit:     1.0,
		PixelSize:         2.8e-4,
		ScaleNames:        ScaleNames{ScaleName: []string{"0", "1"}},
		TileHeight:        256,
		TileWidth:         256,
		YCoordinateFirst:  false,
	})

	assert.Nil(t, err)
}

func TestUpdateGridsetSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "PUT")
		assert.Equal(t, r.URL.Path, "/gridsets/EPSG:3857")

		rawBody, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		var payload *Gridset
		err = xml.Unmarshal(rawBody, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload, &Gridset{
			XMLName: xml.Name{
				Space: "",
				Local: "gridSet",
			},
			Name:              "EPSG:3857",
			Srs:               SRS{SrsNumber: 3857},
			Description:       "",
			Extent:            []float64{-2.003750834e7, -2.003750834e7, 2.003750834e7, 2.003750834e7},
			AlignTopLeft:      true,
			ScaleDenominators: ScaleDenominators{ScaleDenominator: []float64{5.59082264028717e8, 2.79541132014358e8}},
			MetersPerUnit:     1.0,
			PixelSize:         2.8e-4,
			ScaleNames:        ScaleNames{ScaleName: []string{"0", "1"}},
			TileHeight:        256,
			TileWidth:         256,
			YCoordinateFirst:  false,
		})

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.UpdateGridset("EPSG:3857", &Gridset{
		XMLName: xml.Name{
			Space: "",
			Local: "gridSet",
		},
		Name:              "EPSG:3857",
		Srs:               SRS{SrsNumber: 3857},
		Description:       "",
		Extent:            []float64{-2.003750834e7, -2.003750834e7, 2.003750834e7, 2.003750834e7},
		AlignTopLeft:      true,
		ScaleDenominators: ScaleDenominators{ScaleDenominator: []float64{5.59082264028717e8, 2.79541132014358e8}},
		MetersPerUnit:     1.0,
		PixelSize:         2.8e-4,
		ScaleNames:        ScaleNames{ScaleName: []string{"0", "1"}},
		TileHeight:        256,
		TileWidth:         256,
		YCoordinateFirst:  false,
	})

	assert.Nil(t, err)
}

func TestDeleteGridsetSuccess(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE")
		assert.Equal(t, r.URL.Path, "/gridsets/EPSG:3857")

		w.WriteHeader(200)
		w.Write([]byte(``))
	}))
	defer testServer.Close()

	cli := &Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}

	err := cli.DeleteGridset("EPSG:3857")

	assert.Nil(t, err)
}
