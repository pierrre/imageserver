package imageserver

// Handler handles an Image and returns an Image.
type Handler interface {
	Handle(*Image, Params) (*Image, error)
}

// HandlerFunc is a Handler func.
type HandlerFunc func(*Image, Params) (*Image, error)

// Handle implements Handler.
func (f HandlerFunc) Handle(im *Image, params Params) (*Image, error) {
	return f(im, params)
}

// HandlerServer is a Server implementation that calls a Handler.
type HandlerServer struct {
	Server
	Handler Handler
}

// Get implements Server.
func (srv *HandlerServer) Get(params Params) (*Image, error) {
	im, err := srv.Server.Get(params)
	if err != nil {
		return nil, err
	}
	im, err = srv.Handler.Handle(im, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}
