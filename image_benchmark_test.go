package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkImageMarshalBinarySmall(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Small)
}

func BenchmarkImageMarshalBinaryMedium(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Medium)
}

func BenchmarkImageMarshalBinaryLarge(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Large)
}

func BenchmarkImageMarshalBinaryHuge(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Huge)
}

func BenchmarkImageMarshalBinaryAnimated(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Animated)
}

func benchmarkImageMarshalBinary(b *testing.B, image *Image) {
	for i := 0; i < b.N; i++ {
		image.MarshalBinary()
	}

	b.SetBytes(int64(len(image.Data)))
}

func BenchmarkImageUnmarshalBinarySmall(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Small)
}

func BenchmarkImageUnmarshalBinaryMedium(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Medium)
}

func BenchmarkImageUnmarshalBinaryLarge(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Large)
}

func BenchmarkImageUnmarshalBinaryHuge(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Huge)
}

func BenchmarkImageUnmarshalBinaryAnimated(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Animated)
}

func benchmarkImageUnmarshalBinary(b *testing.B, image *Image) {
	data, _ := image.MarshalBinary()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewImageUnmarshalBinary(data)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.SetBytes(int64(len(data)))
}
