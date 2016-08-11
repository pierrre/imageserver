package http

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

var _ http.Handler = &Handler{}

func TestHandler(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		hasETagFunc           bool
		server                imageserver.Server
		method                string
		url                   string
		header                map[string]string
		expectedStatusCode    int
		expectedHeader        map[string]string
		expectErrorFuncCalled bool
	}{
		{
			name:               "Normal",
			hasETagFunc:        true,
			url:                "http://localhost?source=medium.jpg",
			expectedStatusCode: http.StatusOK,
			expectedHeader: map[string]string{
				"Etag": fmt.Sprintf("\"%s\"", NewParamsHashETagFunc(sha256.New)(imageserver.Params{
					imageserver_source.Param: testdata.MediumFileName,
				})),
				"Content-Type":   fmt.Sprintf("image/%s", testdata.Medium.Format),
				"Content-Length": fmt.Sprint(len(testdata.Medium.Data)),
			},
		},
		{
			name:               "MethodHead",
			method:             "HEAD",
			url:                "http://localhost?source=medium.jpg",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "IfNoneMatchNoETagFunc",
			url:  "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": fmt.Sprintf("\"%s\"", NewParamsHashETagFunc(sha256.New)(imageserver.Params{
					imageserver_source.Param: testdata.MediumFileName,
				})),
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "IfNoneMatchNotModified",
			hasETagFunc: true,
			url:         "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": fmt.Sprintf("\"%s\"", NewParamsHashETagFunc(sha256.New)(imageserver.Params{
					imageserver_source.Param: testdata.MediumFileName,
				})),
			},
			expectedStatusCode: http.StatusNotModified,
		},
		{
			name:        "IfNoneMatchInvalidFormat",
			hasETagFunc: true,
			url:         "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": "foobar",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "IfNoneMatchDifferent",
			hasETagFunc: true,
			url:         "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": "\"foobar\"",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "MethodUnsupported",
			method:             "POST",
			url:                "http://localhost?source=medium.jpg",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "RequestParseError",
			url:                "http://localhost?error=foo",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "SourceParamError",
			url:  "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, &imageserver.ParamError{
					Param:   imageserver_source.Param,
					Message: "error",
				}
			}),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "ParamError",
			url:  "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, &imageserver.ParamError{
					Param:   "foobar",
					Message: "error",
				}
			}),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "ImageError",
			url:  "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, &imageserver.ImageError{
					Message: "error",
				}
			}),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "InternalError",
			url:  "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, fmt.Errorf("error")
			}),
			expectedStatusCode:    http.StatusInternalServerError,
			expectErrorFuncCalled: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			errorFuncCalled := false
			h := &Handler{
				Parser: ListParser{
					&SourceParser{},
					&testErrorParser{},
				},
				ErrorFunc: func(err error, req *http.Request) {
					errorFuncCalled = true
				},
			}
			if tc.server != nil {
				h.Server = tc.server
			} else {
				h.Server = testdata.Server
			}
			if tc.hasETagFunc {
				h.ETagFunc = NewParamsHashETagFunc(sha256.New)
			}
			rw := httptest.NewRecorder()
			met := tc.method
			if met == "" {
				met = "GET"
			}
			req, err := http.NewRequest(met, tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tc.header != nil {
				for hd, val := range tc.header {
					req.Header.Set(hd, val)
				}
			}
			h.ServeHTTP(rw, req)
			rw.Flush()
			if tc.expectedStatusCode != 0 && rw.Code != tc.expectedStatusCode {
				t.Fatalf("unexpected statud code: got %d, want %d", rw.Code, tc.expectedStatusCode)
			}
			if tc.expectedHeader != nil {
				for hd, want := range tc.expectedHeader {
					got := rw.Header().Get(hd)
					if got != want {
						t.Fatalf("unexpected value for header \"%s\": got \"%s\", want \"%s\"", hd, got, want)
					}
				}
			}
			if tc.expectErrorFuncCalled && !errorFuncCalled {
				t.Fatal("ErrorFunc not called")
			}
		})
	}
}

func TestNewParamsHashETagFunc(t *testing.T) {
	NewParamsHashETagFunc(sha256.New)(imageserver.Params{
		"foo": "bar",
	})
}
