package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestGet(t *testing.T) {
	srv := &Server{}
	httpSrv := createTestHTTPServer()
	defer httpSrv.Close()
	for _, tc := range []struct {
		name               string
		params             imageserver.Params
		expectedParamError string
		expectedImage      *imageserver.Image
	}{
		{
			name: "Normal",
			params: imageserver.Params{
				imageserver_source.Param: createTestSource(httpSrv, testdata.MediumFileName),
			},
			expectedImage: testdata.Medium,
		},
		{
			name:               "NoSource",
			params:             imageserver.Params{},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "InvalidURL",
			params: imageserver.Params{
				imageserver_source.Param: "%",
			},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "UnreachableURL",
			params: imageserver.Params{
				imageserver_source.Param: "http://localhost:123456",
			},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "NotFound",
			params: imageserver.Params{
				imageserver_source.Param: createTestSource(httpSrv, testdata.MediumFileName) + "foobar",
			},
			expectedParamError: imageserver_source.Param,
		},
	} {
		func() {
			t.Logf("test: %s", tc.name)
			im, err := srv.Get(tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && err.Param == tc.expectedParamError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedParamError != "" {
				t.Fatal("no error")
			}
			if im == nil {
				t.Fatal("no image")
			}
			if im.Format != tc.expectedImage.Format {
				t.Fatalf("unexpected image format: got \"%s\", want \"%s\"", im.Format, tc.expectedImage.Format)
			}
			if !bytes.Equal(im.Data, tc.expectedImage.Data) {
				t.Fatal("data not equal")
			}
		}()
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

func createTestHTTPServer() *httptest.Server {
	return httptest.NewServer(http.FileServer(http.Dir(testdata.Dir)))
}

func createTestSource(srv *httptest.Server, filename string) string {
	return fmt.Sprintf("http://%s/%s", srv.Listener.Addr(), filename)
}
