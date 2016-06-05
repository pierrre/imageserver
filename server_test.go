package imageserver

import "testing"

var _ Server = ServerFunc(nil)

func TestServerFunc(t *testing.T) {
	called := false
	srv := ServerFunc(func(params Params) (*Image, error) {
		called = true
		return &Image{}, nil
	})
	_, _ = srv.Get(Params{})
	if !called {
		t.Fatal("not called")
	}
}

func TestNewLimitServer(t *testing.T) {
	// TODO test limit
	srv := NewLimitServer(ServerFunc(func(params Params) (*Image, error) {
		return &Image{}, nil
	}), 1)
	_, err := srv.Get(Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewLimitServerZero(t *testing.T) {
	// TODO ?
	NewLimitServer(ServerFunc(func(params Params) (*Image, error) {
		return &Image{}, nil
	}), 0)
}
