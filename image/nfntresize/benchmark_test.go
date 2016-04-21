package nfntresize

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
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
	benchmark(b, im, imageserver.Params{})
}

func BenchmarkInterpolationNearestNeighbor(b *testing.B) {
	benchmarkInterpolation(b, "nearest_neighbor")
}

func BenchmarkInterpolationBilinear(b *testing.B) {
	benchmarkInterpolation(b, "bilinear")
}

func BenchmarkInterpolationBicubic(b *testing.B) {
	benchmarkInterpolation(b, "bicubic")
}

func BenchmarkInterpolationMitchellNetravali(b *testing.B) {
	benchmarkInterpolation(b, "mitchell_netravali")
}

func BenchmarkInterpolationLanczos2(b *testing.B) {
	benchmarkInterpolation(b, "lanczos2")
}

func BenchmarkInterpolationLanczos3(b *testing.B) {
	benchmarkInterpolation(b, "lanczos3")
}

func benchmarkInterpolation(b *testing.B, interp string) {
	params := imageserver.Params{
		"interpolation": interp,
	}
	benchmark(b, testdata.Medium, params)
}

func benchmark(b *testing.B, im *imageserver.Image, params imageserver.Params) {
	nim, err := imageserver_image.Decode(im)
	if err != nil {
		b.Fatal(err)
	}
	params.Set("width", 100)
	params = imageserver.Params{
		param: params,
	}
	proc := &Processor{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proc.Process(nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
