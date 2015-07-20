package tiff

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image_test "github.com/pierrre/imageserver/image/_test"
	_ "github.com/pierrre/imageserver/image/jpeg"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkSmall(b *testing.B) {
	benchmark(b, testdata.Small)
}

func BenchmarkMedium(b *testing.B) {
	benchmark(b, testdata.Medium)
}

func BenchmarkLarge(b *testing.B) {
	benchmark(b, testdata.Large)
}

func BenchmarkHuge(b *testing.B) {
	benchmark(b, testdata.Huge)
}

func benchmark(b *testing.B, im *imageserver.Image) {
	enc := &Encoder{}
	params := imageserver.Params{}
	b.ResetTimer()
	imageserver_image_test.BenchmarkEncoder(b, enc, im, params)
}
