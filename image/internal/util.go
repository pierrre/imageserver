package internal

import (
	"image"
	"image/draw"
	"runtime"
	"sync"
)

// NewDrawable returns a new draw.Image with the same type and size as p.
// If p has no size, 1x1 is used.
// See NewDrawableSize.
func NewDrawable(p image.Image) draw.Image {
	width := p.Bounds().Dx()
	height := p.Bounds().Dy()
	if _, ok := p.(*image.Uniform); ok {
		width = 1
		height = 1
	}
	return NewDrawableSize(p, width, height)
}

// NewDrawableSize returns a new draw.Image with the same type as p and the given size .
// If p is not a draw.Image, another type is used.
func NewDrawableSize(p image.Image, width, height int) draw.Image {
	r := image.Rect(0, 0, width, height)
	switch p.(type) {
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
	case *image.CMYK:
		return image.NewCMYK(r)
	default:
		return image.NewRGBA(r)
	}
}

// Copy copies src to dst.
func Copy(dst draw.Image, src image.Image) {
	width := src.Bounds().Dx()
	if dstWidth := dst.Bounds().Dx(); dstWidth < width {
		width = dstWidth
	}
	height := src.Bounds().Dy()
	if dstHeight := dst.Bounds().Dy(); dstHeight < height {
		height = dstHeight
	}
	at := NewAtFunc(src)
	set := NewSetFunc(dst)
	Parallel(height, func(yStart, yEnd int) {
		for y := yStart; y < yEnd; y++ {
			for x := 0; x < width; x++ {
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
