package internal

import (
	"image"
	"image/color"
	"testing"
)

func BenchmarkNewAtFuncRGBA(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewRGBA(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewAtFuncRGBA64(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewRGBA64(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewAtFuncNRGBAOpaque(b *testing.B) {
	p := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	p.SetNRGBA(0, 0, color.NRGBA{255, 255, 255, 255})
	benchmarkNewAtFunc(b, p)
}

func BenchmarkNewAtFuncNRGBATransparent(b *testing.B) {
	p := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	p.SetNRGBA(0, 0, color.NRGBA{255, 255, 255, 0})
	benchmarkNewAtFunc(b, p)
}

func BenchmarkNewAtFuncNRGBATranslucent(b *testing.B) {
	p := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	p.SetNRGBA(0, 0, color.NRGBA{255, 255, 255, 128})
	benchmarkNewAtFunc(b, p)
}

func BenchmarkNewAtFuncNRGBA64Opaque(b *testing.B) {
	p := image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	p.SetNRGBA64(0, 0, color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff})
	benchmarkNewAtFunc(b, p)
}

func BenchmarkNewAtFuncNRGBA64Transparent(b *testing.B) {
	p := image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	p.SetNRGBA64(0, 0, color.NRGBA64{0xffff, 0xffff, 0xffff, 0})
	benchmarkNewAtFunc(b, p)
}

func BenchmarkNewAtFuncNRGBA64Translucent(b *testing.B) {
	p := image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	p.SetNRGBA64(0, 0, color.NRGBA64{0xffff, 0xffff, 0xffff, 0x8000})
	benchmarkNewAtFunc(b, p)
}

func BenchmarkNewAtFuncAlpha(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewAlpha(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewAtFuncAlpha16(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewAlpha16(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewAtFuncGray(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewGray(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewAtFuncGray16(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewGray16(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewAtFuncPaletted(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewPaletted(image.Rect(0, 0, 1, 1), testPalette))
}

func BenchmarkNewAtFuncUniform(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewUniform(color.RGBA{}))
}

func BenchmarkNewAtFuncYCbCr(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewYCbCr(image.Rect(0, 0, 1, 1), image.YCbCrSubsampleRatio444))
}

func BenchmarkNewAtFuncNYCbCrA(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewNYCbCrA(image.Rect(0, 0, 1, 1), image.YCbCrSubsampleRatio444))
}

func BenchmarkNewAtFuncCMYK(b *testing.B) {
	benchmarkNewAtFunc(b, image.NewCMYK(image.Rect(0, 0, 1, 1)))
}

func BenchmarkNewAtFuncDefault(b *testing.B) {
	benchmarkNewAtFunc(b, &testImageDefault{image.NewRGBA(image.Rect(0, 0, 1, 1))})
}

func benchmarkNewAtFunc(b *testing.B, p image.Image) {
	at := NewAtFunc(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		at(0, 0)
	}
}
