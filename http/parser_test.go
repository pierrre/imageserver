package http

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestListParserInterface(t *testing.T) {
	var _ Parser = ListParser{}
}

func TestSourceParserInterface(t *testing.T) {
	var _ Parser = &SourceParser{}
}

func TestSourceParserParse(t *testing.T) {
	source := "foo"

	query := make(url.Values)
	query.Add("source", source)

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

	parser := &SourceParser{}

	err = parser.Parse(request, params)
	if err != nil {
		t.Fatal(err)
	}

	v, err := params.GetString("source")
	if err != nil {
		t.Fatal(err)
	}
	if v != source {
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

func TestSourcePathParserInterface(t *testing.T) {
	var _ Parser = &SourcePathParser{}
}

func TestSourceURLParserInterface(t *testing.T) {
	var _ Parser = &SourceURLParser{}
}
