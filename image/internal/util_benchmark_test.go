package internal

import (
	"image"
	"testing"
)

func BenchmarkCopy(b *testing.B) {
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	testDrawRandom(src)
	dst := image.NewRGBA(image.Rect(0, 0, 100, 100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(dst, src)
	}
}
