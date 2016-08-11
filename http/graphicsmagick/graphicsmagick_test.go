package graphicsmagick

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

var _ imageserver_http.Parser = &Parser{}

func TestParse(t *testing.T) {
	p := &Parser{}
	for _, tc := range []struct {
		name               string
		query              url.Values
		expectedParams     imageserver.Params
		expectedParamError string
	}{
		{
			name: "Empty",
		},
		{
			name:  "Width",
			query: url.Values{"width": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"width": 100,
			}},
		},
		{
			name:  "Height",
			query: url.Values{"height": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"height": 100,
			}},
		},
		{
			name:  "Fill",
			query: url.Values{"fill": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"fill": true,
			}},
		},
		{
			name:  "IgnoreRatio",
			query: url.Values{"ignore_ratio": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"ignore_ratio": true,
			}},
		},
		{
			name:  "OnlyShrinkLarger",
			query: url.Values{"only_shrink_larger": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"only_shrink_larger": true,
			}},
		},
		{
			name:  "OnlyEnlargeSmaller",
			query: url.Values{"only_enlarge_smaller": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"only_enlarge_smaller": true,
			}},
		},
		{
			name:  "Background",
			query: url.Values{"background": {"123abc"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"background": "123abc",
			}},
		},
		{
			name:  "Extent",
			query: url.Values{"extent": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"extent": true,
			}},
		},
		{
			name:  "Format",
			query: url.Values{"format": {"jpeg"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"format": "jpeg",
			}},
		},
		{
			name:  "Quality",
			query: url.Values{"quality": {"75"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"quality": 75,
			}},
		},
		{
			name:               "WidthInvalid",
			query:              url.Values{"width": {"invalid"}},
			expectedParamError: globalParam + ".width",
		},
		{
			name:               "HeightInvalid",
			query:              url.Values{"height": {"invalid"}},
			expectedParamError: globalParam + ".height",
		},
		{
			name:               "FillInvalid",
			query:              url.Values{"fill": {"invalid"}},
			expectedParamError: globalParam + ".fill",
		},
		{
			name:               "IgnoreRatioInvalid",
			query:              url.Values{"ignore_ratio": {"invalid"}},
			expectedParamError: globalParam + ".ignore_ratio",
		},
		{
			name:               "OnlyShrinkLargerInvalid",
			query:              url.Values{"only_shrink_larger": {"invalid"}},
			expectedParamError: globalParam + ".only_shrink_larger",
		},
		{
			name:               "OnlyEnlargeSmallerInvalid",
			query:              url.Values{"only_enlarge_smaller": {"invalid"}},
			expectedParamError: globalParam + ".only_enlarge_smaller",
		},
		{
			name:               "ExtentInvalid",
			query:              url.Values{"extent": {"invalid"}},
			expectedParamError: globalParam + ".extent",
		},
		{
			name:               "QualityInvalid",
			query:              url.Values{"quality": {"invalid"}},
			expectedParamError: globalParam + ".quality",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := &url.URL{
				Scheme:   "http",
				Host:     "localhost",
				RawQuery: tc.query.Encode(),
			}
			req, err := http.NewRequest("GET", u.String(), nil)
			if err != nil {
				t.Fatal(err)
			}
			params := imageserver.Params{}
			err = p.Parse(req, params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && tc.expectedParamError == err.Param {
					return
				}
				t.Fatal(err)
			}
			if params.String() != tc.expectedParams.String() {
				t.Fatalf("unexpected params: got %s, want %s", params, tc.expectedParams)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	p := &Parser{}
	httpParam := p.Resolve(globalParam + ".width")
	if httpParam != "width" {
		t.Fatal("not equal")
	}
}

func TestResolveNoMatch(t *testing.T) {
	p := &Parser{}
	httpParam := p.Resolve("foo")
	if httpParam != "" {
		t.Fatal("not equal")
	}
}
