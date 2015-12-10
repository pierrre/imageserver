package main

import (
	"image"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver/testdata"
)

func TestServer(t *testing.T) {
	h := newHTTPHandler()
	type TC struct {
		path               string
		query              url.Values
		expectedStatusCode int
		expectedFormat     string
		expectedWidth      int
		expectedHeight     int
	}
	for _, tc := range []TC{
		{
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			path: testdata.SmallFileName,
		},
		{
			path: testdata.MediumFileName,
		},
		{
			path: testdata.LargeFileName,
		},
		{
			path: testdata.HugeFileName,
		},
		{
			path: testdata.AnimatedFileName,
		},
		{
			path: testdata.MediumFileName,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"format": {"foobar"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"format": {"png"},
			},
			expectedFormat: "png",
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"format": {"gif"},
			},
			expectedFormat: "gif",
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"format":  {"jpeg"},
				"quality": {"-10"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"format":  {"jpeg"},
				"quality": {"50"},
			},
			expectedFormat: "jpeg",
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"width": {"-100"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"width": {"9001"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"width": {"100"},
			},
			expectedWidth: 100,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"height": {"-100"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"height": {"9001"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			path: testdata.MediumFileName,
			query: url.Values{
				"height": {"100"},
			},
			expectedHeight: 100,
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			u := &url.URL{
				Scheme:   "http",
				Host:     "localhost",
				Path:     "/" + tc.path,
				RawQuery: tc.query.Encode(),
			}
			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			w.Flush()
			if tc.expectedStatusCode != 0 && w.Code != tc.expectedStatusCode {
				t.Fatalf("unexpected http status: %d", w.Code)
			}
			if w.Code != http.StatusOK {
				if tc.expectedStatusCode != 0 {
					return
				}
				t.Fatalf("http status not OK: %d", w.Code)
			}
			im, format, err := image.Decode(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tc.expectedFormat != "" && format != tc.expectedFormat {
				t.Fatalf("unexpected format: %s", format)
			}
			if tc.expectedWidth != 0 && im.Bounds().Dx() != tc.expectedWidth {
				t.Fatalf("unexpected width: %d", im.Bounds().Dx())
			}
			if tc.expectedHeight != 0 && im.Bounds().Dy() != tc.expectedHeight {
				t.Fatalf("unexpected height: %d", im.Bounds().Dy())
			}
		}()
	}
}
