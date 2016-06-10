package http

import (
	"context"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkServerGetSmall(b *testing.B) {
	benchmark(b, testdata.SmallFileName)
}

func BenchmarkServerGetMedium(b *testing.B) {
	benchmark(b, testdata.MediumFileName)
}

func BenchmarkServerGetLarge(b *testing.B) {
	benchmark(b, testdata.LargeFileName)
}

func BenchmarkServerGetHuge(b *testing.B) {
	benchmark(b, testdata.HugeFileName)
}

func benchmark(b *testing.B, filename string) {
	httpSrv := createTestHTTPServer()
	defer httpSrv.Close()
	srv := &Server{}
	ctx := context.Background()
	params := imageserver.Params{
		imageserver_source.Param: createTestSource(httpSrv, filename),
	}
	var bs int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im, err := srv.Get(ctx, params)
		if err != nil {
			b.Fatal(err)
		}
		bs = len(im.Data)
	}
	b.SetBytes(int64(bs))
}
