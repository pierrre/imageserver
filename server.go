// Package imageserver provides an Image server toolkit.
package imageserver

// Server serves an Image.
type Server interface {
	Get(Params) (*Image, error)
}

// ServerFunc is a Server func.
type ServerFunc func(Params) (*Image, error)

// Get implements Server.
func (f ServerFunc) Get(params Params) (*Image, error) {
	return f(params)
}

// NewLimitServer creates a new Server that limits the number of concurrent executions.
//
// It uses a buffered channel to limit the number of concurrent executions.
func NewLimitServer(s Server, limit int) Server {
	return &limitServer{
		Server:  s,
		limitCh: make(chan struct{}, limit),
	}
}

type limitServer struct {
	Server
	limitCh chan struct{}
}

func (s *limitServer) Get(params Params) (*Image, error) {
	s.limitCh <- struct{}{}
	defer func() {
		<-s.limitCh
	}()
	return s.Server.Get(params)
}
