package graphicsmagick

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestParse(t *testing.T) {
	parser := &GraphicsMagickParser{}

	request, err := http.NewRequest("GET", "http://localhost?width=200&height=100&fill=true&ignore_ratio=true&only_shrink_larger=true&only_enlarge_smaller=true&background=ffffff&extent=true&format=jpeg&quality=85", nil)
	if err != nil {
		t.Fatal(err)
	}

	parameters := make(imageserver.Parameters)

	err = parser.Parse(request, parameters)
	if err != nil {
		t.Fatal(err)
	}

	gmParameters, err := parameters.GetParameters("graphicsmagick")
	if err != nil {
		t.Fatal(err)
	}

	width, err := gmParameters.GetInt("width")
	if err != nil {
		t.Fatal(err)
	}
	if width != 200 {
		t.Fatal("Wrong value")
	}

	height, err := gmParameters.GetInt("height")
	if err != nil {
		t.Fatal(err)
	}
	if height != 100 {
		t.Fatal("Wrong value")
	}

	ignoreRatio, err := gmParameters.GetBool("ignore_ratio")
	if err != nil {
		t.Fatal(err)
	}
	if ignoreRatio != true {
		t.Fatal("Wrong value")
	}

	onlyShrinkLarger, err := gmParameters.GetBool("only_shrink_larger")
	if err != nil {
		t.Fatal(err)
	}
	if onlyShrinkLarger != true {
		t.Fatal("Wrong value")
	}

	onlyEnlargeSmaller, err := gmParameters.GetBool("only_enlarge_smaller")
	if err != nil {
		t.Fatal(err)
	}
	if onlyEnlargeSmaller != true {
		t.Fatal("Wrong value")
	}

	background, err := gmParameters.GetString("background")
	if err != nil {
		t.Fatal(err)
	}
	if background != "ffffff" {
		t.Fatal("Wrong value")
	}

	extent, err := gmParameters.GetBool("extent")
	if err != nil {
		t.Fatal(err)
	}
	if extent != true {
		t.Fatal("Wrong value")
	}

	format, err := gmParameters.GetString("format")
	if err != nil {
		t.Fatal(err)
	}
	if format != "jpeg" {
		t.Fatal("Wrong value")
	}

	quality, err := gmParameters.GetString("quality")
	if err != nil {
		t.Fatal(err)
	}
	if quality != "85" {
		t.Fatal("Wrong value")
	}
}

func TestParseEmpty(t *testing.T) {
	parser := &GraphicsMagickParser{}

	request, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	parameters := make(imageserver.Parameters)

	err = parser.Parse(request, parameters)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseError(t *testing.T) {
	parser := &GraphicsMagickParser{}

	for key, value := range map[string]string{
		"width":                "foo",
		"height":               "foo",
		"fill":                 "foo",
		"ignore_ratio":         "foo",
		"only_shrink_larger":   "foo",
		"only_enlarge_smaller": "foo",
		"extent":               "foo",
	} {
		request, err := http.NewRequest("GET", "http://localhost?"+key+"="+value, nil)
		if err != nil {
			t.Fatal(err)
		}

		parameters := make(imageserver.Parameters)

		err = parser.Parse(request, parameters)
		if err == nil {
			t.Fatal("no error")
		}
	}
}

func TestParseErrorDimensionNegative(t *testing.T) {
	parser := &GraphicsMagickParser{}

	request, err := http.NewRequest("GET", "http://localhost?width=-42", nil)
	if err != nil {
		t.Fatal(err)
	}

	parameters := make(imageserver.Parameters)

	err = parser.Parse(request, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}
