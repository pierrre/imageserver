package internal

import "testing"

func BenchmarkRGBAToNRGBAOpaque(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RGBAToNRGBA(0xffff, 0x8000, 0x0000, 0xffff)
	}
}

func BenchmarkRGBAToNRGBATransparent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RGBAToNRGBA(0x0000, 0x0000, 0x0000, 0x0000)
	}
}

func BenchmarkRGBAToNRGBATranslucent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RGBAToNRGBA(0x8000, 0x4000, 0x0000, 0x8000)
	}
}

func BenchmarkNRGBAToRGBAOpaque(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0xffff)
	}
}

func BenchmarkNRGBAToRGBATransparent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0x0000)
	}
}

func BenchmarkNRGBAToRGBATranslucent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0x8000)
	}
}
