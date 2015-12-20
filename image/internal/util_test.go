package internal

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestRGBAToNRGBA(t *testing.T) {
	test := func(r, g, b, a uint16) {
		r1, g1, b1, a1 := RGBAToNRGBA(uint32(r), uint32(g), uint32(b), uint32(a))
		c := color.NRGBA64Model.Convert(color.RGBA64{r, g, b, a}).(color.NRGBA64)
		r2, g2, b2, a2 := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
		if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
			t.Fatalf("different color: {%d %d %d %d}: got {%d %d %d %d}, want {%d %d %d %d}", r, g, b, a, r1, g1, b1, a1, r2, g2, b2, a2)
		}
	}
	vals := []uint16{0, 0x4000, 0x8000, 0xc000, 0xffff}
	for i, a := range vals {
		for _, r := range vals[:i+1] {
			for _, g := range vals[:i+1] {
				for _, b := range vals[:i+1] {
					test(r, g, b, a)
				}
			}
		}
	}
	for i := 0; i < 1000; i++ {
		a := uint16(rand.Int31n(1 << 16))
		r := uint16(rand.Int31n(int32(a) + 1))
		g := uint16(rand.Int31n(int32(a) + 1))
		b := uint16(rand.Int31n(int32(a) + 1))
		test(r, g, b, a)
	}
}

func TestNRGBAToRGBA(t *testing.T) {
	test := func(r, g, b, a uint16) {
		r1, g1, b1, a1 := NRGBAToRGBA(uint32(r), uint32(g), uint32(b), uint32(a))
		c := color.RGBA64Model.Convert(color.NRGBA64{r, g, b, a}).(color.RGBA64)
		r2, g2, b2, a2 := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
		if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
			t.Fatalf("different color: {%d %d %d %d}: got {%d %d %d %d}, want {%d %d %d %d}", r, g, b, a, r1, g1, b1, a1, r2, g2, b2, a2)
		}
	}
	vals := []uint16{0, 0x4000, 0x8000, 0xc000, 0xffff}
	for _, r := range vals {
		for _, g := range vals {
			for _, b := range vals {
				for _, a := range vals {
					test(r, g, b, a)
				}
			}
		}
	}
	for i := 0; i < 1000; i++ {
		r := uint16(rand.Int31n(1 << 16))
		g := uint16(rand.Int31n(1 << 16))
		b := uint16(rand.Int31n(1 << 16))
		a := uint16(rand.Int31n(1 << 16))
		test(r, g, b, a)
	}
}

func TestNewDrawable(t *testing.T) {
	r := image.Rect(0, 0, 1, 1)
	for _, newImage := range testNewImageFuncs {
		p := newImage(r)
		NewDrawable(p)
	}
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

func TestParallel(t *testing.T) {
	r := image.Rect(100, 100, 200, 200)
	Parallel(r, func(sub image.Rectangle) {
		if !sub.In(r) {
			t.Fatalf("%s is not in %s", sub, r)
		}
	})
}

type testImageDefault struct {
	*image.RGBA
}

var testNewImageFuncs = []func(image.Rectangle) image.Image{
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
		return image.NewPaletted(r, testPalette)
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
}

var testNewImageDrawFuncs = []func(image.Rectangle) draw.Image{
	func(r image.Rectangle) draw.Image {
		return image.NewRGBA(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewRGBA64(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewNRGBA(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewNRGBA64(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewAlpha(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewAlpha16(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewGray(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewGray16(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewCMYK(r)
	},
	func(r image.Rectangle) draw.Image {
		return image.NewPaletted(r, testPalette)
	},
	func(r image.Rectangle) draw.Image {
		return &testImageDefault{image.NewRGBA(r)}
	},
}

var testColors []color.Color

func init() {
	vals := []uint8{0x00, 0x40, 0x80, 0xc0, 0xff}
	for _, r := range vals {
		for _, g := range vals {
			for _, b := range vals {
				for _, a := range vals {
					testColors = append(testColors, color.NRGBA{r, g, b, a})
				}
			}
		}
	}
	for i := 0; i < 100; i++ {
		testColors = append(testColors, testRandomColor())
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

var testPalette = color.Palette{
	color.RGBA{0, 0, 0, 255},
	color.RGBA{255, 0, 0, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{0, 0, 255, 255},
	color.RGBA{255, 255, 255, 255},
}
