package png

import (
	"image/png"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image_test "github.com/pierrre/imageserver/image/_test"
	_ "github.com/pierrre/imageserver/image/jpeg"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkSizeSmall(b *testing.B) {
	benchmarkSize(b, testdata.Small)
}

func BenchmarkSizeMedium(b *testing.B) {
	benchmarkSize(b, testdata.Medium)
}

func BenchmarkSizeLarge(b *testing.B) {
	benchmarkSize(b, testdata.Large)
}

func BenchmarkSizeHuge(b *testing.B) {
	benchmarkSize(b, testdata.Huge)
}

func benchmarkSize(b *testing.B, im *imageserver.Image) {
	enc := &Encoder{}
	benchmark(b, enc, im)
}

func BenchmarkCompressionLevelDefaultCompression(b *testing.B) {
	benchmarkCompressionLevel(b, png.DefaultCompression)
}

func BenchmarkCompressionLevelNoCompression(b *testing.B) {
	benchmarkCompressionLevel(b, png.NoCompression)
}

func BenchmarkCompressionLevelBestSpeed(b *testing.B) {
	benchmarkCompressionLevel(b, png.BestSpeed)
}

func BenchmarkCompressionLevelBestCompression(b *testing.B) {
	benchmarkCompressionLevel(b, png.BestCompression)
}

func benchmarkCompressionLevel(b *testing.B, cl png.CompressionLevel) {
	enc := &Encoder{
		CompressionLevel: cl,
	}
	benchmark(b, enc, testdata.Medium)
}

func benchmark(b *testing.B, enc *Encoder, im *imageserver.Image) {
	params := imageserver.Params{}
	imageserver_image_test.BenchmarkEncoder(b, enc, im, params)
}
