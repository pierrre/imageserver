package http

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ http.Handler = &Handler{}

func TestHandler(t *testing.T) {
	type TC struct {
		hasETagFunc           bool
		server                imageserver.Server
		method                string
		url                   string
		header                map[string]string
		responseWriter        http.ResponseWriter
		expectedStatusCode    int
		expectedHeader        map[string]string
		expectErrorFuncCalled bool
	}
	for _, tc := range []TC{
		{
			hasETagFunc:        true,
			url:                "http://localhost?source=medium.jpg",
			expectedStatusCode: http.StatusOK,
			expectedHeader: map[string]string{
				"Etag": fmt.Sprintf("\"%s\"", NewParamsHashETagFunc(sha256.New)(imageserver.Params{
					imageserver.SourceParam: testdata.MediumFileName,
				})),
				"Content-Type":   fmt.Sprintf("image/%s", testdata.Medium.Format),
				"Content-Length": fmt.Sprint(len(testdata.Medium.Data)),
			},
		},
		{
			method:             "HEAD",
			url:                "http://localhost?source=medium.jpg",
			expectedStatusCode: http.StatusOK,
		},
		{
			url: "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": fmt.Sprintf("\"%s\"", NewParamsHashETagFunc(sha256.New)(imageserver.Params{
					imageserver.SourceParam: testdata.MediumFileName,
				})),
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			hasETagFunc: true,
			url:         "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": fmt.Sprintf("\"%s\"", NewParamsHashETagFunc(sha256.New)(imageserver.Params{
					imageserver.SourceParam: testdata.MediumFileName,
				})),
			},
			expectedStatusCode: http.StatusNotModified,
		},
		{
			hasETagFunc: true,
			url:         "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": "foobar",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			hasETagFunc: true,
			url:         "http://localhost?source=medium.jpg",
			header: map[string]string{
				"If-None-Match": "\"foobar\"",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             "POST",
			url:                "http://localhost?source=medium.jpg",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			url:                "http://localhost?error=foo",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			url: "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, &imageserver.ParamError{
					Param:   imageserver.SourceParam,
					Message: "error",
				}
			}),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			url: "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, &imageserver.ParamError{
					Param:   "foobar",
					Message: "error",
				}
			}),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			url: "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, &imageserver.ImageError{
					Message: "error",
				}
			}),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			url: "http://localhost",
			server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
				return nil, fmt.Errorf("error")
			}),
			expectedStatusCode:    http.StatusInternalServerError,
			expectErrorFuncCalled: true,
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
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
			rw := tc.responseWriter
			if rw == nil {
				rw = httptest.NewRecorder()
			}
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
			if rw, ok := rw.(*httptest.ResponseRecorder); ok {
				rw.Flush()
				if tc.expectedStatusCode != 0 && rw.Code != tc.expectedStatusCode {
					t.Fatalf("unexpected statud code: got %d, want %d", rw.Code, tc.expectedStatusCode)
				}
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
		}()
	}
}

func TestNewParamsHashETagFunc(t *testing.T) {
	NewParamsHashETagFunc(sha256.New)(imageserver.Params{
		"foo": "bar",
	})
}
