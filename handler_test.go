package imageserver

import (
	"fmt"
	"testing"
)

var _ Handler = HandlerFunc(nil)

func TestHandlerFunc(t *testing.T) {
	called := false
	hdr := HandlerFunc(func(im *Image, params Params) (*Image, error) {
		called = true
		return im, nil
	})
	_, _ = hdr.Handle(&Image{}, Params{})
	if !called {
		t.Fatal("not called")
	}
}

var _ Server = &HandlerServer{}

func TestHandlerServer(t *testing.T) {
	srv := &HandlerServer{
		Server: ServerFunc(func(params Params) (*Image, error) {
			return &Image{}, nil
		}),
		Handler: HandlerFunc(func(im *Image, params Params) (*Image, error) {
			return im, nil
		}),
	}
	_, err := srv.Get(Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerServerErrorServer(t *testing.T) {
	srv := &HandlerServer{
		Server: ServerFunc(func(params Params) (*Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := srv.Get(Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestHandlerServerErrorHandler(t *testing.T) {
	srv := &HandlerServer{
		Server: ServerFunc(func(params Params) (*Image, error) {
			return &Image{}, nil
		}),
		Handler: HandlerFunc(func(im *Image, params Params) (*Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := srv.Get(Params{})
	if err == nil {
		t.Fatal("no error")
	}
}
