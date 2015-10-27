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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Copy(dst, src)
		}
	})
}
