package internal

import (
	"image/color"
	"math/rand"
	"testing"
)

func TestRGBAToNRGBA(t *testing.T) {
	test := func(r, g, b, a uint16) {
		r1, g1, b1, a1 := RGBAToNRGBA(uint32(r), uint32(g), uint32(b), uint32(a))
		c := color.NRGBA64Model.Convert(color.RGBA64{r, g, b, a}).(color.NRGBA64)
		r2, g2, b2, a2 := uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
		if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
			t.Errorf("different color: {%d %d %d %d}: got {%d %d %d %d}, want {%d %d %d %d}", r, g, b, a, r1, g1, b1, a1, r2, g2, b2, a2)
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
			t.Errorf("different color: {%d %d %d %d}: got {%d %d %d %d}, want {%d %d %d %d}", r, g, b, a, r1, g1, b1, a1, r2, g2, b2, a2)
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
