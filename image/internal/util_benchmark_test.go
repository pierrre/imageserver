package internal

import (
	"image"
	"testing"
)

var benchResR, benchResG, benchResB, benchResA uint32

func BenchmarkRGBAToNRGBAOpaque(b *testing.B) {
	var resR, resG, resB, resA uint32
	for i := 0; i < b.N; i++ {
		resR, resG, resB, resA = RGBAToNRGBA(0xffff, 0x8000, 0x0000, 0xffff)
	}
	benchResR, benchResG, benchResB, benchResA = resR, resG, resB, resA
}

func BenchmarkRGBAToNRGBATransparent(b *testing.B) {
	var resR, resG, resB, resA uint32
	for i := 0; i < b.N; i++ {
		resR, resG, resB, resA = RGBAToNRGBA(0x0000, 0x0000, 0x0000, 0x0000)
	}
	benchResR, benchResG, benchResB, benchResA = resR, resG, resB, resA
}

func BenchmarkRGBAToNRGBATranslucent(b *testing.B) {
	var resR, resG, resB, resA uint32
	for i := 0; i < b.N; i++ {
		resR, resG, resB, resA = RGBAToNRGBA(0x8000, 0x4000, 0x0000, 0x8000)
	}
	benchResR, benchResG, benchResB, benchResA = resR, resG, resB, resA
}

func BenchmarkNRGBAToRGBAOpaque(b *testing.B) {
	var resR, resG, resB, resA uint32
	for i := 0; i < b.N; i++ {
		resR, resG, resB, resA = NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0xffff)
	}
	benchResR, benchResG, benchResB, benchResA = resR, resG, resB, resA
}

func BenchmarkNRGBAToRGBATransparent(b *testing.B) {
	var resR, resG, resB, resA uint32
	for i := 0; i < b.N; i++ {
		resR, resG, resB, resA = NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0x0000)
	}
	benchResR, benchResG, benchResB, benchResA = resR, resG, resB, resA
}

func BenchmarkNRGBAToRGBATranslucent(b *testing.B) {
	var resR, resG, resB, resA uint32
	for i := 0; i < b.N; i++ {
		resR, resG, resB, resA = NRGBAToRGBA(0xffff, 0x8000, 0x0000, 0x8000)
	}
	benchResR, benchResG, benchResB, benchResA = resR, resG, resB, resA
}

func BenchmarkCopy(b *testing.B) {
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	testDrawRandom(src)
	dst := image.NewRGBA(image.Rect(0, 0, 100, 100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(dst, src)
	}
}
