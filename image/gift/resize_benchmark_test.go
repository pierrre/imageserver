package gift

import (
	"context"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkResizeProcessorSize(b *testing.B) {
	for _, tc := range []struct {
		name string
		im   *imageserver.Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
		{"Huge", testdata.Huge},
	} {
		benchmarkResizeProcessor(b, tc.name, tc.im, imageserver.Params{})
	}
}

func BenchmarkResizeProcessorResampling(b *testing.B) {
	for _, r := range []string{
		"nearest_neighbor",
		"box",
		"linear",
		"cubic",
		"lanczos",
	} {
		benchmarkResizeProcessor(b, r, testdata.Medium, imageserver.Params{
			"resampling": r,
		})
	}
}

func benchmarkResizeProcessor(b *testing.B, name string, im *imageserver.Image, params imageserver.Params) {
	nim, err := imageserver_image.Decode(im)
	if err != nil {
		b.Fatal(err)
	}
	params.Set("width", 100)
	params = imageserver.Params{
		resizeParam: params,
	}
	prc := &ResizeProcessor{}
	ctx := context.Background()
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := prc.Process(ctx, nim, params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
