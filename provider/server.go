package provider

import (
	"github.com/pierrre/imageserver"
)

// Server is a Provider Server
type Server struct {
	Provider Provider
}

// Get get an Image from the Provider using the "source" param
func (s *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	source, err := params.Get("source")
	if err != nil {
		return nil, newSourceParamError("missing")
	}

	image, err := s.Provider.Get(source, params)
	if err != nil {
		if err, ok := err.(*SourceError); ok {
			return nil, newSourceParamError(err.Message)
		}
		return nil, err
	}

	return image, nil
}

func newSourceParamError(message string) error {
	return &imageserver.ParamError{Param: "source", Message: message}
}
