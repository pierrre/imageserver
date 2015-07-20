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

func BenchmarkFullSmall(b *testing.B) {
	benchmarkFull(b, testdata.Small)
}

func BenchmarkFullMedium(b *testing.B) {
	benchmarkFull(b, testdata.Medium)
}

func BenchmarkFullLarge(b *testing.B) {
	benchmarkFull(b, testdata.Large)
}

func BenchmarkFullHuge(b *testing.B) {
	benchmarkFull(b, testdata.Huge)
}

func benchmarkFull(b *testing.B, im *imageserver.Image) {
	srv := &imageserver_image.Server{
		Server: &imageserver.StaticServer{
			Image: im,
		},
		Processor: &Processor{},
	}
	params := imageserver.Params{
		Param: imageserver.Params{
			"width": 100,
		},
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := srv.Get(params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
