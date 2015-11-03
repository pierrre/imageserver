package httpsource

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestGet(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()
	params := imageserver.Params{imageserver.SourceParam: createTestSource(listener, testdata.MediumFileName)}
	server := &Server{}
	im, err := server.Get(params)
	if err != nil {
		t.Fatal(err)
	}
	if im == nil {
		t.Fatal("no image")
	}
	if im.Format != testdata.Medium.Format {
		t.Fatalf("unexpected image format: got \"%s\", want \"%s\"", im.Format, testdata.Medium.Format)
	}
	if len(im.Data) != len(testdata.Medium.Data) {
		t.Fatalf("unexpected image data length: got %d, want %d", len(im.Data), len(testdata.Medium.Data))
	}
}

func TestGetErrorNoSource(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()
	params := imageserver.Params{}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorNotFound(t *testing.T) {
	listener := createTestHTTPServer(t)
	defer listener.Close()
	source := createTestSource(listener, testdata.MediumFileName)
	source += "foobar"
	params := imageserver.Params{imageserver.SourceParam: source}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorInvalidUrl(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "foobar"}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorInvalidUrlScheme(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "custom://foobar"}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorRequest(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "http://localhost:123456"}
	server := &Server{}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

type errorReadCloser struct{}

func (erc *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func (erc *errorReadCloser) Close() error {
	return errors.New("error")
}

func TestParseResponseErrorData(t *testing.T) {
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       &errorReadCloser{},
	}
	_, err := parseResponse(response)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
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

func createTestSource(listener *net.TCPListener, filename string) string {
	return fmt.Sprintf("http://%s/%s", listener.Addr(), filename)
}
