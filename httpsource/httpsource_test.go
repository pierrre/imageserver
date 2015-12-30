package httpsource

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestGet(t *testing.T) {
	httpSrv := createTestHTTPServer(t)
	defer httpSrv.Close()
	params := imageserver.Params{
		imageserver.SourceParam: createTestSource(httpSrv, testdata.MediumFileName),
	}
	srv := &Server{}
	im, err := srv.Get(params)
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
	httpSrv := createTestHTTPServer(t)
	defer httpSrv.Close()
	params := imageserver.Params{}
	srv := &Server{}
	_, err := srv.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorNotFound(t *testing.T) {
	httpSrv := createTestHTTPServer(t)
	defer httpSrv.Close()
	source := createTestSource(httpSrv, testdata.MediumFileName)
	source += "foobar"
	params := imageserver.Params{imageserver.SourceParam: source}
	srv := &Server{}
	_, err := srv.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorInvalidUrl(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "foobar"}
	srv := &Server{}
	_, err := srv.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorInvalidUrlScheme(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "custom://foobar"}
	srv := &Server{}
	_, err := srv.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestGetErrorRequest(t *testing.T) {
	params := imageserver.Params{imageserver.SourceParam: "http://localhost:123456"}
	srv := &Server{}
	_, err := srv.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

type errorReadCloser struct{}

func (erc *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error")
}

func (erc *errorReadCloser) Close() error {
	return fmt.Errorf("error")
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

func createTestHTTPServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.FileServer(http.Dir(testdata.Dir)))
}

func createTestSource(srv *httptest.Server, filename string) string {
	return fmt.Sprintf("http://%s/%s", srv.Listener.Addr(), filename)
}
