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
			query: url.Values{"width": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"width": 100,
			}},
		},
		{
			query: url.Values{"height": {"100"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"height": 100,
			}},
		},
		{
			query: url.Values{"fill": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"fill": true,
			}},
		},
		{
			query: url.Values{"ignore_ratio": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"ignore_ratio": true,
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
			query: url.Values{"background": {"123abc"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"background": "123abc",
			}},
		},
		{
			query: url.Values{"extent": {"true"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"extent": true,
			}},
		},
		{
			query: url.Values{"format": {"jpeg"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"format": "jpeg",
			}},
		},
		{
			query: url.Values{"quality": {"75"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"quality": 75,
			}},
		},
		{
			query:              url.Values{"width": {"invalid"}},
			expectedParamError: globalParam + ".width",
		},
		{
			query:              url.Values{"height": {"invalid"}},
			expectedParamError: globalParam + ".height",
		},
		{
			query:              url.Values{"fill": {"invalid"}},
			expectedParamError: globalParam + ".fill",
		},
		{
			query:              url.Values{"ignore_ratio": {"invalid"}},
			expectedParamError: globalParam + ".ignore_ratio",
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
			query:              url.Values{"quality": {"invalid"}},
			expectedParamError: globalParam + ".quality",
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
