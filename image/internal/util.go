package internal

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"
	"sync"
)

// RGBAToNRGBA converts RGBA to NRGBA.
func RGBAToNRGBA(r, g, b, a uint32) (uint32, uint32, uint32, uint32) {
	if a == 0xffff {
		return r, g, b, a
	}
	if a == 0 {
		return 0, 0, 0, 0
	}
	r = r * 0xffff / a
	g = g * 0xffff / a
	b = b * 0xffff / a
	return r, g, b, a
}

// NRGBAToRGBA converts NRGBA to RGBA.
func NRGBAToRGBA(r, g, b, a uint32) (uint32, uint32, uint32, uint32) {
	if a == 0xffff {
		return r, g, b, a
	}
	if a == 0 {
		return 0, 0, 0, 0
	}
	r = r * a / 0xffff
	g = g * a / 0xffff
	b = b * a / 0xffff
	return r, g, b, a
}

// NewDrawable returns a new draw.Image with the same type and size as p.
// If p has no size, 1x1 is used.
// See NewDrawableSize.
func NewDrawable(p image.Image) draw.Image {
	r := p.Bounds()
	if _, ok := p.(*image.Uniform); ok {
		r = image.Rect(0, 0, 1, 1)
	}
	return NewDrawableSize(p, r)
}

// NewDrawableSize returns a new draw.Image with the same type as p and the given bounds.
// If p is not a draw.Image, another type is used.
func NewDrawableSize(p image.Image, r image.Rectangle) draw.Image {
	switch p := p.(type) {
	case *image.RGBA:
		return image.NewRGBA(r)
	case *image.RGBA64:
		return image.NewRGBA64(r)
	case *image.NRGBA:
		return image.NewNRGBA(r)
	case *image.NRGBA64:
		return image.NewNRGBA64(r)
	case *image.Alpha:
		return image.NewAlpha(r)
	case *image.Alpha16:
		return image.NewAlpha16(r)
	case *image.Gray:
		return image.NewGray(r)
	case *image.Gray16:
		return image.NewGray16(r)
	case *image.Paletted:
		pl := make(color.Palette, len(p.Palette))
		copy(pl, p.Palette)
		return image.NewPaletted(r, pl)
	case *image.CMYK:
		return image.NewCMYK(r)
	default:
		return image.NewRGBA(r)
	}
}

// Copy copies src to dst.
func Copy(dst draw.Image, src image.Image) {
	bd := src.Bounds().Intersect(dst.Bounds())
	at := NewAtFunc(src)
	set := NewSetFunc(dst)
	Parallel(bd.Dy(), func(yOffStart, yOffEnd int) {
		for y, yEnd := bd.Min.Y+yOffStart, bd.Min.Y+yOffEnd; y < yEnd; y++ {
			for x, xEnd := bd.Min.X, bd.Max.X; x < xEnd; x++ {
				r, g, b, a := at(x, y)
				set(x, y, r, g, b, a)
			}
		}
	})
}

// Parallel helps to dispatch tasks concurrently.
// It calls f with arguments (0,a) (a,b) ... (x,n).
// Currently, it starts GOMAXPROCS goroutines.
func Parallel(n int, f func(start, end int)) {
	parallel(n, runtime.GOMAXPROCS(0), f)
}

func parallel(n int, p int, f func(start, end int)) {
	if n < 1 {
		return
	}
	// n >= 1
	if p > n {
		p = n
	} else if p < 1 {
		p = 1
	}
	// n >= p >= 1
	if p == 1 {
		f(0, n)
		return
	}
	// n >= p > 1
	wg := new(sync.WaitGroup)
	wg.Add(p)
	for i := 0; i < p; i++ {
		go func(i int) {
			defer wg.Done()
			start := n * i / p
			end := n * (i + 1) / p
			f(start, end)
		}(i)
	}
	wg.Wait()
}

type colorRGBA struct {
	r, g, b, a uint32
}

type paletteRGBA []colorRGBA

func newPaletteRGBA(pl color.Palette) paletteRGBA {
	pa := make(paletteRGBA, len(pl))
	for i, c := range pl {
		r, g, b, a := c.RGBA()
		pa[i] = colorRGBA{r, g, b, a}
	}
	return pa
}

func (pa paletteRGBA) index(c colorRGBA) int {
	ret, bestSum := 0, uint32(1<<32-1)
	for i, ca := range pa {
		sum := sqDiff(c.r, ca.r) + sqDiff(c.g, ca.g) + sqDiff(c.b, ca.b) + sqDiff(c.a, ca.a)
		if sum < bestSum {
			if sum == 0 {
				return i
			}
			ret, bestSum = i, sum
		}
	}
	return ret
}

func sqDiff(x, y uint32) uint32 {
	var d uint32
	if x > y {
		d = x - y
	} else {
		d = y - x
	}
	return (d * d) >> 2
}
