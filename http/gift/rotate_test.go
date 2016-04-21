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
	type TC struct {
		query              url.Values
		expectedParams     imageserver.Params
		expectedParamError string
	}
	for _, tc := range []TC{
		{},
		{
			query: url.Values{"rotation": {"90"}},
			expectedParams: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 90.0,
			}},
		},
		{
			query: url.Values{"interpolation": {"cubic"}},
			expectedParams: imageserver.Params{rotateParam: imageserver.Params{
				"interpolation": "cubic",
			}},
		},
		{
			query: url.Values{"background": {"FF0000"}},
			expectedParams: imageserver.Params{rotateParam: imageserver.Params{
				"background": "FF0000",
			}},
		},
		{
			query:              url.Values{"rotation": {"invalid"}},
			expectedParamError: rotateParam + ".rotation",
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
			prc := &RotateParser{}
			params := imageserver.Params{}
			err = prc.Parse(req, params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && tc.expectedParamError == err.Param {
					return
				}
				t.Fatal(err)
			}
			if params.String() != tc.expectedParams.String() {
				t.Fatalf("unexpected params: got %s, want %s", params, tc.expectedParams)
			}
		}()
	}
}

func TestRotateParserResolve(t *testing.T) {
	prc := &RotateParser{}
	httpParam := prc.Resolve(rotateParam + ".rotation")
	if httpParam != "rotation" {
		t.Fatal("not equal")
	}
}

func TestRotateParserResolveNoMatch(t *testing.T) {
	prc := &RotateParser{}
	httpParam := prc.Resolve("foo")
	if httpParam != "" {
		t.Fatal("not equal")
	}
}
