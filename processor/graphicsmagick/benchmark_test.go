package graphicsmagick

import (
	"sync"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

// Worker 1

func BenchmarkResizeWorker1Small(b *testing.B) {
	benchmarkResizeWorker1(b, testdata.Small)
}

func BenchmarkResizeWorker1Medium(b *testing.B) {
	benchmarkResizeWorker1(b, testdata.Medium)
}

func BenchmarkResizeWorker1Large(b *testing.B) {
	benchmarkResizeWorker1(b, testdata.Large)
}

func BenchmarkResizeWorker1Huge(b *testing.B) {
	benchmarkResizeWorker1(b, testdata.Huge)
}

func BenchmarkResizeWorker1Animated(b *testing.B) {
	benchmarkResizeWorker1(b, testdata.Animated)
}

func benchmarkResizeWorker1(b *testing.B, image *imageserver.Image) {
	benchmarkResize(b, 1, image)
}

// Worker 2

func BenchmarkResizeWorker2Small(b *testing.B) {
	benchmarkResizeWorker2(b, testdata.Small)
}

func BenchmarkResizeWorker2Medium(b *testing.B) {
	benchmarkResizeWorker2(b, testdata.Medium)
}

func BenchmarkResizeWorker2Large(b *testing.B) {
	benchmarkResizeWorker2(b, testdata.Large)
}

func BenchmarkResizeWorker2Huge(b *testing.B) {
	benchmarkResizeWorker2(b, testdata.Huge)
}

func BenchmarkResizeWorker2Animated(b *testing.B) {
	benchmarkResizeWorker2(b, testdata.Animated)
}

func benchmarkResizeWorker2(b *testing.B, image *imageserver.Image) {
	benchmarkResize(b, 2, image)
}

// Worker 4

func BenchmarkResizeWorker4Small(b *testing.B) {
	benchmarkResizeWorker4(b, testdata.Small)
}

func BenchmarkResizeWorker4Medium(b *testing.B) {
	benchmarkResizeWorker4(b, testdata.Medium)
}

func BenchmarkResizeWorker4Large(b *testing.B) {
	benchmarkResizeWorker4(b, testdata.Large)
}

func BenchmarkResizeWorker4Huge(b *testing.B) {
	benchmarkResizeWorker4(b, testdata.Huge)
}

func BenchmarkResizeWorker4Animated(b *testing.B) {
	benchmarkResizeWorker4(b, testdata.Animated)
}

func benchmarkResizeWorker4(b *testing.B, image *imageserver.Image) {
	benchmarkResize(b, 4, image)
}

// Worker 8

func BenchmarkResizeWorker8Small(b *testing.B) {
	benchmarkResizeWorker8(b, testdata.Small)
}

func BenchmarkResizeWorker8Medium(b *testing.B) {
	benchmarkResizeWorker8(b, testdata.Medium)
}

func BenchmarkResizeWorker8Large(b *testing.B) {
	benchmarkResizeWorker8(b, testdata.Large)
}

func BenchmarkResizeWorker8Huge(b *testing.B) {
	benchmarkResizeWorker8(b, testdata.Huge)
}

func BenchmarkResizeWorker8Animated(b *testing.B) {
	benchmarkResizeWorker8(b, testdata.Animated)
}

func benchmarkResizeWorker8(b *testing.B, image *imageserver.Image) {
	benchmarkResize(b, 8, image)
}

func benchmarkResize(b *testing.B, workerCount int, image *imageserver.Image) {
	parameters := imageserver.Parameters{
		"graphicsmagick": imageserver.Parameters{
			"width":  100,
			"height": 100,
		},
	}

	processor := &GraphicsMagickProcessor{
		Executable: "gm",
	}

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
				_, err := processor.Process(image, parameters)
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
