package redis

import (
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"sync"
	"testing"
)

// Worker 1

func BenchmarkGetSmallWorker1(b *testing.B) {
	benchmarkGet(b, testdata.Small, 1)
}

func BenchmarkGetMediumWorker1(b *testing.B) {
	benchmarkGet(b, testdata.Medium, 1)
}

func BenchmarkGetLargeWorker1(b *testing.B) {
	benchmarkGet(b, testdata.Large, 1)
}

func BenchmarkGetHugeWorker1(b *testing.B) {
	benchmarkGet(b, testdata.Huge, 1)
}

func BenchmarkGetAnimatedWorker1(b *testing.B) {
	benchmarkGet(b, testdata.Animated, 1)
}

// Worker 2

func BenchmarkGetSmallWorker2(b *testing.B) {
	benchmarkGet(b, testdata.Small, 2)
}

func BenchmarkGetMediumWorker2(b *testing.B) {
	benchmarkGet(b, testdata.Medium, 2)
}

func BenchmarkGetLargeWorker2(b *testing.B) {
	benchmarkGet(b, testdata.Large, 2)
}

func BenchmarkGetHugeWorker2(b *testing.B) {
	benchmarkGet(b, testdata.Huge, 2)
}

func BenchmarkGetAnimatedWorker2(b *testing.B) {
	benchmarkGet(b, testdata.Animated, 2)
}

// Worker 4

func BenchmarkGetSmallWorker4(b *testing.B) {
	benchmarkGet(b, testdata.Small, 4)
}

func BenchmarkGetMediumWorker4(b *testing.B) {
	benchmarkGet(b, testdata.Medium, 4)
}

func BenchmarkGetLargeWorker4(b *testing.B) {
	benchmarkGet(b, testdata.Large, 4)
}

func BenchmarkGetHugeWorker4(b *testing.B) {
	benchmarkGet(b, testdata.Huge, 4)
}

func BenchmarkGetAnimatedWorker4(b *testing.B) {
	benchmarkGet(b, testdata.Animated, 4)
}

// Worker 8

func BenchmarkGetSmallWorker8(b *testing.B) {
	benchmarkGet(b, testdata.Small, 8)
}

func BenchmarkGetMediumWorker8(b *testing.B) {
	benchmarkGet(b, testdata.Medium, 8)
}

func BenchmarkGetLargeWorker8(b *testing.B) {
	benchmarkGet(b, testdata.Large, 8)
}

func BenchmarkGetHugeWorker8(b *testing.B) {
	benchmarkGet(b, testdata.Huge, 8)
}

func BenchmarkGetAnimatedWorker8(b *testing.B) {
	benchmarkGet(b, testdata.Animated, 8)
}

// Worker 16

func BenchmarkGetSmallWorker16(b *testing.B) {
	benchmarkGet(b, testdata.Small, 16)
}

func BenchmarkGetMediumWorker16(b *testing.B) {
	benchmarkGet(b, testdata.Medium, 16)
}

func BenchmarkGetLargeWorker16(b *testing.B) {
	benchmarkGet(b, testdata.Large, 16)
}

func BenchmarkGetHugeWorker16(b *testing.B) {
	benchmarkGet(b, testdata.Huge, 16)
}

func BenchmarkGetAnimatedWorker16(b *testing.B) {
	benchmarkGet(b, testdata.Animated, 16)
}

func benchmarkGet(b *testing.B, image *imageserver.Image, workerCount int) {
	pool := &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", "localhost:6379")
		},
		MaxIdle: 50,
	}
	defer pool.Close()

	cache := &RedisCache{
		Pool: pool,
	}

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
