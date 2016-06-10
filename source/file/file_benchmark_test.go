package file

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkServerGetSmall(b *testing.B) {
	benchmarkServerGet(b, testdata.SmallFileName)
}

func BenchmarkServerGetMedium(b *testing.B) {
	benchmarkServerGet(b, testdata.MediumFileName)
}

func BenchmarkServerGetLarge(b *testing.B) {
	benchmarkServerGet(b, testdata.LargeFileName)
}

func BenchmarkServerGetHuge(b *testing.B) {
	benchmarkServerGet(b, testdata.HugeFileName)
}

func benchmarkServerGet(b *testing.B, filename string) {
	srv := &Server{
		Root: testdata.Dir,
	}
	params := imageserver.Params{
		imageserver_source.Param: filename,
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
}
