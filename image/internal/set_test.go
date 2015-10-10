package internal

import (
	"image"
	"image/color"
	"testing"
)

func TestNewSetFunc(t *testing.T) {
	width := 3
	height := 3
	r := image.Rect(0, 0, width, height)
	for _, newImageDrawFunc := range testNewImageDrawFuncs {
		p := newImageDrawFunc(r)
		set := NewSetFunc(p)
		for _, c := range testColors {
			c := color.RGBA64Model.Convert(c).(color.RGBA64)
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					set(x, y, uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A))
					r1, g1, b1, a1 := p.At(x, y).RGBA()
					p.Set(x, y, c)
					r2, g2, b2, a2 := p.At(x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Errorf("different color: image %T, pixel %dx%d, color %#v: got {%d %d %d %d}, want {%d %d %d %d}", p, x, y, c, r1, g1, b1, a1, r2, g2, b2, a2)
					}
				}
			}
		}
	}
}
