package _test

import (
	"context"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// BenchmarkGet is a helper to benchmark imageserver/cache.Cache.Get().
func BenchmarkGet(b *testing.B, cache imageserver_cache.Cache, parallelism int, im *imageserver.Image) {
	ctx := context.Background()
	key := "test"
	err := cache.Set(ctx, key, im, imageserver.Params{})
	if err != nil {
		b.Fatal(err)
	}
	params := imageserver.Params{}
	b.SetParallelism(parallelism)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			im, err := cache.Get(ctx, key, params)
			if err != nil {
				b.Fatal(err)
			}
			if im == nil {
				b.Fatal("image nil")
			}
		}
	})
	b.SetBytes(int64(len(im.Data)))
}
