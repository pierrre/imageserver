package native

import (
	"bytes"
	"image"

	"github.com/pierrre/imageserver"
)

/*
NativeProcessor is an Image Processor that uses the native Go Image

Steps:

- decode (from raw data to Go image)

- process (Go image)

- encode (from Go image to raw data)
*/
type NativeProcessor struct {
	Decoder   Decoder
	Processor Processor
	Encoder   Encoder
}

// Process processes an Image using native Go Image
func (processor *NativeProcessor) Process(rawImage *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	nativeImage, decodedFormat, err := processor.Decoder.Decode(rawImage, parameters)
	if err != nil {
		return nil, err
	}

	nativeImage, err = processor.Processor.Process(nativeImage, parameters)
	if err != nil {
		return nil, err
	}

	rawImage, err = processor.Encoder.Encode(nativeImage, decodedFormat, parameters)
	if err != nil {
		return nil, err
	}

	return rawImage, nil
}

// Decoder decodes a raw Image to a native Image
type Decoder interface {
	Decode(rawImage *imageserver.Image, parameters imageserver.Parameters) (nativeImage image.Image, decodedFormat string, err error)
}

// DecoderFunc is a Decoder func
type DecoderFunc func(rawImage *imageserver.Image, parameters imageserver.Parameters) (nativeImage image.Image, decodedFormat string, err error)

// Decode calls the func
func (f DecoderFunc) Decode(rawImage *imageserver.Image, parameters imageserver.Parameters) (nativeImage image.Image, decodedFormat string, err error) {
	return f(rawImage, parameters)
}

var baseDecoder = DecoderFunc(func(rawImage *imageserver.Image, parameters imageserver.Parameters) (nativeImage image.Image, decodedFormat string, err error) {
	return image.Decode(bytes.NewReader(rawImage.Data))
})

// GetBaseDecoder returns a base Decoder
//
// It decodes Image using image.Decode()
func GetBaseDecoder() Decoder {
	return baseDecoder
}

// Processor processes a native Go Image
type Processor interface {
	Process(image.Image, imageserver.Parameters) (image.Image, error)
}

// ProcessorFunc is a Processor func
type ProcessorFunc func(nativeImage image.Image, parameters imageserver.Parameters) (image.Image, error)

// Process calls the func
func (f ProcessorFunc) Process(nativeImage image.Image, parameters imageserver.Parameters) (image.Image, error) {
	return f(nativeImage, parameters)
}

// Encoder encodes a native Image to a raw Image
type Encoder interface {
	Encode(nativeImage image.Image, decodedFormat string, parameters imageserver.Parameters) (rawImage *imageserver.Image, err error)
}

// EncoderFunc is a Encoder func
type EncoderFunc func(nativeImage image.Image, decodedFormat string, parameters imageserver.Parameters) (rawImage *imageserver.Image, err error)

// Encode calls the func
func (f EncoderFunc) Encode(nativeImage image.Image, decodedFormat string, parameters imageserver.Parameters) (rawImage *imageserver.Image, err error) {
	return f(nativeImage, decodedFormat, parameters)
}
