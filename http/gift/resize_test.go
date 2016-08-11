package gift

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

var _ imageserver_http.Parser = &ResizeParser{}

func TestResizeParserParse(t *testing.T) {
	prs := &ResizeParser{}
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
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
		},
		{
			name:  "Height",
			query: url.Values{"height": {"100"}},
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"height": 100,
			}},
		},
		{
			name:  "Resampling",
			query: url.Values{"resampling": {"lanczos"}},
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"resampling": "lanczos",
			}},
		},
		{
			name:  "Mode",
			query: url.Values{"mode": {"fit"}},
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"mode": "fit",
			}},
		},
		{
			name:               "WidthInvalid",
			query:              url.Values{"width": {"invalid"}},
			expectedParamError: resizeParam + ".width",
		},
		{
			name:               "HeightInvalid",
			query:              url.Values{"height": {"invalid"}},
			expectedParamError: resizeParam + ".height",
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
			err = prs.Parse(req, params)
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

func TestResizeParserResolve(t *testing.T) {
	prs := &ResizeParser{}
	httpParam := prs.Resolve(resizeParam + ".width")
	if httpParam != "width" {
		t.Fatal("not equal")
	}
}

func TestResizeParserResolveNoMatch(t *testing.T) {
	prs := &ResizeParser{}
	httpParam := prs.Resolve("foo")
	if httpParam != "" {
		t.Fatal("not equal")
	}
}
