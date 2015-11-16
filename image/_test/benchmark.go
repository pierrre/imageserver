package _test

import (
	"io/ioutil"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

// BenchmarkEncoder is a helper to benchmark Encoder.
func BenchmarkEncoder(b *testing.B, enc imageserver_image.Encoder, im *imageserver.Image, params imageserver.Params) {
	nim, err := imageserver_image.Decode(im)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := enc.Encode(ioutil.Discard, nim, params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
