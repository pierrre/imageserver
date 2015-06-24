// Package imageserver provides an Image server toolkit.
package imageserver

// Server represents an Image server
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

// SourceServer is a source Server.
//
// It forwards to the underlying Server with only the source param.
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
	if limit <= 0 {
		return s
	}
	limitCh := make(chan struct{}, limit)
	return ServerFunc(func(params Params) (*Image, error) {
		limitCh <- struct{}{}
		defer func() {
			<-limitCh
		}()
		return s.Get(params)
	})
}

// StaticServer is an Image Server that always returns the same Image and error.
type StaticServer struct {
	Image *Image
	Error error
}

// Get implements Server.
func (s *StaticServer) Get(params Params) (*Image, error) {
	return s.Image, s.Error
}
