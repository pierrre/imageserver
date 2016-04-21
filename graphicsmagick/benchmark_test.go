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

func benchmarkResize(b *testing.B, im *imageserver.Image) {
	testCheckAvailable(b)
	hdr := &Handler{
		Executable: testExecutable,
	}
	params := imageserver.Params{
		param: imageserver.Params{
			"width": 100,
		},
	}
	for i := 0; i < b.N; i++ {
		_, err := hdr.Handle(im, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
