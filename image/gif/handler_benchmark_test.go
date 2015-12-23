package gif

import (
	"image/gif"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkHandler(b *testing.B) {
	hdr := &Handler{
		Processor: ProcessorFunc(func(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
			return g, nil
		}),
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := hdr.Handle(testdata.Animated, imageserver.Params{})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
