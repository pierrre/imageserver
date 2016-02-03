package image

import (
	"github.com/pierrre/imageserver"
)

// Server is an Image Server for an Go Image Provider.
type Server struct {
	Provider Provider
}

// Get implements Server
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
