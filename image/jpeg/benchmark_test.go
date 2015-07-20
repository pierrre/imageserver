package jpeg

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image_test "github.com/pierrre/imageserver/image/_test"
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
	params := imageserver.Params{}
	benchmark(b, im, params)
}

func BenchmarkQuality1(b *testing.B) {
	benchmarkQuality(b, 1)
}

func BenchmarkQuality25(b *testing.B) {
	benchmarkQuality(b, 25)
}

func BenchmarkQuality50(b *testing.B) {
	benchmarkQuality(b, 50)
}

func BenchmarkQuality75(b *testing.B) {
	benchmarkQuality(b, 75)
}

func BenchmarkQuality85(b *testing.B) {
	benchmarkQuality(b, 85)
}

func BenchmarkQuality90(b *testing.B) {
	benchmarkQuality(b, 90)
}

func BenchmarkQuality95(b *testing.B) {
	benchmarkQuality(b, 95)
}

func BenchmarkQuality100(b *testing.B) {
	benchmarkQuality(b, 100)
}

func benchmarkQuality(b *testing.B, quality int) {
	params := imageserver.Params{
		"quality": quality,
	}
	benchmark(b, testdata.Medium, params)
}

func benchmark(b *testing.B, im *imageserver.Image, params imageserver.Params) {
	enc := &Encoder{}
	imageserver_image_test.BenchmarkEncoder(b, enc, im, params)
}
