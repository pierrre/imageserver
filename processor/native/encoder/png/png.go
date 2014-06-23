package png

import (
	"bytes"
	"image"
	"image/png"

	"github.com/pierrre/imageserver"
)

// Encode encode a native Go image to an Image
func Encode(nativeImage image.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, nativeImage)
	if err != nil {
		return nil, err
	}

	image := &imageserver.Image{
		Format: "png",
		Data:   buf.Bytes(),
	}

	return image, nil
}
