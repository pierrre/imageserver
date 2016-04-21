package gift

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkResizeProcessorSizeSmall(b *testing.B) {
	benchmarkResizeProcessorSize(b, testdata.Small)
}

func BenchmarkResizeProcessorSizeMedium(b *testing.B) {
	benchmarkResizeProcessorSize(b, testdata.Medium)
}

func BenchmarkResizeProcessorSizeLarge(b *testing.B) {
	benchmarkResizeProcessorSize(b, testdata.Large)
}

func BenchmarkResizeProcessorSizeHuge(b *testing.B) {
	benchmarkResizeProcessorSize(b, testdata.Huge)
}

func benchmarkResizeProcessorSize(b *testing.B, im *imageserver.Image) {
	benchmarkResizeProcessor(b, im, imageserver.Params{})
}

func BenchmarkResizeProcessorResamplingNearestNeighbor(b *testing.B) {
	benchmarkResizeProcessorResampling(b, "nearest_neighbor")
}

func BenchmarkResizeProcessorResamplingBox(b *testing.B) {
	benchmarkResizeProcessorResampling(b, "box")
}

func BenchmarkResizeProcessorResamplingLinear(b *testing.B) {
	benchmarkResizeProcessorResampling(b, "linear")
}

func BenchmarkResizeProcessorResamplingCubic(b *testing.B) {
	benchmarkResizeProcessorResampling(b, "cubic")
}

func BenchmarkResizeProcessorResamplingLanczos(b *testing.B) {
	benchmarkResizeProcessorResampling(b, "lanczos")
}

func benchmarkResizeProcessorResampling(b *testing.B, rsp string) {
	params := imageserver.Params{
		"resampling": rsp,
	}
	benchmarkResizeProcessor(b, testdata.Medium, params)
}

func benchmarkResizeProcessor(b *testing.B, im *imageserver.Image, params imageserver.Params) {
	nim, err := imageserver_image.Decode(im)
	if err != nil {
		b.Fatal(err)
	}
	params.Set("width", 100)
	params = imageserver.Params{
		resizeParam: params,
	}
	prc := &ResizeProcessor{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
