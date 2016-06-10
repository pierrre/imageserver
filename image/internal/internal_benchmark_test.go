package internal

import (
	"context"
	"image"
	"testing"
)

func BenchmarkCopy(b *testing.B) {
	ctx := context.Background()
	src := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	testDrawRandom(src)
	dst := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(ctx, dst, src)
	}
}
