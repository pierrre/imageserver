package graphicsmagick

import (
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkResizeSmall(b *testing.B) {
	benchmarkResize(b, testdata.Small)
}

func BenchmarkResizeMedium(b *testing.B) {
	benchmarkResize(b, testdata.Medium)
}

func BenchmarkResizeLarge(b *testing.B) {
	benchmarkResize(b, testdata.Large)
}

func BenchmarkResizeHuge(b *testing.B) {
	benchmarkResize(b, testdata.Huge)
}

func BenchmarkResizeAnimated(b *testing.B) {
	benchmarkResize(b, testdata.Animated)
}

func benchmarkResize(b *testing.B, image *imageserver.Image) {
	parameters := imageserver.Parameters{
		"graphicsmagick": imageserver.Parameters{
			"width":  100,
			"height": 100,
		},
	}

	processor := &Processor{
		Executable: "gm",
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := processor.Process(image, parameters)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.SetBytes(int64(len(image.Data)))
}
