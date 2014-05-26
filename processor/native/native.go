package native

import (
	"bytes"
	native_image "image"

	"github.com/pierrre/imageserver"
)

/*
NativeProcessor is an Image Processor that uses the natives Go images

Steps:

- decode (from raw data to Go image)

- process (Go image)

- encode (from Go image to raw data)
*/
type NativeProcessor struct {
	DecodeFunc func(*imageserver.Image, imageserver.Parameters) (native_image.Image, error)
	Processor
	EncodeFunc func(native_image.Image, imageserver.Parameters) (*imageserver.Image, error)
}

// Process processes an Image using natives Go images
func (processor *NativeProcessor) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	nativeImage, err := processor.DecodeFunc(image, parameters)
	if err != nil {
		return nil, err
	}

	nativeImage, err = processor.Processor.Process(nativeImage, parameters)
	if err != nil {
		return nil, err
	}

	image, err = processor.EncodeFunc(nativeImage, parameters)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// DefaultDecode is the default Decode function
func DefaultDecode(image *imageserver.Image, parameters imageserver.Parameters) (native_image.Image, error) {
	nativeImage, _, err := native_image.Decode(bytes.NewReader(image.Data))
	return nativeImage, err
}

// Processor represents a native Go image processor
type Processor interface {
	Process(native_image.Image, imageserver.Parameters) (native_image.Image, error)
}
