package graphicsmagick

import (
	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"sync"
	"testing"
)

// Worker 1

func BenchmarkResizeSmallWorker1(b *testing.B) {
	benchmarkResize(b, testdata.Small, 1)
}

func BenchmarkResizeMediumWorker1(b *testing.B) {
	benchmarkResize(b, testdata.Medium, 1)
}

func BenchmarkResizeLargeWorker1(b *testing.B) {
	benchmarkResize(b, testdata.Large, 1)
}

func BenchmarkResizeHugeWorker1(b *testing.B) {
	benchmarkResize(b, testdata.Huge, 1)
}

func BenchmarkResizeAnimatedWorker1(b *testing.B) {
	benchmarkResize(b, testdata.Animated, 1)
}

// Worker 2

func BenchmarkResizeSmallWorker2(b *testing.B) {
	benchmarkResize(b, testdata.Small, 2)
}

func BenchmarkResizeMediumWorker2(b *testing.B) {
	benchmarkResize(b, testdata.Medium, 2)
}

func BenchmarkResizeLargeWorker2(b *testing.B) {
	benchmarkResize(b, testdata.Large, 2)
}

func BenchmarkResizeHugeWorker2(b *testing.B) {
	benchmarkResize(b, testdata.Huge, 2)
}

func BenchmarkResizeAnimatedWorker2(b *testing.B) {
	benchmarkResize(b, testdata.Animated, 2)
}

// Worker 4

func BenchmarkResizeSmallWorker4(b *testing.B) {
	benchmarkResize(b, testdata.Small, 4)
}

func BenchmarkResizeMediumWorker4(b *testing.B) {
	benchmarkResize(b, testdata.Medium, 4)
}

func BenchmarkResizeLargeWorker4(b *testing.B) {
	benchmarkResize(b, testdata.Large, 4)
}

func BenchmarkResizeHugeWorker4(b *testing.B) {
	benchmarkResize(b, testdata.Huge, 4)
}

func BenchmarkResizeAnimatedWorker4(b *testing.B) {
	benchmarkResize(b, testdata.Animated, 4)
}

// Worker 8

func BenchmarkResizeSmallWorker8(b *testing.B) {
	benchmarkResize(b, testdata.Small, 8)
}

func BenchmarkResizeMediumWorker8(b *testing.B) {
	benchmarkResize(b, testdata.Medium, 8)
}

func BenchmarkResizeLargeWorker8(b *testing.B) {
	benchmarkResize(b, testdata.Large, 8)
}

func BenchmarkResizeHugeWorker8(b *testing.B) {
	benchmarkResize(b, testdata.Huge, 8)
}

func BenchmarkResizeAnimatedWorker8(b *testing.B) {
	benchmarkResize(b, testdata.Animated, 8)
}

func benchmarkResize(b *testing.B, image *imageserver.Image, workerCount int) {
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
