package imageserver

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"testing"
)

func TestImage(t *testing.T) {
	image1 := CreateImage(500, 400)
	data, err := image1.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	image2 := new(Image)
	err = image2.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}

	if image2.Type != image1.Type {
		t.Fatalf("Types not equals: %s %s", image1.Type, image2.Type)
	}
	if len(image2.Data) != len(image1.Data) {
		t.Fatalf("Data not the same size %d %d", len(image1.Data), len(image2.Data))
	}
	for i, b2 := range image2.Data {
		b1 := image1.Data[i]
		if b2 != b1 {
			t.Fatalf("Data not equals at index %d: %d %d", i, b1, b2)
		}
	}
}

func CreateImage(width, height int) *Image {
	baseImage := CreateBaseImage(width, height)
	buffer := new(bytes.Buffer)
	png.Encode(buffer, baseImage)
	return &Image{
		Type: "png",
		Data: buffer.Bytes(),
	}
}

func CreateBaseImage(width, height int) *image.NRGBA {
	i := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y, height := 0, i.Bounds().Dy(); y < height; y++ {
		for x, width := 0, i.Bounds().Dx(); x < width; x++ {
			i.Set(x, y, randColor())
		}
	}
	return i
}

func randColor() color.RGBA {
	return color.RGBA{
		R: randColorComponent(),
		G: randColorComponent(),
		B: randColorComponent(),
		A: randColorComponent(),
	}
}

func randColorComponent() uint8 {
	return uint8(rand.Int31n(256))
}
