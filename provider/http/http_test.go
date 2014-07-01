package http

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_provider "github.com/pierrre/imageserver/provider"
	"github.com/pierrre/imageserver/testdata"
)

var (
	testSourceFileName = testdata.SmallFileName
)

func TestInterface(t *testing.T) {
	var _ imageserver_provider.Provider = &HTTPProvider{}
}

func TestGet(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := createTestURL(listener)
	parameters := make(imageserver.Parameters)

	provider := &HTTPProvider{}

	image, err := provider.Get(source, parameters)
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("no image")
	}
	if len(image.Data) == 0 {
		t.Fatal("no data")
	}
	if len(image.Format) == 0 {
		t.Fatal("no format")
	}
}

func TestGetErrorNotFound(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := createTestURL(listener)
	source.Path += "foobar"
	parameters := make(imageserver.Parameters)

	provider := &HTTPProvider{}

	_, err := provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrl(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := "foobar"
	parameters := make(imageserver.Parameters)

	provider := &HTTPProvider{}

	_, err := provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrlScheme(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := "custom://foobar"
	parameters := make(imageserver.Parameters)

	provider := &HTTPProvider{}

	_, err := provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidHost(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := "http://invalid.localhost"
	parameters := make(imageserver.Parameters)

	provider := &HTTPProvider{}

	_, err := provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseFormatEmpty(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := createTestURL(listener)

	provider := &HTTPProvider{}

	response, err := provider.getResponse(source)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	response.Header.Del("Content-Type")

	image, err := provider.createImage(response)
	if err != nil {
		t.Fatal(err)
	}

	if len(image.Format) != 0 {
		t.Fatal("format not empty")
	}
}

func TestParseFormatInvalid(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := createTestURL(listener)

	provider := &HTTPProvider{}

	response, err := provider.getResponse(source)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	response.Header.Set("Content-Type", "foobar")

	image, err := provider.createImage(response)
	if err != nil {
		t.Fatal(err)
	}

	if len(image.Format) != 0 {
		t.Fatal("format not empty")
	}
}

func TestParseDataError(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()

	source := createTestURL(listener)

	provider := &HTTPProvider{}

	response, err := provider.getResponse(source)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()

	_, err = provider.createImage(response)
	if err == nil {
		t.Fatal("no error")
	}
}

func createTestHTTPServer(t *testing.T) *net.TCPListener {
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		t.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	server := &http.Server{
		Handler: http.FileServer(http.Dir(testdata.Dir)),
	}
	go server.Serve(listener)

	return listener
}

func createTestURL(listener *net.TCPListener) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   listener.Addr().String(),
		Path:   testSourceFileName,
	}
}
