package main

import (
	"image"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestServer(t *testing.T) {
	h := newImageHTTPHandler()
	type TC struct {
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
			query: url.Values{
				imageserver.SourceParam: {testdata.SmallFileName},
			},
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
			},
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.LargeFileName},
			},
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.HugeFileName},
			},
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.AnimatedFileName},
			},
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
			},
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"format":                {"foobar"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"format":                {"png"},
			},
			expectedFormat: "png",
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"format":                {"gif"},
			},
			expectedFormat: "gif",
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"format":                {"jpeg"},
				"quality":               {"-10"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"format":                {"jpeg"},
				"quality":               {"50"},
			},
			expectedFormat: "jpeg",
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"width":                 {"-100"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"width":                 {"9001"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"width":                 {"100"},
			},
			expectedWidth: 100,
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"height":                {"-100"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"height":                {"9001"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			query: url.Values{
				imageserver.SourceParam: {testdata.MediumFileName},
				"height":                {"100"},
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
