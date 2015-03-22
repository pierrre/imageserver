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

	v, err := params.GetString("source")
	if err != nil {
		t.Fatal(err)
	}
	if v != "foo" {
		t.Fatal("wrong value")
	}
}

func TestSourceParserResolve(t *testing.T) {
	parser := &SourceParser{}

	httpParam := parser.Resolve("source")
	if httpParam != "source" {
		t.Fatal("not equals")
	}

	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourcePathParser{}

var _ Parser = &SourceURLParser{}

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
