package internal

import "image/color"

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
