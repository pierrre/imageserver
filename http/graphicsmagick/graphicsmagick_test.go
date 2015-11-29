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
	type TC struct {
		query              url.Values
		expectedParams     imageserver.Params
		expectedParamError string
	}
	for _, tc := range []TC{
		{},
		{
			query: url.Values{"ignore_ratio": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"ignore_ratio": true,
			}},
		},
		{
			query: url.Values{"fill": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"fill": true,
			}},
		},
		{
			query: url.Values{"only_shrink_larger": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"only_shrink_larger": true,
			}},
		},
		{
			query: url.Values{"only_enlarge_smaller": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"only_enlarge_smaller": true,
			}},
		},
		{
			query: url.Values{"extent": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"extent": true,
			}},
		},

		{
			query: url.Values{"monochrome": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"monochrome": true,
			}},
		},
		{
			query: url.Values{"grey": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"grey": true,
			}},
		},

		{
			query: url.Values{"no_strip": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"no_strip": true,
			}},
		},
		{
			query: url.Values{"trim": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"trim": true,
			}},
		},
		{
			query: url.Values{"no_interlace": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"no_interlace": true,
			}},
		},
		{
			query: url.Values{"flip": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"flip": true,
			}},
		},
		{
			query: url.Values{"flop": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"flop": true,
			}},
		},
		{
			query: url.Values{"w": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"w": 100,
			}},
		},
		{
			query: url.Values{"h": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"h": 100,
			}},
		},
		{
			query: url.Values{"rotate": {"90"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"rotate": 90,
			}},
		},
		{
			query: url.Values{"q": {"75"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"q": 75,
			}},
		},
		{
			query: url.Values{"bg": {"fafafa"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"bg": "fafafa",
			}},
		},
		{
			query: url.Values{"gravity": {"se"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"gravity": "se",
			}},
		},
		{
			query: url.Values{"crop": {"200,200,0,0"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"crop": "200,200,0,0",
			}},
		},

		{
			query: url.Values{"format": {"jpeg"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"format": "jpeg",
			}},
		},
		{
			query:              url.Values{"ignore_ratio": {"invalid"}},
			expectedParamError: globalParam + ".ignore_ratio",
		},

		{
			query:              url.Values{"fill": {"invalid"}},
			expectedParamError: globalParam + ".fill",
		},
		{
			query:              url.Values{"only_shrink_larger": {"invalid"}},
			expectedParamError: globalParam + ".only_shrink_larger",
		},
		{
			query:              url.Values{"only_enlarge_smaller": {"invalid"}},
			expectedParamError: globalParam + ".only_enlarge_smaller",
		},
		{
			query:              url.Values{"extent": {"invalid"}},
			expectedParamError: globalParam + ".extent",
		},
		{
			query:              url.Values{"monochrome": {"invalid"}},
			expectedParamError: globalParam + ".monochrome",
		},
		{
			query:              url.Values{"grey": {"invalid"}},
			expectedParamError: globalParam + ".grey",
		},
		{
			query:              url.Values{"no_strip": {"invalid"}},
			expectedParamError: globalParam + ".no_strip",
		},
		{
			query:              url.Values{"trim": {"invalid"}},
			expectedParamError: globalParam + ".trim",
		},
		{
			query:              url.Values{"no_interlace": {"invalid"}},
			expectedParamError: globalParam + ".no_interlace",
		},
		{
			query:              url.Values{"flip": {"invalid"}},
			expectedParamError: globalParam + ".flip",
		},
		{
			query:              url.Values{"flop": {"invalid"}},
			expectedParamError: globalParam + ".flop",
		},

		{
			query:              url.Values{"w": {"invalid"}},
			expectedParamError: globalParam + ".w",
		},
		{
			query:              url.Values{"h": {"invalid"}},
			expectedParamError: globalParam + ".h",
		},
		{
			query:              url.Values{"rotate": {"invalid"}},
			expectedParamError: globalParam + ".rotate",
		},
		{
			query:              url.Values{"density": {"invalid"}},
			expectedParamError: globalParam + ".density",
		},

		{
			query:              url.Values{"q": {"invalid"}},
			expectedParamError: globalParam + ".q",
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
			p := &Parser{}
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
		}()
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
