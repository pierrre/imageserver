package internal

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
