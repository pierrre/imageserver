package gif

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image_test "github.com/pierrre/imageserver/image/_test"
	_ "github.com/pierrre/imageserver/image/jpeg"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkEncoderSmall(b *testing.B) {
	benchmarkEncoder(b, testdata.Small)
}

func BenchmarkEncoderMedium(b *testing.B) {
	benchmarkEncoder(b, testdata.Medium)
}

func BenchmarkEncoderLarge(b *testing.B) {
	benchmarkEncoder(b, testdata.Large)
}

func BenchmarkEncoderHuge(b *testing.B) {
	benchmarkEncoder(b, testdata.Huge)
}

func benchmarkEncoder(b *testing.B, im *imageserver.Image) {
	enc := &Encoder{}
	params := imageserver.Params{}
	b.ResetTimer()
	imageserver_image_test.BenchmarkEncoder(b, enc, im, params)
}
