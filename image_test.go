package imageserver

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"reflect"
	"testing"
)

func BenchmarkImageMarshalGob(b *testing.B) {
	image := &Image{
		Format: "png",
		Data:   make([]byte, 50*1024),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := image.Marshal()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkImageMarshalBinaryExp(b *testing.B) {
	image := &Image{
		Format: "png",
		Data:   make([]byte, 50*1024),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := image.MarshalBinaryExp()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkImageUnmarshalGob(b *testing.B) {
	image := &Image{
		Format: "png",
		Data:   make([]byte, 50*1024),
	}
	data, err := image.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewImageUnmarshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkImageUnmarshalBinaryExp(b *testing.B) {
	image := &Image{
		Format: "png",
		Data:   make([]byte, 50*1024),
	}
	data, err := image.MarshalBinaryExp()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewImageUnmarshalBinaryExp(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestImage(t *testing.T) {
	image1 := CreateImage(500, 400)
	data, err := image1.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	image2, err := NewImageUnmarshal(data)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(image2, image1) {
		t.Fatal("image not equals")
	}
}

func TestImageUnmarshalError(t *testing.T) {
	_, err := NewImageUnmarshal(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func CreateImage(width, height int) *Image {
	baseImage := CreateBaseImage(width, height)
	buffer := new(bytes.Buffer)
	png.Encode(buffer, baseImage)
	return &Image{
		Format: "png",
		Data:   buffer.Bytes(),
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
