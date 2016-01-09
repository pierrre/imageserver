package gamma

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

var _ imageserver_http.Parser = &CorrectionParser{}

func TestCorrectionParserParse(t *testing.T) {
	parser := &CorrectionParser{}
	req, err := http.NewRequest("GET", "http://localhost?gamma_correction=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	res, err := params.GetBool("gamma_correction")
	if err != nil {
		t.Fatal(err)
	}
	if res != true {
		t.Fatalf("unexpected result: got %t, want %t", res, true)
	}
}

func TestCorrectionParserResolve(t *testing.T) {
	parser := &CorrectionParser{}

	res := parser.Resolve("gamma_correction")
	expected := "gamma_correction"
	if res != expected {
		t.Fatalf("got %s, want %s", res, expected)
	}

	res = parser.Resolve("foobar")
	expected = ""
	if res != expected {
		t.Fatalf("got %s, want %s", res, expected)
	}
}
