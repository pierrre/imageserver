package gift

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

var _ imageserver_http.Parser = &RotateParser{}

func TestRotateParserParse(t *testing.T) {
	prs := &RotateParser{}
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
			name:  "Rotation",
			query: url.Values{"rotation": {"90"}},
			expectedParams: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 90.0,
			}},
		},
		{
			name:  "Interpolation",
			query: url.Values{"interpolation": {"cubic"}},
			expectedParams: imageserver.Params{rotateParam: imageserver.Params{
				"interpolation": "cubic",
			}},
		},
		{
			name:  "Background",
			query: url.Values{"background": {"FF0000"}},
			expectedParams: imageserver.Params{rotateParam: imageserver.Params{
				"background": "FF0000",
			}},
		},
		{
			name:               "RotationInvalid",
			query:              url.Values{"rotation": {"invalid"}},
			expectedParamError: rotateParam + ".rotation",
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

func TestRotateParserResolve(t *testing.T) {
	prs := &RotateParser{}
	httpParam := prs.Resolve(rotateParam + ".rotation")
	if httpParam != "rotation" {
		t.Fatal("not equal")
	}
}

func TestRotateParserResolveNoMatch(t *testing.T) {
	prs := &RotateParser{}
	httpParam := prs.Resolve("foo")
	if httpParam != "" {
		t.Fatal("not equal")
	}
}
