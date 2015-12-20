package internal

import "image"

// AtFunc returns a RGBA value of the pixel at (x, y).
type AtFunc func(x, y int) (r, g, b, a uint32)

// NewAtFunc returns an AtFunc for an Image.
func NewAtFunc(p image.Image) AtFunc {
	switch p := p.(type) {
	case *image.RGBA:
		return newAtFuncRGBA(p)
	case *image.RGBA64:
		return newAtFuncRGBA64(p)
	case *image.NRGBA:
		return newAtFuncNRGBA(p)
	case *image.NRGBA64:
		return newAtFuncNRGBA64(p)
	case *image.Alpha:
		return newAtFuncAlpha(p)
	case *image.Alpha16:
		return newAtFuncAlpha16(p)
	case *image.Gray:
		return newAtFuncGray(p)
	case *image.Gray16:
		return newAtFuncGray16(p)
	case *image.Paletted:
		return newAtFuncPaletted(p)
	case *image.YCbCr:
		return newAtFuncYCbCr(p)
	case *image.NYCbCrA:
		return newAtFuncNYCbCrA(p)
	case *image.CMYK:
		return newAtFuncCMYK(p)
	case *image.Uniform:
		return newAtFuncUniform(p)
	default:
		return newAtFuncDefault(p)
	}
}

func newAtFuncRGBA(p *image.RGBA) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
		r = uint32(p.Pix[i+0])
		r |= r << 8
		g = uint32(p.Pix[i+1])
		g |= g << 8
		b = uint32(p.Pix[i+2])
		b |= b << 8
		a = uint32(p.Pix[i+3])
		a |= a << 8
		return
	}
}

func newAtFuncRGBA64(p *image.RGBA64) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
		r = uint32(p.Pix[i+0])<<8 | uint32(p.Pix[i+1])
		g = uint32(p.Pix[i+2])<<8 | uint32(p.Pix[i+3])
		b = uint32(p.Pix[i+4])<<8 | uint32(p.Pix[i+5])
		a = uint32(p.Pix[i+6])<<8 | uint32(p.Pix[i+7])
		return
	}
}

func newAtFuncNRGBA(p *image.NRGBA) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
		a = uint32(p.Pix[i+3])
		a |= a << 8
		if a == 0 {
			return
		}
		r = uint32(p.Pix[i+0])
		r |= r << 8
		g = uint32(p.Pix[i+1])
		g |= g << 8
		b = uint32(p.Pix[i+2])
		b |= b << 8
		if a == 0xffff {
			return
		}
		r = r * a / 0xffff
		g = g * a / 0xffff
		b = b * a / 0xffff
		return
	}
}

func newAtFuncNRGBA64(p *image.NRGBA64) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
		a = uint32(p.Pix[i+6])<<8 | uint32(p.Pix[i+7])
		if a == 0 {
			return
		}
		r = uint32(p.Pix[i+0])<<8 | uint32(p.Pix[i+1])
		g = uint32(p.Pix[i+2])<<8 | uint32(p.Pix[i+3])
		b = uint32(p.Pix[i+4])<<8 | uint32(p.Pix[i+5])
		if a == 0xffff {
			return
		}
		r = r * a / 0xffff
		g = g * a / 0xffff
		b = b * a / 0xffff
		return
	}
}

func newAtFuncAlpha(p *image.Alpha) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
		a = uint32(p.Pix[i])
		a |= a << 8
		return a, a, a, a
	}
}

func newAtFuncAlpha16(p *image.Alpha16) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
		a = uint32(p.Pix[i+0])<<8 | uint32(p.Pix[i+1])
		return a, a, a, a
	}
}

func newAtFuncGray(p *image.Gray) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
		yy := uint32(p.Pix[i])
		yy |= yy << 8
		return yy, yy, yy, 0xffff
	}
}

func newAtFuncGray16(p *image.Gray16) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*2
		yy := uint32(p.Pix[i+0])<<8 | uint32(p.Pix[i+1])
		return yy, yy, yy, 0xffff
	}
}

func newAtFuncPaletted(p *image.Paletted) AtFunc {
	pa := newPaletteRGBA(p.Palette)
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*1
		c := pa[p.Pix[i]]
		return c.r, c.g, c.b, c.a
	}
}

func newAtFuncUniform(p *image.Uniform) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		return p.C.RGBA()
	}
}

