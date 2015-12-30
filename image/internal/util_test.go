package internal

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
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

func TestMin(t *testing.T) {
	type TC struct {
		val1     int
		val2     int
		expected int
	}
	for _, tc := range []TC{
		{
			val1:     1,
			val2:     2,
			expected: 1,
		},
		{
			val1:     5,
			val2:     4,
			expected: 4,
		},
		{
			val1:     6,
			val2:     6,
			expected: 6,
		},
	} {
		m := min(tc.val1, tc.val2)
		if m != tc.expected {
			t.Fatalf("val1=%d, val2=%d: got %d, want %d", tc.val1, tc.val2, m, tc.expected)
		}
	}
}

func TestParallel(t *testing.T) {
	type TC struct {
		n        int
		p        int
		expected map[int]int
	}
	for _, tc := range []TC{
		{
			n: 0,
			p: 0,
		},
		{
			n: -1,
			p: 0,
		},
		{
			n: 1,
			p: 4,
			expected: map[int]int{
				0: 1,
			},
		},
		{
			n: 1,
			p: 0,
			expected: map[int]int{
				0: 1,
			},
		},
		{
			n: 8,
			p: 4,
			expected: map[int]int{
				0: 2,
				2: 4,
				4: 6,
				6: 8,
			},
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			var called int32
			parallel(tc.n, tc.p, func(start, end int) {
				expectedEnd, ok := tc.expected[start]
				if !ok {
					t.Fatalf("unexpected start: %d", start)
				}
				if end != expectedEnd {
					t.Fatalf("unexpected end for start %d: got %d want %d", start, end, expectedEnd)
				}
				atomic.AddInt32(&called, 1)
			})
			if int(called) != len(tc.expected) {
				t.Fatalf("unexpected call count: got %d want %d", called, len(tc.expected))
			}
		}()
	}
	Parallel(1, func(start, end int) {})
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
