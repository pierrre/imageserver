package native

import (
	"bytes"
	"image"

	"github.com/pierrre/imageserver"
)

/*
Processor is an Image Processor that uses the native Go Image

Steps:

- decode (from raw data to Go image)

- process (Go image)

- encode (from Go image to raw data)
*/
type Processor struct {
	Decoder   Decoder
	Processor ProcessorNative
	Encoder   Encoder
}

// Process processes an Image using native Go Image
func (processor *Processor) Process(rawImage *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	nativeImage, decodedFormat, err := processor.Decoder.Decode(rawImage, params)
	if err != nil {
		return nil, err
	}

	nativeImage, err = processor.Processor.Process(nativeImage, params)
	if err != nil {
		return nil, err
	}

	rawImage, err = processor.Encoder.Encode(nativeImage, decodedFormat, params)
	if err != nil {
		return nil, err
	}

	return rawImage, nil
}

// Decoder decodes a raw Image to a native Image
type Decoder interface {
	Decode(rawImage *imageserver.Image, params imageserver.Params) (nativeImage image.Image, decodedFormat string, err error)
}

// DecoderFunc is a Decoder func
type DecoderFunc func(rawImage *imageserver.Image, params imageserver.Params) (nativeImage image.Image, decodedFormat string, err error)

// Decode calls the func
func (f DecoderFunc) Decode(rawImage *imageserver.Image, params imageserver.Params) (nativeImage image.Image, decodedFormat string, err error) {
	return f(rawImage, params)
}

var baseDecoder = DecoderFunc(func(rawImage *imageserver.Image, params imageserver.Params) (nativeImage image.Image, decodedFormat string, err error) {
	return image.Decode(bytes.NewReader(rawImage.Data))
})

// GetBaseDecoder returns a base Decoder
//
// It decodes Image using image.Decode()
func GetBaseDecoder() Decoder {
	return baseDecoder
}

// ProcessorNative processes a native Go Image
type ProcessorNative interface {
	Process(image.Image, imageserver.Params) (image.Image, error)
}

// ProcessorNativeFunc is a Processor func
type ProcessorNativeFunc func(nativeImage image.Image, params imageserver.Params) (image.Image, error)

// Process calls the func
func (f ProcessorNativeFunc) Process(nativeImage image.Image, params imageserver.Params) (image.Image, error) {
	return f(nativeImage, params)
}

// Encoder encodes a native Image to a raw Image
type Encoder interface {
	Encode(nativeImage image.Image, decodedFormat string, params imageserver.Params) (rawImage *imageserver.Image, err error)
}

// EncoderFunc is a Encoder func
type EncoderFunc func(nativeImage image.Image, decodedFormat string, params imageserver.Params) (rawImage *imageserver.Image, err error)

// Encode calls the func
func (f EncoderFunc) Encode(nativeImage image.Image, decodedFormat string, params imageserver.Params) (rawImage *imageserver.Image, err error) {
	return f(nativeImage, decodedFormat, params)
}
