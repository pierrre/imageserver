// Package imageserver provides an Image server toolkit.
package imageserver

// Server serves an Image.
type Server interface {
	Get(Params) (*Image, error)
}

// ServerFunc is a Server func.
type ServerFunc func(params Params) (*Image, error)

// Get implements Server.
func (f ServerFunc) Get(params Params) (*Image, error) {
	return f(params)
}

// SourceParam is the source Param name.
const SourceParam = "source"

// SourceServer is a Server implementation that forwards calls to the underlying Server with only the "source" param.
//
// It should be used to cache the source Image.
type SourceServer struct {
	Server
}

// Get implements Server.
func (s *SourceServer) Get(params Params) (*Image, error) {
	source, err := params.Get(SourceParam)
	if err != nil {
		return nil, err
	}
	params = Params{SourceParam: source}
	return s.Server.Get(params)
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
