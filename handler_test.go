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
	hdr.Handle(&Image{}, Params{})
	if !called {
		t.Fatal("not called")
	}
}

var _ Handler = &IdentityHandler{}

func TestIdentityHandler(t *testing.T) {
	im := &Image{}
	hdr := &IdentityHandler{}
	out, err := hdr.Handle(im, Params{})
	if err != nil {
		t.Fatal(err)
	}
	if out != im {
		t.Fatal("not equal")
	}
}

var _ Server = &HandlerServer{}

func TestHandlerServer(t *testing.T) {
	srv := &HandlerServer{
		Server: &StaticServer{
			Image: &Image{},
		},
		Handler: &IdentityHandler{},
	}
	_, err := srv.Get(Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerServerErrorServer(t *testing.T) {
	srv := &HandlerServer{
		Server: &StaticServer{
			Error: fmt.Errorf("error"),
		},
	}
	_, err := srv.Get(Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestHandlerServerErrorHandler(t *testing.T) {
	srv := &HandlerServer{
		Server: &StaticServer{
			Image: &Image{},
		},
		Handler: HandlerFunc(func(im *Image, params Params) (*Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := srv.Get(Params{})
	if err == nil {
		t.Fatal("no error")
	}
}
