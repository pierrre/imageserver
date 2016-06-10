package imageserver

import (
	"context"
)

// Handler handles an Image and returns an Image.
type Handler interface {
	Handle(context.Context, *Image, Params) (*Image, error)
}

// HandlerFunc is a Handler func.
type HandlerFunc func(context.Context, *Image, Params) (*Image, error)

// Handle implements Handler.
func (f HandlerFunc) Handle(ctx context.Context, im *Image, params Params) (*Image, error) {
	return f(ctx, im, params)
}

// HandlerServer is a Server implementation that calls a Handler.
type HandlerServer struct {
	Server
	Handler Handler
}

// Get implements Server.
func (srv *HandlerServer) Get(ctx context.Context, params Params) (*Image, error) {
	im, err := srv.Server.Get(ctx, params)
	if err != nil {
		return nil, err
	}
	im, err = srv.Handler.Handle(ctx, im, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}
