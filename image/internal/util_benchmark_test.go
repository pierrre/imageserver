package internal

import (
	"image"
	"testing"
)

func BenchmarkRGBAToNRGBAOpaque(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RGBAToNRGBA(0xffff, 0x8000, 0x0000, 0xffff)
		}
	})
}

func BenchmarkRGBAToNRGBATransparent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RGBAToNRGBA(0x0000, 0x0000, 0x0000, 0x0000)
		}
	})
}

func BenchmarkRGBAToNRGBATranslucent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RGBAToNRGBA(0x8000, 0x4000, 0x0000, 0x8000)
		}
	})
}

func BenchmarkNRGBAToRGBAOpaque(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0xffff)
		}
	})
}

func BenchmarkNRGBAToRGBATransparent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0x0000)
		}
	})
}

func BenchmarkNRGBAToRGBATranslucent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0x8000)
		}
	})
}

func BenchmarkCopy(b *testing.B) {
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	testDrawRandom(src)
	dst := image.NewRGBA(image.Rect(0, 0, 100, 100))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Copy(dst, src)
		}
	})
}
