package internal

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"testing"
)

func TestNewDrawable(t *testing.T) {
	r := image.Rect(0, 0, 1, 1)
	for _, newImage := range []func(image.Rectangle) image.Image{
		func(r image.Rectangle) image.Image {
			return image.NewRGBA(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewRGBA64(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNRGBA(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNRGBA64(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewAlpha(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewAlpha16(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewGray(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewGray16(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewCMYK(r)
		},
		func(r image.Rectangle) image.Image {
			return image.NewPaletted(r, color.Palette{
				color.RGBA{0, 0, 0, 255},
				color.RGBA{255, 0, 0, 255},
				color.RGBA{0, 255, 0, 255},
				color.RGBA{0, 0, 255, 255},
				color.RGBA{255, 255, 255, 255},
			})
		},
		func(r image.Rectangle) image.Image {
			return image.NewYCbCr(r, image.YCbCrSubsampleRatio444)
		},
		func(r image.Rectangle) image.Image {
			return image.NewYCbCr(r, image.YCbCrSubsampleRatio422)
		},
		func(r image.Rectangle) image.Image {
			return image.NewYCbCr(r, image.YCbCrSubsampleRatio420)
		},
		func(r image.Rectangle) image.Image {
			return image.NewYCbCr(r, image.YCbCrSubsampleRatio440)
		},
		func(r image.Rectangle) image.Image {
			return image.NewYCbCr(r, image.YCbCrSubsampleRatio411)
		},
		func(r image.Rectangle) image.Image {
			return image.NewYCbCr(r, image.YCbCrSubsampleRatio410)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNYCbCrA(r, image.YCbCrSubsampleRatio444)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNYCbCrA(r, image.YCbCrSubsampleRatio422)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNYCbCrA(r, image.YCbCrSubsampleRatio420)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNYCbCrA(r, image.YCbCrSubsampleRatio440)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNYCbCrA(r, image.YCbCrSubsampleRatio411)
		},
		func(r image.Rectangle) image.Image {
			return image.NewNYCbCrA(r, image.YCbCrSubsampleRatio410)
		},
		func(r image.Rectangle) image.Image {
			return image.NewUniform(color.RGBA{})
		},
		func(r image.Rectangle) image.Image {
			return &testImageDefault{image.NewRGBA(r)}
		},
	} {
		p := newImage(r)
		NewDrawable(p)
	}
}

type testImageDefault struct {
	*image.RGBA
}

func TestCopy(t *testing.T) {
	type TC struct {
		srcSize image.Rectangle
		dstSize image.Rectangle
	}
	for _, tc := range []TC{
		{
			srcSize: image.Rect(0, 0, 10, 10),
			dstSize: image.Rect(0, 0, 10, 10),
		},
		{
			srcSize: image.Rect(0, 0, 5, 5),
			dstSize: image.Rect(0, 0, 10, 10),
		},
		{
			srcSize: image.Rect(0, 0, 10, 10),
			dstSize: image.Rect(0, 0, 5, 5),
		},
		{
			srcSize: image.Rect(0, 0, 5, 5),
			dstSize: image.Rect(5, 5, 10, 10),
		},
		{
			srcSize: image.Rect(0, 0, 10, 10),
			dstSize: image.Rect(5, 5, 15, 15),
		},
		{
			srcSize: image.Rect(5, 5, 15, 15),
			dstSize: image.Rect(0, 0, 10, 10),
		},
	} {
		src := image.NewRGBA(tc.srcSize)
		testDrawRandom(src)
		dst := image.NewRGBA(tc.dstSize)
		Copy(dst, src)
		bd := src.Bounds().Intersect(dst.Bounds())
		for y, yEnd := bd.Min.Y, bd.Max.Y; y < yEnd; y++ {
			for x, xEnd := bd.Min.X, bd.Max.X; x < xEnd; x++ {
				cSrc := src.At(x, y)
				cDst := dst.At(x, y)
				if cSrc != cDst {
					t.Fatalf("different color: %#v, pixel %dx%d: src=%#v, dst=%#v", tc, x, y, cSrc, cDst)
				}
			}
		}
	}
}

func testRandomColor() color.Color {
	return color.NRGBA{
		R: uint8(rand.Intn(1 << 8)),
		G: uint8(rand.Intn(1 << 8)),
		B: uint8(rand.Intn(1 << 8)),
		A: uint8(rand.Intn(1 << 8)),
	}
}

func testDrawRandom(p draw.Image) {
	bd := p.Bounds()
	for y, yEnd := bd.Min.Y, bd.Max.Y; y < yEnd; y++ {
		for x, xEnd := bd.Min.X, bd.Max.X; x < xEnd; x++ {
			p.Set(x, y, testRandomColor())
		}
	}
}
