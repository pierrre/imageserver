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
	for i := 0; i < b.N; i++ {
		_, err := hdr.Handle(testdata.Animated, imageserver.Params{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
