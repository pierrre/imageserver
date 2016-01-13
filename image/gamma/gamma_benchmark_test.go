package gamma

import (
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
	params := imageserver.Params{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prc.Process(nim, params)
	}
}

func BenchmarkProcessorHighQuality(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewProcessor(2.2, true)
	params := imageserver.Params{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prc.Process(nim, params)
	}
}

func BenchmarkCorrectionProcessor(b *testing.B) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	prc := NewCorrectionProcessor(
		imageserver_image.ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nim, nil
		}),
		true,
	)
	params := imageserver.Params{}
	for i := 0; i < b.N; i++ {
		prc.Process(nim, params)
	}
}
