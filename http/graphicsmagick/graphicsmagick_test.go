package graphicsmagick

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

var _ imageserver_http.Parser = &Parser{}

func TestParse(t *testing.T) {
	urlParams := map[string]interface{}{
		"width":                200,
		"height":               100,
		"fill":                 true,
		"ignore_ratio":         true,
		"only_shrink_larger":   true,
		"only_enlarge_smaller": true,
		"background":           "ffffff",
		"extent":               true,
		"format":               "jpeg",
		"quality":              85,
	}
	query := make(url.Values)
	for k, v := range urlParams {
		query.Add(k, fmt.Sprint(v))
	}
	request, err := http.NewRequest(
		"GET",
		(&url.URL{
			Scheme:   "http",
			Host:     "localhost",
			RawQuery: query.Encode(),
		}).String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	params := make(imageserver.Params)
	parser := &Parser{}
	err = parser.Parse(request, params)
	if err != nil {
		t.Fatal(err)
	}
	gmParams, err := params.GetParams("graphicsmagick")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range urlParams {
		param, err := gmParams.Get(k)
		if err != nil {
			t.Fatal(err)
		}
		if param != v {
			t.Fatal(fmt.Errorf("wrong value: got %#v, expected %#v", param, v))
		}
	}
}

func TestParseEmpty(t *testing.T) {
	request, err := http.NewRequest(
		"GET",
		(&url.URL{
			Scheme: "http",
			Host:   "localhost",
		}).String(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	params := make(imageserver.Params)
	parser := &Parser{}
	err = parser.Parse(request, params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseError(t *testing.T) {
	parser := &Parser{}
	for k, v := range map[string]interface{}{
		"width":                "foo",
		"height":               "foo",
		"fill":                 "foo",
		"ignore_ratio":         "foo",
		"only_shrink_larger":   "foo",
		"only_enlarge_smaller": "foo",
		"extent":               "foo",
		"quality":              "foo",
	} {
		query := make(url.Values)
		query.Add(k, fmt.Sprint(v))
		request, err := http.NewRequest(
			"GET",
			(&url.URL{
				Scheme:   "http",
				Host:     "localhost",
				RawQuery: query.Encode(),
			}).String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}
		params := make(imageserver.Params)
		err = parser.Parse(request, params)
		if err == nil {
			t.Fatal("no error")
		}
	}
}

func TestResolve(t *testing.T) {
	parser := &Parser{}
	httpParam := parser.Resolve("graphicsmagick.foo")
	if httpParam != "foo" {
		t.Fatal("not equals")
	}
	httpParam = parser.Resolve("bar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}
