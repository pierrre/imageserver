package internal

import (
	"image"
	"image/color"
	"testing"
)

func TestNewSetFunc(t *testing.T) {
	bd := image.Rect(0, 0, 3, 3)
	for _, newImageDrawFunc := range testNewImageDrawFuncs {
		p := newImageDrawFunc(bd)
		set := NewSetFunc(p)
		for _, c := range testColors {
			c := color.RGBA64Model.Convert(c).(color.RGBA64)
			for y := bd.Min.Y; y < bd.Max.Y; y++ {
				for x := bd.Min.X; x < bd.Max.X; x++ {
					set(x, y, uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A))
					r1, g1, b1, a1 := p.At(x, y).RGBA()
					p.Set(x, y, c)
					r2, g2, b2, a2 := p.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("different color: image %T, pixel %dx%d, color %#v: got {%d %d %d %d}, want {%d %d %d %d}", p, x, y, c, r1, g1, b1, a1, r2, g2, b2, a2)
					}
				}
			}
		}
	}
}
