package gift

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

func BenchmarkResamplingNearestNeighbor(b *testing.B) {
	benchmarkResampling(b, "nearest_neighbor")
}

func BenchmarkResamplingBox(b *testing.B) {
	benchmarkResampling(b, "box")
}

func BenchmarkResamplingLinear(b *testing.B) {
	benchmarkResampling(b, "linear")
}

func BenchmarkResamplingCubic(b *testing.B) {
	benchmarkResampling(b, "cubic")
}

func BenchmarkResamplingLanczos(b *testing.B) {
	benchmarkResampling(b, "lanczos")
}

func benchmarkResampling(b *testing.B, rsp string) {
	params := imageserver.Params{
		"resampling": rsp,
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
		Param: params,
	}
	proc := &Processor{}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := proc.Process(nim, params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
