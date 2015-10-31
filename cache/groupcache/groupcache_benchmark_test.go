package groupcache

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkServerSmall(b *testing.B) {
	benchmarkServer(b, testdata.Small)
}

func BenchmarkServerMedium(b *testing.B) {
	benchmarkServer(b, testdata.Medium)
}

func BenchmarkServerLarge(b *testing.B) {
	benchmarkServer(b, testdata.Large)
}

func BenchmarkServerHuge(b *testing.B) {
	benchmarkServer(b, testdata.Huge)
}

func benchmarkServer(b *testing.B, im *imageserver.Image) {
	srv := newTestServer(
		&imageserver.StaticServer{
			Image: im,
		},
		imageserver_cache.StringKeyGenerator("test"),
	)
	params := imageserver.Params{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := srv.Get(params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.SetBytes(int64(len(im.Data)))
}
