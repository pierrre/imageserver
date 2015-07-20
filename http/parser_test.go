package http

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

var _ Parser = ListParser{}

var _ Parser = &SourceParser{}

func TestSourceParserParse(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost?source=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := make(imageserver.Params)
	parser := &SourceParser{}
	err = parser.Parse(request, params)
	if err != nil {
		t.Fatal(err)
	}
	v, err := params.GetString(imageserver.SourceParam)
	if err != nil {
		t.Fatal(err)
	}
	if v != "foo" {
		t.Fatal("wrong value")
	}
}

func TestSourceParserResolve(t *testing.T) {
	parser := &SourceParser{}
	httpParam := parser.Resolve(imageserver.SourceParam)
	if httpParam != imageserver.SourceParam {
		t.Fatal("not equals")
	}
	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourcePathParser{}

var _ Parser = &SourceURLParser{}

func TestFormatParserInterface(t *testing.T) {
	var _ Parser = &FormatParser{}
}

func TestFormatParserParse(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost?format=jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	params := make(imageserver.Params)
	parser := &FormatParser{}
	err = parser.Parse(request, params)
	if err != nil {
		t.Fatal(err)
	}

	format, err := params.GetString("format")
	if err != nil {
		t.Fatal(err)
	}
	if format != "jpeg" {
		t.Fatal("wrong value")
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

func TestQualityParserInterface(t *testing.T) {
	var _ Parser = &QualityParser{}
}

func TestQualityParserParse(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost?quality=50", nil)
	if err != nil {
		t.Fatal(err)
	}

	params := make(imageserver.Params)
	parser := &QualityParser{}
	err = parser.Parse(request, params)
	if err != nil {
		t.Fatal(err)
	}

	quality, err := params.GetInt("quality")
	if err != nil {
		t.Fatal(err)
	}
	if quality != 50 {
		t.Fatal("wrong value")
	}
}

func TestQualityParserParseError(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost?quality=foobar", nil)
	if err != nil {
		t.Fatal(err)
	}

	params := make(imageserver.Params)
	parser := &QualityParser{}
	err = parser.Parse(request, params)
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

func TestParseQueryString(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost?string=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := make(imageserver.Params)
	ParseQueryString("string", request, params)
	s, err := params.GetString("string")
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatal("not equals")
	}
}

func TestParseQueryInt(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost?int=42", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := make(imageserver.Params)
	ParseQueryInt("int", request, params)
	i, err := params.GetInt("int")
	if err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal("not equals")
	}
}

func TestParseQueryFloat(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost?float=12.34", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := make(imageserver.Params)
	ParseQueryFloat("float", request, params)
	f, err := params.GetFloat("float")
	if err != nil {
		t.Fatal(err)
	}
	if f != 12.34 {
		t.Fatal("not equals")
	}
}
