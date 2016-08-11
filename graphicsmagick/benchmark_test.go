package graphicsmagick

import (
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkResize(b *testing.B) {
	testCheckAvailable(b)
	hdr := &Handler{
		Executable: testExecutable,
	}
	params := imageserver.Params{
		param: imageserver.Params{
			"width": 100,
		},
	}
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
			for i := 0; i < b.N; i++ {
				_, err := hdr.Handle(tc.im, params)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
