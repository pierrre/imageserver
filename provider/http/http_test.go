package http

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var (
	testSourceFileName = testdata.SmallFileName
)

func TestInterfaceProvider(t *testing.T) {
	var _ imageserver.Provider = &HTTPProvider{}
}

func TestGet(t *testing.T) {
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := createTestURL(listener)
	parameters := make(imageserver.Parameters)

	image, err := provider.Get(source, parameters)
	if err != nil {
		t.Fatal(err)
	}

	if len(image.Data) == 0 {
		t.Fatal("no data")
	}

	if len(image.Format) == 0 {
		t.Fatal("no format")
	}
}

func TestGetErrorNotFound(t *testing.T) {
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := createTestURL(listener)
	source.Path += "foobar"
	parameters := make(imageserver.Parameters)

	_, err = provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrl(t *testing.T) {
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := "foobar"
	parameters := make(imageserver.Parameters)

	_, err = provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrlScheme(t *testing.T) {
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := "custom://foobar"
	parameters := make(imageserver.Parameters)

	_, err = provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidHost(t *testing.T) {
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := "http://invalid.foobar.com"
	parameters := make(imageserver.Parameters)

	_, err = provider.Get(source, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseFormatEmpty(t *testing.T) {
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := createTestURL(listener)

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
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := createTestURL(listener)

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
	provider, listener, err := createTestHTTPProvider()
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	source := createTestURL(listener)

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

func createTestHTTPProvider() (provider *HTTPProvider, listener *net.TCPListener, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "")
	if err != nil {
		return
	}

	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}

	server := &http.Server{
		Handler: http.FileServer(http.Dir(testdata.Dir)),
	}
	go server.Serve(listener)

	provider = &HTTPProvider{}

	return
}

func createTestURL(listener *net.TCPListener) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   listener.Addr().String(),
		Path:   testSourceFileName,
	}
}
