package graphicsmagick

import (
	"testing"
	"time"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Handler = &Handler{}

func TestHandle(t *testing.T) {
	hdr := &Handler{
		Executable: "gm",
	}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := hdr.Handle(testdata.Medium, params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleErrorTimeout(t *testing.T) {
	hdr := &Handler{
		Executable: "gm",
		Timeout:    1 * time.Nanosecond,
	}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := hdr.Handle(testdata.Medium, params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}
