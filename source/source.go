// Package source provides a imageserver.Server implementation that forwards calls to the underlying Server with only the "source" param.
package source

import (
	"github.com/pierrre/imageserver"
)

// Param is the source param name.
const Param = "source"

// Server is a imageserver.Server implementation that forwards calls to the underlying Server with only the "source" param.
//
// It should be used to cache the source Image.
type Server struct {
	imageserver.Server
}

// Get implements imageserver.Server.
func (s *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	src, err := params.Get(Param)
	if err != nil {
		return nil, err
	}
	params = imageserver.Params{Param: src}
	return s.Server.Get(params)
}
