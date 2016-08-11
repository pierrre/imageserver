package http

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkServerGet(b *testing.B) {
	httpSrv := createTestHTTPServer()
	defer httpSrv.Close()
	srv := &Server{}
	for _, tc := range []struct {
		name     string
		filename string
	}{
		{"Small", testdata.SmallFileName},
		{"Medium", testdata.MediumFileName},
		{"Large", testdata.LargeFileName},
		{"Huge", testdata.HugeFileName},
	} {
		b.Run(tc.name, func(b *testing.B) {
			params := imageserver.Params{
				imageserver_source.Param: createTestSource(httpSrv, tc.filename),
			}
			var bs int
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				im, err := srv.Get(params)
				if err != nil {
					b.Fatal(err)
				}
				bs = len(im.Data)
			}
			b.SetBytes(int64(bs))
		})
	}
}
