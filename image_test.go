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
	i1 := CreateBaseImage(500, 400)
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
	if im2.Type != im1.Type {
		t.Fatalf("Types not equals: %s %s", im1.Type, im2.Type)
	}
	if len(im2.Data) != len(im1.Data) {
		t.Fatalf("Data not the same size %d %d", len(im1.Data), len(im2.Data))
	}
	for i, b := range im2.Data {
		if b != im1.Data[i] {
			t.Fatalf("Data not equals at index %d: %d %d", i, im1.Data[i], b)
		}
	}
	buffer2 := bytes.NewBuffer(im2.Data)
	i2, err := png.Decode(buffer2)
	if err != nil {
		t.Fatal(err)
	}
	if !i2.Bounds().Eq(i1.Bounds()) {
		t.Fatal("Image bounds not equals %s %s", i1.Bounds(), i2.Bounds())
	}
	for y, height := 0, i2.Bounds().Dy(); y < height; y++ {
		for x, width := 0, i2.Bounds().Dx(); x < width; x++ {
			c1 := i1.At(x, y)
			r1, g1, b1, a1 := c1.RGBA()
			c2 := i2.At(x, y)
			r2, g2, b2, a2 := i2.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("Colors not equals at %d,%d: %s %s", x, y, c1, c2)
			}
		}
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
