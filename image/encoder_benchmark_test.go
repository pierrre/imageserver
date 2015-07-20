package image_test

import (
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/image"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkDecodeSmall(b *testing.B) {
	benchmarkDecode(b, testdata.Small)
}

func BenchmarkDecodeMedium(b *testing.B) {
	benchmarkDecode(b, testdata.Medium)
}

func BenchmarkDecodeLarge(b *testing.B) {
	benchmarkDecode(b, testdata.Large)
}

func BenchmarkDecodeHuge(b *testing.B) {
	benchmarkDecode(b, testdata.Huge)
}

func benchmarkDecode(b *testing.B, im *imageserver.Image) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Decode(im)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecodeCheckServerSmall(b *testing.B) {
	benchmarkDecodeCheckServer(b, testdata.Small)
}

func BenchmarkDecodeCheckServerMedium(b *testing.B) {
	benchmarkDecodeCheckServer(b, testdata.Medium)
}

func BenchmarkDecodeCheckServerLarge(b *testing.B) {
	benchmarkDecodeCheckServer(b, testdata.Large)
}

func BenchmarkDecodeCheckServerHuge(b *testing.B) {
	benchmarkDecodeCheckServer(b, testdata.Huge)
}

func benchmarkDecodeCheckServer(b *testing.B, im *imageserver.Image) {
	srv := &DecodeCheckServer{
		Server: &imageserver.StaticServer{
			Image: im,
		},
	}
	params := imageserver.Params{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := srv.Get(params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