func newAtFuncYCbCr(p *image.YCbCr) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		yi := (y-p.Rect.Min.Y)*p.YStride + (x - p.Rect.Min.X)
		var ci int
		switch p.SubsampleRatio {
		case image.YCbCrSubsampleRatio422:
			ci = (y-p.Rect.Min.Y)*p.CStride + (x/2 - p.Rect.Min.X/2)
		case image.YCbCrSubsampleRatio420:
			ci = (y/2-p.Rect.Min.Y/2)*p.CStride + (x/2 - p.Rect.Min.X/2)
		case image.YCbCrSubsampleRatio440:
			ci = (y/2-p.Rect.Min.Y/2)*p.CStride + (x - p.Rect.Min.X)
		case image.YCbCrSubsampleRatio411:
			ci = (y-p.Rect.Min.Y)*p.CStride + (x/4 - p.Rect.Min.X/4)
		case image.YCbCrSubsampleRatio410:
			ci = (y/2-p.Rect.Min.Y/2)*p.CStride + (x/4 - p.Rect.Min.X/4)
		default:
			ci = (y-p.Rect.Min.Y)*p.CStride + (x - p.Rect.Min.X)
		}
		yy1 := int32(p.Y[yi]) * 0x10100
		cb1 := int32(p.Cb[ci]) - 128
		cr1 := int32(p.Cr[ci]) - 128
		r1 := (yy1 + 91881*cr1) >> 8
		g1 := (yy1 - 22554*cb1 - 46802*cr1) >> 8
		b1 := (yy1 + 116130*cb1) >> 8
		if r1 < 0 {
			r1 = 0
		} else if r1 > 0xffff {
			r1 = 0xffff
		}
		if g1 < 0 {
			g1 = 0
		} else if g1 > 0xffff {
			g1 = 0xffff
		}
		if b1 < 0 {
			b1 = 0
		} else if b1 > 0xffff {
			b1 = 0xffff
		}
		return uint32(r1), uint32(g1), uint32(b1), 0xffff
	}
}

func newAtFuncNYCbCrA(p *image.NYCbCrA) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		yi := (y-p.Rect.Min.Y)*p.YStride + (x - p.Rect.Min.X)
		var ci int
		switch p.SubsampleRatio {
		case image.YCbCrSubsampleRatio422:
			ci = (y-p.Rect.Min.Y)*p.CStride + (x/2 - p.Rect.Min.X/2)
		case image.YCbCrSubsampleRatio420:
			ci = (y/2-p.Rect.Min.Y/2)*p.CStride + (x/2 - p.Rect.Min.X/2)
		case image.YCbCrSubsampleRatio440:
			ci = (y/2-p.Rect.Min.Y/2)*p.CStride + (x - p.Rect.Min.X)
		case image.YCbCrSubsampleRatio411:
			ci = (y-p.Rect.Min.Y)*p.CStride + (x/4 - p.Rect.Min.X/4)
		case image.YCbCrSubsampleRatio410:
			ci = (y/2-p.Rect.Min.Y/2)*p.CStride + (x/4 - p.Rect.Min.X/4)
		default:
			ci = (y-p.Rect.Min.Y)*p.CStride + (x - p.Rect.Min.X)
		}
		ai := (y-p.Rect.Min.Y)*p.AStride + (x - p.Rect.Min.X)
		yy1 := int32(p.Y[yi]) * 0x10100
		cb1 := int32(p.Cb[ci]) - 128
		cr1 := int32(p.Cr[ci]) - 128
		r1 := (yy1 + 91881*cr1) >> 8
		g1 := (yy1 - 22554*cb1 - 46802*cr1) >> 8
		b1 := (yy1 + 116130*cb1) >> 8
		if r1 < 0 {
			r1 = 0
		} else if r1 > 0xffff {
			r1 = 0xffff
		}
		if g1 < 0 {
			g1 = 0
		} else if g1 > 0xffff {
			g1 = 0xffff
		}
		if b1 < 0 {
			b1 = 0
		} else if b1 > 0xffff {
			b1 = 0xffff
		}
		a = uint32(p.A[ai]) * 0x101
		r = uint32(r1) * a / 0xffff
		g = uint32(g1) * a / 0xffff
		b = uint32(b1) * a / 0xffff
		return
	}
}

func newAtFuncCMYK(p *image.CMYK) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
		w := uint32(0xffff - uint32(p.Pix[i+3])*0x101)
		r = uint32(0xffff-uint32(p.Pix[i+0])*0x101) * w / 0xffff
		g = uint32(0xffff-uint32(p.Pix[i+1])*0x101) * w / 0xffff
		b = uint32(0xffff-uint32(p.Pix[i+2])*0x101) * w / 0xffff
		a = 0xffff
		return
	}
}

func newAtFuncDefault(p image.Image) AtFunc {
	return func(x, y int) (r, g, b, a uint32) {
		return p.At(x, y).RGBA()
	}
}
