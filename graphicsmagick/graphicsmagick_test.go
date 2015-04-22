package graphicsmagick

import (
	"testing"
	"time"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestGet(t *testing.T) {
	im := testdata.Medium
	var server imageserver.Server
	server = imageserver.ServerFunc(func(imageserver.Params) (*imageserver.Image, error) {
		return im, nil
	})
	server = &Server{
		Server:     server,
		Executable: "gm",
	}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := server.Get(params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetErrorTimeout(t *testing.T) {
	im := testdata.Medium
	var server imageserver.Server
	server = imageserver.ServerFunc(func(imageserver.Params) (*imageserver.Image, error) {
		return im, nil
	})
	server = &Server{
		Server:     server,
		Executable: "gm",
		Timeout:    1 * time.Nanosecond,
	}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
}
