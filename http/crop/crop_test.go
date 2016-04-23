package crop

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

var _ imageserver_http.Parser = &Parser{}

func TestParse(t *testing.T) {
	ps := &Parser{}
	type TC struct {
		url                string
		expectedParams     imageserver.Params
		expectedParamError string
	}
	for _, tc := range []TC{
		{
			url:            "http://localhost",
			expectedParams: imageserver.Params{},
		},
		{
			url: "http://localhost?crop=1,2|3,4",
			expectedParams: imageserver.Params{param: imageserver.Params{
				"min_x": 1,
				"min_y": 2,
				"max_x": 3,
				"max_y": 4,
			}},
		},
		{
			url:                "http://localhost?crop=invalid",
			expectedParamError: "crop",
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			params := imageserver.Params{}
			err = ps.Parse(req, params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && tc.expectedParamError == err.Param {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedParamError != "" {
				t.Fatalf("no error, expected: %s", tc.expectedParamError)
			}
			if !reflect.DeepEqual(params, tc.expectedParams) {
				t.Fatalf("unexpected params: got %s, want %s", params, tc.expectedParams)
			}
		}()
	}
}

func TestResolve(t *testing.T) {
	ps := &Parser{}
	type TC struct {
		param    string
		expected string
	}
	for _, tc := range []TC{
		{
			param:    param,
			expected: param,
		},
		{
			param:    param + ".min_x",
			expected: param,
		},
		{
			param:    "foobar",
			expected: "",
		},
	} {
		httpParam := ps.Resolve(tc.param)
		if httpParam != tc.expected {
			t.Logf("param %s", tc.param)
			t.Fatalf("unexpected result: got '%s', want %s''", httpParam, tc.expected)
		}
	}
}
