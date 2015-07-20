package nfntresize

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
			query: url.Values{"interpolation": {"lanczos3"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"interpolation": "lanczos3",
			}},
		},
		{
			query: url.Values{"mode": {"resize"}},
			expectedParams: imageserver.Params{globalParam: imageserver.Params{
				"mode": "resize",
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
				t.Fatalf("unexpected params: wanted %s, got %s", tc.expectedParams, params)
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
