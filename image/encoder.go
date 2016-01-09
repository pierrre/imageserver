package image

import (
	"bytes"
	"fmt"
	"image"
	"io"

	"github.com/pierrre/imageserver"
)

// Encoder encodes an Image.
//
// An Encoder must encode to only one specific format.
type Encoder interface {
	Encode(io.Writer, image.Image, imageserver.Params) error
	Changer
}

// EncoderFunc is an Encoder func.
type EncoderFunc func(io.Writer, image.Image, imageserver.Params) error

// Encode implements Encoder.
func (f EncoderFunc) Encode(w io.Writer, nim image.Image, params imageserver.Params) error {
	return f(w, nim, params)
}

// Change implements Encoder.
func (f EncoderFunc) Change(params imageserver.Params) bool {
	return true
}

var encoders = make(map[string]Encoder)

// RegisterEncoder registers an Encoder for a format.
func RegisterEncoder(format string, enc Encoder) {
	encoders[format] = enc
}

func getEncoder(format string) (Encoder, error) {
	enc, ok := encoders[format]
	if !ok {
		return nil, fmt.Errorf("no registered encoder for format \"%s\"", format)
	}
	return enc, nil
}

func getEncoderFormat(defaultFormat string, params imageserver.Params) (Encoder, string, error) {
	fromParams := false
	format := defaultFormat
	if params.Has("format") || defaultFormat == "" {
		fromParams = true
		var err error
		format, err = params.GetString("format")
		if err != nil {
			return nil, "", err
		}
	}
	enc, err := getEncoder(format)
	if err != nil {
		if fromParams {
			err = &imageserver.ParamError{Param: "format", Message: err.Error()}
		}
		return nil, "", err
	}
	return enc, format, nil
}

func encode(nim image.Image, format string, enc Encoder, params imageserver.Params) (*imageserver.Image, error) {
	buf := new(bytes.Buffer)
	err := enc.Encode(buf, nim, params)
	if err != nil {
		return nil, err
	}
	im := &imageserver.Image{
		Format: format,
		Data:   buf.Bytes(),
	}
	return im, nil
}

// Decode decodes a raw Image to a Go Image.
//
// It returns an error if the decoded Image format does not match the raw Image format.
func Decode(im *imageserver.Image) (image.Image, error) {
	nim, format, err := image.Decode(bytes.NewReader(im.Data))
	if err != nil {
		return nil, &imageserver.ImageError{Message: err.Error()}
	}
	if format != im.Format {
		return nil, &imageserver.ImageError{Message: fmt.Sprintf("decoded format \"%s\" does not match image format \"%s\"", format, im.Format)}
	}
	return nim, nil
}
