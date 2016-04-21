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
	type TC struct {
		query              url.Values
		expectedParams     imageserver.Params
		expectedParamError string
	}
	for _, tc := range []TC{
		{},
		{
			query: url.Values{"width": {"100"}},
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
		},
		{
			query: url.Values{"height": {"100"}},
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"height": 100,
			}},
		},
		{
			query: url.Values{"resampling": {"lanczos"}},
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"resampling": "lanczos",
			}},
		},
		{
			query: url.Values{"mode": {"fit"}},
			expectedParams: imageserver.Params{resizeParam: imageserver.Params{
				"mode": "fit",
			}},
		},
		{
			query:              url.Values{"width": {"invalid"}},
			expectedParamError: resizeParam + ".width",
		},
		{
			query:              url.Values{"height": {"invalid"}},
			expectedParamError: resizeParam + ".height",
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
			prc := &ResizeParser{}
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

func TestResizeParserResolve(t *testing.T) {
	prc := &ResizeParser{}
	httpParam := prc.Resolve(resizeParam + ".width")
	if httpParam != "width" {
		t.Fatal("not equal")
	}
}

func TestResizeParserResolveNoMatch(t *testing.T) {
	prc := &ResizeParser{}
	httpParam := prc.Resolve("foo")
	if httpParam != "" {
		t.Fatal("not equal")
	}
}
