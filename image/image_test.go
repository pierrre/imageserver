package image

import "image"

var testImageBounds = image.Rect(0, 0, 256, 256)

func NewTestImage() image.Image {
	return image.NewRGBA(testImageBounds)
}
