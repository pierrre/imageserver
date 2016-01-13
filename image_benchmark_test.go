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

func benchmarkImageMarshalBinary(b *testing.B, im *Image) {
	for i := 0; i < b.N; i++ {
		_, err := im.MarshalBinary()
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(im.Data)))
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

func benchmarkImageUnmarshalBinary(b *testing.B, im *Image) {
	data, err := im.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}
	imNew := new(Image)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := imNew.UnmarshalBinary(data)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(im.Data)))
}

func BenchmarkImageUnmarshalBinaryNoCopySmall(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Small)
}

func BenchmarkImageUnmarshalBinaryNoCopyMedium(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Medium)
}

func BenchmarkImageUnmarshalBinaryNoCopyLarge(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Large)
}

func BenchmarkImageUnmarshalBinaryNoCopyHuge(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Huge)
}

func benchmarkImageUnmarshalBinaryNoCopy(b *testing.B, im *Image) {
	data, err := im.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}
	imNew := new(Image)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := imNew.UnmarshalBinaryNoCopy(data)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.SetBytes(int64(len(im.Data)))
}
