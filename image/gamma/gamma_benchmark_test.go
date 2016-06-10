package gamma

import (
	"context"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkProcessor(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewProcessor(2.2, false)
	ctx := context.Background()
	params := imageserver.Params{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(ctx, nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessorHighQuality(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewProcessor(2.2, true)
	ctx := context.Background()
	params := imageserver.Params{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(ctx, nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCorrectionProcessor(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewCorrectionProcessor(
		imageserver_image.ProcessorFunc(func(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
			return nim, nil
		}),
		true,
	)
	ctx := context.Background()
	params := imageserver.Params{}
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(ctx, nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
