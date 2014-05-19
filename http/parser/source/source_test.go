package source

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestParse(t *testing.T) {
	source := "foo"

	query := make(url.Values)
	query.Add("source", "foo")

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

	parser := &SourceParser{}

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
