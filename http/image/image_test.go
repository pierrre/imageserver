package image

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

var _ imageserver_http.Parser = &FormatParser{}

func TestFormatParserParse(t *testing.T) {
	parser := &FormatParser{}
	req, err := http.NewRequest("GET", "http://localhost?format=jpg", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	format, err := params.GetString("format")
	if err != nil {
		t.Fatal(err)
	}
	if format != "jpeg" {
		t.Fatal("not equals")
	}
}

func TestFormatParserParseUndefined(t *testing.T) {
	parser := &FormatParser{}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("format") {
		t.Fatal("should not be set")
	}
}

func TestFormatParserParseError(t *testing.T) {
	parser := &FormatParser{}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{"format": 666}
	err = parser.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestFormatParserResolve(t *testing.T) {
	parser := &FormatParser{}

	httpParam := parser.Resolve("format")
	if httpParam != "format" {
		t.Fatal("not equals")
	}

	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ imageserver_http.Parser = &QualityParser{}

func TestQualityParserParse(t *testing.T) {
	parser := &QualityParser{}
	req, err := http.NewRequest("GET", "http://localhost?quality=50", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	quality, err := params.GetInt("quality")
	if err != nil {
		t.Fatal(err)
	}
	if quality != 50 {
		t.Fatal("not equals")
	}
}

func TestQualityParserParseError(t *testing.T) {
	parser := &QualityParser{}
	req, err := http.NewRequest("GET", "http://localhost?quality=foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
	if err, ok := err.(*imageserver.ParamError); !ok {
		t.Fatal("wrong error type")
	} else {
		param := err.Param
		if param != "quality" {
			t.Fatal("wrong param")
		}
	}
}

func TestQualityParserResolve(t *testing.T) {
	parser := &QualityParser{}
	httpParam := parser.Resolve("quality")
	if httpParam != "quality" {
		t.Fatal("not equals")
	}
	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}
