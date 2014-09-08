package source

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

func TestInterface(t *testing.T) {
	var _ imageserver_http.Parser = &Parser{}
	var _ imageserver_http.Resolver = &Parser{}
}

func TestParse(t *testing.T) {
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

	parameters := make(imageserver.Parameters)

	parser := &Parser{}

	err = parser.Parse(request, parameters)
	if err != nil {
		t.Fatal(err)
	}

	v, err := parameters.GetString("source")
	if err != nil {
		t.Fatal(err)
	}
	if v != source {
		t.Fatal("wrong value")
	}
}

func TestResolve(t *testing.T) {
	parser := &Parser{}

	httpParameter := parser.Resolve("source")
	if httpParameter != "source" {
		t.Fatal("not equals")
	}

	httpParameter = parser.Resolve("foobar")
	if httpParameter != "" {
		t.Fatal("not equals")
	}
}
