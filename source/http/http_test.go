package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestServerGet(t *testing.T) {
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
			name:               "ErrorNoSource",
			params:             imageserver.Params{},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "ErrorInvalidURL",
			params: imageserver.Params{
				imageserver_source.Param: "%",
			},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "ErrorUnreachableURL",
			params: imageserver.Params{
				imageserver_source.Param: "http://localhost:123456",
			},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "ErrorNotFound",
			params: imageserver.Params{
				imageserver_source.Param: createTestSource(httpSrv, testdata.MediumFileName) + "foobar",
			},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "ErrorIdentify",
			params: imageserver.Params{
				imageserver_source.Param: createTestSource(httpSrv, "testdata.go"),
			},
			expectedParamError: imageserver_source.Param,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			im, err := srv.Get(context.Background(), tc.params)
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
		})
	}
}

type errorReadCloser struct{}

func (erc *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("error")
}

func (erc *errorReadCloser) Close() error {
	return fmt.Errorf("error")
}

func TestLoadDataError(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       &errorReadCloser{},
	}
	_, err := loadData(resp)
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
	return srv.URL + "/" + filename
}

func TestIdentifyHeader(t *testing.T) {
	for _, tc := range []struct {
		name           string
		resp           *http.Response
		data           []byte
		expectedFormat string
		expectedError  bool
	}{
		{
			name: "JPEG",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": {"image/" + testdata.Medium.Format},
				},
			},
			data:           testdata.Medium.Data,
			expectedFormat: testdata.Medium.Format,
			expectedError:  false,
		},
		{
			name: "PNG",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": {"image/" + testdata.Random.Format},
				},
			},
			data:           testdata.Random.Data,
			expectedFormat: testdata.Random.Format,
			expectedError:  false,
		},
		{
			name: "GIF",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": {"image/" + testdata.Animated.Format},
				},
			},
			data:           testdata.Animated.Data,
			expectedFormat: testdata.Animated.Format,
			expectedError:  false,
		},
		{
			name: "ErrorNoHeader",
			resp: &http.Response{
				StatusCode: http.StatusOK,
			},
			data:          testdata.Medium.Data,
			expectedError: true,
		},
		{
			name: "InvalidHeader",
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": {"invalid"},
				},
			},
			data:          testdata.Medium.Data,
			expectedError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			format, err := IdentifyHeader(tc.resp, tc.data)
			if err != nil {
				if tc.expectedError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedError {
				t.Fatal("no error")
			}
			if format != tc.expectedFormat {
				t.Fatalf("unexpected format: got %s, want %s", format, tc.expectedFormat)
			}
		})
	}
}
