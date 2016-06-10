package gif

import (
	"context"
	"image/gif"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkHandler(b *testing.B) {
	hdr := &Handler{
		Processor: ProcessorFunc(func(ctx context.Context, g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
			return g, nil
		}),
	}
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := hdr.Handle(ctx, testdata.Animated, imageserver.Params{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
