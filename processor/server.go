package processor

import (
	"github.com/pierrre/imageserver"
)

// Server is a Processor Server
type Server struct {
	imageserver.Server
	Processor Processor
}

// Get gets an Image from the underlying Server, then processes it with the Processor
func (s *Server) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	image, err := s.Server.Get(parameters)
	if err != nil {
		return nil, err
	}

	image, err = s.Processor.Process(image, parameters)
	if err != nil {
		return nil, err
	}

	return image, err
}
