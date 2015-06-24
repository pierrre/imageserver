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

func benchmarkResize(b *testing.B, im *imageserver.Image) {
	server := &Server{
		Server:     &imageserver.StaticServer{Image: im},
		Executable: "gm",
	}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width": 100,
		},
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := server.Get(params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.SetBytes(int64(len(im.Data)))
}
