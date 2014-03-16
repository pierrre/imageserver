package http

import (
	"testing"

	"github.com/pierrre/imageserver"
)

const (
	testSource = "https://raw.github.com/pierrre/imageserver/master/testdata/small.jpg"
)

func TestGet(t *testing.T) {
	provider := createTestHTTPProvider()
	parameters := make(imageserver.Parameters)

	image, err := provider.Get(testSource, parameters)
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
	provider := createTestHTTPProvider()
	parameters := make(imageserver.Parameters)

	_, err := provider.Get(testSource+"foobar", parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrl(t *testing.T) {
	provider := createTestHTTPProvider()
	parameters := make(imageserver.Parameters)

	_, err := provider.Get("foobar", parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidUrlScheme(t *testing.T) {
	provider := createTestHTTPProvider()
	parameters := make(imageserver.Parameters)

	_, err := provider.Get("custom://foobar", parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorInvalidHost(t *testing.T) {
	provider := createTestHTTPProvider()
	parameters := make(imageserver.Parameters)

	_, err := provider.Get("http://invalid.foobar.com", parameters)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseFormatEmpty(t *testing.T) {
	provider := createTestHTTPProvider()

	response, err := provider.getResponse(testSource)
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
	provider := createTestHTTPProvider()

	response, err := provider.getResponse(testSource)
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
	provider := createTestHTTPProvider()

	response, err := provider.getResponse(testSource)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()

	_, err = provider.createImage(response)
	if err == nil {
		t.Fatal("no error")
	}
}

func createTestHTTPProvider() *HTTPProvider {
	return &HTTPProvider{}
}
