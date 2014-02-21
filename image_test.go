package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestImage(t *testing.T) {
	for _, image := range testdata.Images {
		data, _ := image.MarshalBinary()

		newImage, err := NewImageUnmarshalBinary(data)
		if err != nil {
			t.Fatal(err)
		}

		if !ImageEqual(newImage, image) {
			t.Fatal("image not equals")
		}
	}
}

func TestImageUnmarshalBinaryError(t *testing.T) {
	for _, image := range testdata.Images {
		data, _ := image.MarshalBinary()

		index := -1 // Always truncate 1 byte
		for _, offset := range []int{
			4,
			len(image.Format),
			4,
			len(image.Data),
		} {
			index += offset
			errorData := data[0:index]
			_, err := NewImageUnmarshalBinary(errorData)
			if err == nil {
				t.Fatal("no error")
			}
		}
	}
}

func TestImageEqual(t *testing.T) {
	for _, test := range []struct {
		image1      *Image
		image2      *Image
		result      bool
		description string
	}{
		{nil, nil, true, "both nil"},
		{nil, testdata.Medium, false, "nil / not nil"},
		{testdata.Medium, nil, false, "not nil / nil"},
		{testdata.Medium, testdata.Medium, true, "same"},
		{testdata.Medium, imageCopy(testdata.Medium), true, "copy"},
		{testdata.Medium, testdata.Animated, false, "different format"},
		{testdata.Medium, testdata.Small, false, "different data"},
	} {
		if ImageEqual(test.image1, test.image2) != test.result {
			t.Fatalf("invalid result for test: %s", test.description)
		}
	}
}

func imageCopy(image *Image) *Image {
	value := *image
	return &value
}
