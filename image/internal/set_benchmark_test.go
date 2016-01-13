package internal

import (
	"image"
	"image/draw"
	"testing"
)

func BenchmarkNewSetFuncRGBA(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewRGBA(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewSetFuncRGBA64(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewRGBA64(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewSetFuncNRGBAOpaque(b *testing.B) {
	benchmarkNewSetFuncColor(b, image.NewNRGBA(image.Rect(0, 0, 1, 1)), 0xffff, 0xffff, 0xffff, 0xffff)
}

func BenchmarkNewSetFuncNRGBATransparent(b *testing.B) {
	benchmarkNewSetFuncColor(b, image.NewNRGBA(image.Rect(0, 0, 1, 1)), 0, 0, 0, 0)
}

func BenchmarkNewSetFuncNRGBATranslucent(b *testing.B) {
	benchmarkNewSetFuncColor(b, image.NewNRGBA(image.Rect(0, 0, 1, 1)), 0x8000, 0x8000, 0x8000, 0x8000)
}

func BenchmarkNewSetFuncNRGBA64Opaque(b *testing.B) {
	benchmarkNewSetFuncColor(b, image.NewNRGBA64(image.Rect(0, 0, 1, 1)), 0xffff, 0xffff, 0xffff, 0xffff)
}

func BenchmarkNewSetFuncNRGBA64Transparent(b *testing.B) {
	benchmarkNewSetFuncColor(b, image.NewNRGBA64(image.Rect(0, 0, 1, 1)), 0, 0, 0, 0)
}

func BenchmarkNewSetFuncNRGBA64Translucent(b *testing.B) {
	benchmarkNewSetFuncColor(b, image.NewNRGBA64(image.Rect(0, 0, 1, 1)), 0x8000, 0x8000, 0x8000, 0x8000)
}

func BenchmarkNewSetFuncAlpha(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewAlpha(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewSetFuncAlpha16(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewAlpha16(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewSetFuncGray(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewGray(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewSetFuncGray16(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewGray16(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewSetFuncPaletted(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewPaletted(image.Rect(0, 0, 1, 1), testPalette))
}

func BenchmarkNewSetFuncCMYK(b *testing.B) {
	benchmarkNewSetFunc(b, image.NewCMYK(image.Rect(0, 0, 1, 1)))
}

func benchmarkNewSetFunc(b *testing.B, p draw.Image) {
	benchmarkNewSetFuncColor(b, p, 0xffff, 0xffff, 0xffff, 0xffff)
}

func BenchmarkNewSetFuncDefault(b *testing.B) {
	benchmarkNewSetFunc(b, &testImageDefault{image.NewRGBA(image.Rect(0, 0, 1, 1))})
}

func benchmarkNewSetFuncColor(b *testing.B, p draw.Image, rr, gg, bb, aa uint32) {
	set := NewSetFunc(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(0, 0, rr, gg, bb, aa)
	}
}
