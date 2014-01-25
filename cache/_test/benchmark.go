package _test

import (
	"github.com/pierrre/imageserver"
	"sync"
	"testing"
)

// CacheBenchmarkGet is a helper to benchmark cache Get()
func CacheBenchmarkGet(b *testing.B, cache imageserver.Cache, workerCount int, image *imageserver.Image) {
	key := "test"
	parameters := make(imageserver.Parameters)
	err := cache.Set(key, image, parameters)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	work := make(chan bool)
	go func() {
		for i := 0; i < b.N; i++ {
			work <- true
		}
		close(work)
	}()
	wg := new(sync.WaitGroup)
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			for _ = range work {
				_, err := cache.Get(key, parameters)
				if err != nil {
					b.Fatal(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	b.SetBytes(int64(len(image.Data)))
}
