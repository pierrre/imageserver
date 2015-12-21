package internal

import (
	"image"
	"image/color"
	"image/draw"
)

// SetFunc sets a RGBA value to the pixel at (x, y).
type SetFunc func(x, y int, r, g, b, a uint32)

// NewSetFunc returns a SetFunc for an Image.
func NewSetFunc(p draw.Image) SetFunc {
	switch p := p.(type) {
	case *image.RGBA:
		return newSetFuncRGBA(p)
	case *image.RGBA64:
		return newSetFuncRGBA64(p)
	case *image.NRGBA:
		return newSetFuncNRGBA(p)
	case *image.NRGBA64:
		return newSetFuncNRGBA64(p)
	case *image.Alpha:
		return newSetFuncAlpha(p)
	case *image.Alpha16:
		return newSetFuncAlpha16(p)
	case *image.Gray:
		return newSetFuncGray(p)
	case *image.Gray16:
		return newSetFuncGray16(p)
	case *image.Paletted:
		return newSetFuncPaletted(p)
	case *image.CMYK:
		return newSetFuncCMYK(p)
	default:
		return newSetFuncDefault(p)
	}
}

func newSetFuncRGBA(p *image.RGBA) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
		p.Pix[i+0] = uint8(r >> 8)
		p.Pix[i+1] = uint8(g >> 8)
		p.Pix[i+2] = uint8(b >> 8)
		p.Pix[i+3] = uint8(a >> 8)
	}
}

func newSetFuncRGBA64(p *image.RGBA64) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
		r16, g16, b16, a16 := uint16(r), uint16(g), uint16(b), uint16(a)
		p.Pix[i+0] = uint8(r16 >> 8)
		p.Pix[i+1] = uint8(r16)
		p.Pix[i+2] = uint8(g16 >> 8)
		p.Pix[i+3] = uint8(g16)
		p.Pix[i+4] = uint8(b16 >> 8)
		p.Pix[i+5] = uint8(b16)
		p.Pix[i+6] = uint8(a16 >> 8)
		p.Pix[i+7] = uint8(a16)
	}
}

func newSetFuncNRGBA(p *image.NRGBA) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		r, g, b, a = RGBAToNRGBA(r, g, b, a)
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
		p.Pix[i+0] = uint8(r >> 8)
		p.Pix[i+1] = uint8(g >> 8)
		p.Pix[i+2] = uint8(b >> 8)
		p.Pix[i+3] = uint8(a >> 8)
	}
}

func newSetFuncNRGBA64(p *image.NRGBA64) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		r, g, b, a = RGBAToNRGBA(r, g, b, a)
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
		p.Pix[i+0] = uint8(r >> 8)
		p.Pix[i+1] = uint8(r)
		p.Pix[i+2] = uint8(g >> 8)
		p.Pix[i+3] = uint8(g)
		p.Pix[i+4] = uint8(b >> 8)
		p.Pix[i+5] = uint8(b)
		p.Pix[i+6] = uint8(a >> 8)
		p.Pix[i+7] = uint8(a)
	}
}

func newSetFuncAlpha(p *image.Alpha) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
		p.Pix[i] = uint8(a >> 8)
	}
}

func newSetFuncAlpha16(p *image.Alpha16) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
		a16 := uint16(a)
		p.Pix[i+0] = uint8(a16 >> 8)
		p.Pix[i+1] = uint8(a16)
	}
}

func newSetFuncGray(p *image.Gray) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
		p.Pix[i] = uint8(((299*r + 587*g + 114*b + 500) / 1000) >> 8)
	}
}

func newSetFuncGray16(p *image.Gray16) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
		y16 := uint16((299*r + 587*g + 114*b + 500) / 1000)
		p.Pix[i+0] = uint8(y16 >> 8)
		p.Pix[i+1] = uint8(y16)
	}
}

func newSetFuncPaletted(p *image.Paletted) SetFunc {
	pa := newPaletteRGBA(p.Palette)
	return func(x, y int, r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
		p.Pix[i] = uint8(pa.index(colorRGBA{r, g, b, a}))
	}
}

func newSetFuncCMYK(p *image.CMYK) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		rr := uint32(r >> 8)
		gg := uint32(g >> 8)
		bb := uint32(b >> 8)
		w := rr
		if w < gg {
			w = gg
		}
		if w < bb {
			w = bb
		}
		var c8, m8, y8, k8 uint8
		if w == 0 {
			k8 = 0xff
		} else {
			c8 = uint8((w - rr) * 0xff / w)
			m8 = uint8((w - gg) * 0xff / w)
			y8 = uint8((w - bb) * 0xff / w)
			k8 = uint8(0xff - w)
		}
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
		p.Pix[i+0] = c8
		p.Pix[i+1] = m8
		p.Pix[i+2] = y8
		p.Pix[i+3] = k8
	}
}

func newSetFuncDefault(p draw.Image) SetFunc {
	return func(x, y int, r, g, b, a uint32) {
		p.Set(x, y, color.RGBA64{
			R: uint16(r),
			G: uint16(g),
			B: uint16(b),
			A: uint16(a),
		})
	}
}
