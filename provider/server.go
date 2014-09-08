package provider

import (
	"github.com/pierrre/imageserver"
)

// Server is a Provider Server
type Server struct {
	Provider Provider
}

// Get get an Image from the Provider using the "source" parameter
func (s *Server) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	source, err := parameters.Get("source")
	if err != nil {
		return nil, newSourceParameterError("missing")
	}

	image, err := s.Provider.Get(source, parameters)
	if err != nil {
		if err, ok := err.(*SourceError); ok {
			return nil, newSourceParameterError(err.Message)
		}
		return nil, err
	}

	return image, nil
}

func newSourceParameterError(message string) error {
	return &imageserver.ParameterError{Parameter: "source", Message: message}
}
