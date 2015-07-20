package image

import (
	"github.com/pierrre/imageserver"
)

// Server is an Image Server that uses Go Image.
type Server struct {
	imageserver.Server

	// Force to handle the Image
	Force bool

	// Optional Processor
	Processor Processor
}

// Get implements Server.
func (srv *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	im, err := srv.Server.Get(params)
	if err != nil {
		return nil, err
	}
	enc, format, err := getEncoderFormat(im.Format, params)
	if err != nil {
		if _, ok := err.(*imageserver.ParamError); !ok {
			err = &imageserver.ImageError{Message: err.Error()}
		}
		return nil, err
	}
	if !srv.change(im, format, enc, params) {
		return im, nil
	}
	nim, err := Decode(im)
	if err != nil {
		return nil, err
	}
	if srv.Processor != nil {
		nim, err = srv.Processor.Process(nim, params)
		if err != nil {
			return nil, err
		}
	}
	im, err = encode(nim, format, enc, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}

func (srv *Server) change(im *imageserver.Image, format string, enc Encoder, params imageserver.Params) bool {
	if srv.Force {
		return true
	}
	if format != im.Format {
		return true
	}
	if srv.Processor != nil && srv.Processor.Change(params) {
		return true
	}
	if enc.Change(params) {
		return true
	}
	return false
}
