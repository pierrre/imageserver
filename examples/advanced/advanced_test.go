package main

import (
	"image"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver/testdata"
)

func Test(t *testing.T) {
	h := newHTTPHandler()
	for _, tc := range []struct {
		name               string
		path               string
		query              url.Values
		expectedStatusCode int
		expectedFormat     string
		expectedWidth      int
		expectedHeight     int
	}{
		{
			name:               "NoPath",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Small",
			path: testdata.SmallFileName,
		},
		{
			name: "Medium",
			path: testdata.MediumFileName,
		},
		{
			name: "Large",
			path: testdata.LargeFileName,
		},
		{
			name: "Huge",
			path: testdata.HugeFileName,
		},
		{
			name: "Animated",
			path: testdata.AnimatedFileName,
		},
		{
			name: "MediumFromCache",
			path: testdata.MediumFileName,
		},
		{
			name: "InvalidFormat",
			path: testdata.MediumFileName,
			query: url.Values{
				"format": {"foobar"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "PNG",
			path: testdata.MediumFileName,
			query: url.Values{
				"format": {"png"},
			},
			expectedFormat: "png",
		},
		{
			name: "GIF",
			path: testdata.MediumFileName,
			query: url.Values{
				"format": {"gif"},
			},
			expectedFormat: "gif",
		},
		{
			name: "JPEGInvalidQuality",
			path: testdata.MediumFileName,
			query: url.Values{
				"format":  {"jpeg"},
				"quality": {"-10"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "JPEG",
			path: testdata.MediumFileName,
			query: url.Values{
				"format":  {"jpeg"},
				"quality": {"50"},
			},
			expectedFormat: "jpeg",
		},
		{
			name: "WidthInvalidNegative",
			path: testdata.MediumFileName,
			query: url.Values{
				"width": {"-100"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "WidthInvalidTooLarge",
			path: testdata.MediumFileName,
			query: url.Values{
				"width": {"9001"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Width",
			path: testdata.MediumFileName,
			query: url.Values{
				"width": {"100"},
			},
			expectedWidth: 100,
		},
		{
			name: "HeightInvalidNegative",
			path: testdata.MediumFileName,
			query: url.Values{
				"height": {"-100"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "HeightInvalidTooLarge",
			path: testdata.MediumFileName,
			query: url.Values{
				"height": {"9001"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Height",
			path: testdata.MediumFileName,
			query: url.Values{
				"height": {"100"},
			},
			expectedHeight: 100,
		},
		{
			name: "RotationInvalid",
			path: testdata.MediumFileName,
			query: url.Values{
				"rotation": {"invalid"},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Rotation45",
			path: testdata.MediumFileName,
			query: url.Values{
				"rotation": {"45"},
			},
			expectedWidth:  1304,
			expectedHeight: 1304,
		},
		{
			name: "Rotation90",
			path: testdata.MediumFileName,
			query: url.Values{
				"rotation": {"90"},
			},
			expectedWidth:  819,
			expectedHeight: 1024,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
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
				t.Fatalf("unexpected http status: got %d, want %d", w.Code, tc.expectedStatusCode)
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
				t.Fatalf("unexpected format: got %s, want %s", format, tc.expectedFormat)
			}
			if tc.expectedWidth != 0 && im.Bounds().Dx() != tc.expectedWidth {
				t.Fatalf("unexpected width: got %d, want %d", im.Bounds().Dx(), tc.expectedWidth)
			}
			if tc.expectedHeight != 0 && im.Bounds().Dy() != tc.expectedHeight {
				t.Fatalf("unexpected height: got %d, want %d", im.Bounds().Dy(), tc.expectedHeight)
			}
		})
	}
}
