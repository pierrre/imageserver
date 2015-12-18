package imageserver

import (
	"fmt"
	"testing"
)

var _ Server = ServerFunc(nil)

func TestServerFunc(t *testing.T) {
	called := false
	srv := ServerFunc(func(params Params) (*Image, error) {
		called = true
		return &Image{}, nil
	})
	srv.Get(Params{})
	if !called {
		t.Fatal("not called")
	}
}

var _ Server = &StaticServer{}

func TestStaticServer(t *testing.T) {
	im := &Image{}
	srv := &StaticServer{
		Image: im,
		Error: nil,
	}
	out, err := srv.Get(Params{})
	if err != nil {
		t.Fatal(err)
	}
	if out != im {
		t.Fatal("not equal")
	}
}

var _ Server = &SourceServer{}

func TestSourceServer(t *testing.T) {
	srv := &SourceServer{
		Server: ServerFunc(func(params Params) (*Image, error) {
			if !params.Has(SourceParam) {
				t.Fatal("no source param")
			}
			if params.Has("foo") {
				t.Fatal("unexpected param")
			}
			return &Image{}, nil
		}),
	}
	_, err := srv.Get(Params{
		SourceParam: "source",
		"foo":       "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSourceServerErrorServer(t *testing.T) {
	srv := &SourceServer{
		Server: &StaticServer{
			Error: fmt.Errorf("error"),
		},
	}
	_, err := srv.Get(Params{SourceParam: "source"})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSourceServerErrorNoSource(t *testing.T) {
	srv := &SourceServer{}
	_, err := srv.Get(Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestNewLimitServer(t *testing.T) {
	// TODO test limit
	srv := NewLimitServer(&StaticServer{Image: &Image{}}, 1)
	_, err := srv.Get(Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewLimitServerZero(t *testing.T) {
	// TODO ?
	NewLimitServer(&StaticServer{Image: &Image{}}, 0)
}
