// Package internal provides utilities functions used in the image package.
package internal

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/pierrre/imageutil"
)

// NewDrawable returns a new draw.Image with the same type and size as p.
//
// If p has no size, 1x1 is used.
//
// See NewDrawableSize.
func NewDrawable(p image.Image) draw.Image {
	r := p.Bounds()
	if _, ok := p.(*image.Uniform); ok {
		r = image.Rect(0, 0, 1, 1)
	}
	return NewDrawableSize(p, r)
}

// NewDrawableSize returns a new draw.Image with the same type as p and the given bounds.
//
// If p is not a draw.Image, another type is used.
//
// nolint: gocyclo
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
	at := imageutil.NewAtFunc(src)
	set := imageutil.NewSetFunc(dst)
	imageutil.Parallel1D(bd, func(bd image.Rectangle) {
		for y := bd.Min.Y; y < bd.Max.Y; y++ {
			for x := bd.Min.X; x < bd.Max.X; x++ {
				r, g, b, a := at(x, y)
				set(x, y, r, g, b, a)
			}
		}
	})
}
