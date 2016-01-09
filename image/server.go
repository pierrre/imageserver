package image

import (
	"github.com/pierrre/imageserver"
)

// Server is a imageserver.Server implementation that gets the Image from a Provider.
//
// It uses the "format" param to determine which Encoder is used.
type Server struct {
	Provider Provider
}

// Get implements Server.
func (srv *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	enc, format, err := getEncoderFormat("", params)
	if err != nil {
		return nil, err
	}
	nim, err := srv.Provider.Get(params)
	if err != nil {
		return nil, err
	}
	im, err := encode(nim, format, enc, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}
