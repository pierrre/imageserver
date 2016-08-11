package tiff

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image_test "github.com/pierrre/imageserver/image/_test"
	_ "github.com/pierrre/imageserver/image/jpeg"
	"github.com/pierrre/imageserver/testdata"
)

func Benchmark(b *testing.B) {
	for _, tc := range []struct {
		name string
		im   *imageserver.Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
		{"Huge", testdata.Huge},
	} {
		b.Run(tc.name, func(b *testing.B) {
			imageserver_image_test.BenchmarkEncoder(b, &Encoder{}, tc.im, imageserver.Params{})
		})
	}
}
