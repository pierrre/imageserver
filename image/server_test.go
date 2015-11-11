package image

import (
	"fmt"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestServer(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
	}
	_, err := srv.Get(imageserver.Params{"quality": 85})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerNoChange(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerFormat(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
	}
	_, err := srv.Get(imageserver.Params{"format": "test"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerProcessor(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
		Processor: ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nim, nil
		}),
	}
	_, err := srv.Get(imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerErrorServer(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Error: fmt.Errorf("error"),
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorFormatParam(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
	}
	_, err := srv.Get(imageserver.Params{"format": "unknown"})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestServerErrorFormatImage(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: &imageserver.Image{Format: "unknown"},
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestServerErrorDecode(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Invalid,
		},
	}
	_, err := srv.Get(imageserver.Params{"format": "test"})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestServerErrorProcessor(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
		Processor: ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorEncode(t *testing.T) {
	srv := &Server{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
	}
	_, err := srv.Get(imageserver.Params{"quality": 9001})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}
