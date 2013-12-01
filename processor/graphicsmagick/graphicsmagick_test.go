package graphicsmagick

import (
	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"testing"
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
	processor := &GraphicsMagickProcessor{
		Executable: "gm",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.Process(image, parameters)
		if err != nil {
			b.Fatal(err)
		}
	}
}
