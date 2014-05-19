package source

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestParse(t *testing.T) {
	parser := &SourceParser{}

	request, err := http.NewRequest("GET", "http://localhost?source=foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	parameters := make(imageserver.Parameters)

	err = parser.Parse(request, parameters)
	if err != nil {
		t.Fatal(err)
	}

	v, err := parameters.GetString("source")
	if err != nil {
		t.Fatal(err)
	}
	if v != "foo" {
		t.Fatal("wrong value")
	}
}
