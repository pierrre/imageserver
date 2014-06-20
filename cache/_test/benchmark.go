package _test

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// CacheBenchmarkGet is a helper to benchmark cache Get()
func CacheBenchmarkGet(b *testing.B, cache imageserver_cache.Cache, image *imageserver.Image) {
	key := "test"
	parameters := make(imageserver.Parameters)
	err := cache.Set(key, image, parameters)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := cache.Get(key, parameters)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.SetBytes(int64(len(image.Data)))
}
