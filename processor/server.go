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
func (s *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	image, err := s.Server.Get(params)
	if err != nil {
		return nil, err
	}

	image, err = s.Processor.Process(image, params)
	if err != nil {
		return nil, err
	}

	return image, err
}
