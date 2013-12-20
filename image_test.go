package imageserver_test

import (
	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
	"reflect"
	"testing"
)

func TestImage(t *testing.T) {
	for _, image := range testdata.Images {
		data, err := image.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		newImage, err := NewImageUnmarshalBinary(data)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(newImage, image) {
			t.Fatal("image not equals")
		}
	}
}

func TestImageUnmarshalBinaryError(t *testing.T) {
	for _, image := range testdata.Images {
		data, err := image.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		index := -1 // Always truncate 1 byte
		for _, offset := range []int{
			4,
			len(image.Format),
			4,
			len(image.Data),
		} {
			index += offset
			errorData := data[0:index]
			_, err = NewImageUnmarshalBinary(errorData)
			if err == nil {
				t.Fatal("no error")
			}
		}
	}
}

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
		_, err := image.MarshalBinary()
		if err != nil {
			b.Fatal(err)
		}
	}
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
	data, err := image.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewImageUnmarshalBinary(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestImageOptimized(t *testing.T) {
	for _, image := range testdata.Images {
		data, err := image.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		newImage, err := NewImageUnmarshalBinaryOptimized(data)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(newImage, image) {
			t.Fatal("image not equals")
		}
	}
}

func TestImageUnmarshalBinaryOptimizedError(t *testing.T) {
	for _, image := range testdata.Images {
		data, err := image.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		index := -1 // Always truncate 1 byte
		for _, offset := range []int{
			4,
			len(image.Format),
			4,
			len(image.Data),
		} {
			index += offset
			errorData := data[0:index]
			_, err = NewImageUnmarshalBinaryOptimized(errorData)
			if err == nil {
				t.Fatal("no error")
			}
		}
	}
}

func BenchmarkImageUnmarshalBinaryOptimizedSmall(b *testing.B) {
	benchmarkImageUnmarshalBinaryOptimized(b, testdata.Small)
}

func BenchmarkImageUnmarshalBinaryOptimizedMedium(b *testing.B) {
	benchmarkImageUnmarshalBinaryOptimized(b, testdata.Medium)
}

func BenchmarkImageUnmarshalBinaryOptimizedLarge(b *testing.B) {
	benchmarkImageUnmarshalBinaryOptimized(b, testdata.Large)
}

func BenchmarkImageUnmarshalBinaryOptimizedHuge(b *testing.B) {
	benchmarkImageUnmarshalBinaryOptimized(b, testdata.Huge)
}

func BenchmarkImageUnmarshalBinaryOptimizedAnimated(b *testing.B) {
	benchmarkImageUnmarshalBinaryOptimized(b, testdata.Animated)
}

func benchmarkImageUnmarshalBinaryOptimized(b *testing.B, image *Image) {
	data, err := image.MarshalBinary()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewImageUnmarshalBinaryOptimized(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
