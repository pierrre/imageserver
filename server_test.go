package imageserver

import (
	"context"
	"testing"
)

var _ Server = ServerFunc(nil)

func TestServerFunc(t *testing.T) {
	called := false
	srv := ServerFunc(func(ctx context.Context, params Params) (*Image, error) {
		called = true
		return &Image{}, nil
	})
	_, _ = srv.Get(context.Background(), Params{})
	if !called {
		t.Fatal("not called")
	}
}

func TestNewLimitServer(t *testing.T) {
	// TODO test limit
	srv := NewLimitServer(ServerFunc(func(ctx context.Context, params Params) (*Image, error) {
		return &Image{}, nil
	}), 1)
	_, err := srv.Get(context.Background(), Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewLimitServerZero(t *testing.T) {
	// TODO ?
	NewLimitServer(ServerFunc(func(ctx context.Context, params Params) (*Image, error) {
		return &Image{}, nil
	}), 0)
}
