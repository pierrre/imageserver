package imageserver

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestImage(t *testing.T) {
	i1 := image.NewRGBA(image.Rect(0, 0, 500, 400))
	buffer1 := new(bytes.Buffer)
	png.Encode(buffer1, i1)

	im1 := &Image{
		Type: "png",
		Data: buffer1.Bytes(),
	}

	data, err := im1.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	im2 := new(Image)
	err = im2.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}

	buffer2 := bytes.NewBuffer(im2.Data)
	i2, err := png.Decode(buffer2)
	if err != nil {
		t.Fatal(err)
	}

	if !i2.Bounds().Eq(i1.Bounds()) {
		t.Fatal("Image bounds not equals %s %s", i1.Bounds(), i2.Bounds())
	}
}
