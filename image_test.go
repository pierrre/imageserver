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

func BenchmarkNewImageUnmarshalBinarySmall(b *testing.B) {
	benchmarkNewImageUnmarshalBinary(b, testdata.Small)
}

func BenchmarkNewImageUnmarshalBinaryMedium(b *testing.B) {
	benchmarkNewImageUnmarshalBinary(b, testdata.Medium)
}

func BenchmarkNewImageUnmarshalBinaryLarge(b *testing.B) {
	benchmarkNewImageUnmarshalBinary(b, testdata.Large)
}

func BenchmarkNewImageUnmarshalBinaryHuge(b *testing.B) {
	benchmarkNewImageUnmarshalBinary(b, testdata.Huge)
}

func BenchmarkNewImageUnmarshalBinaryAnimated(b *testing.B) {
	benchmarkNewImageUnmarshalBinary(b, testdata.Animated)
}

func benchmarkNewImageUnmarshalBinary(b *testing.B, image *Image) {
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
