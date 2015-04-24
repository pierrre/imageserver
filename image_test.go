package imageserver_test

import (
	"encoding"
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ encoding.BinaryMarshaler = new(Image)
var _ encoding.BinaryUnmarshaler = new(Image)

func TestImageMarshal(t *testing.T) {
	for _, im := range testdata.Images {
		data, _ := im.MarshalBinary()
		imNew := new(Image)
		err := imNew.UnmarshalBinary(data)
		if err != nil {
			t.Fatal(err)
		}
		if !ImageEqual(imNew, im) {
			t.Fatal("image not equals")
		}
	}
}

func TestImageUnmarshalBinaryErrorEndOfData(t *testing.T) {
	for _, im := range testdata.Images {
		data, _ := im.MarshalBinary()
		index := -1 // Always truncate 1 byte
		for _, offset := range []int{
			4,
			len(im.Format),
			4,
			len(im.Data),
		} {
			index += offset
			errorData := data[0:index]
			imNew := new(Image)
			err := imNew.UnmarshalBinary(errorData)
			if err == nil {
				t.Fatal("no error")
			}
		}
	}
}

func TestImageEqual(t *testing.T) {
	for _, tc := range []struct {
		im1         *Image
		im2         *Image
		equal       bool
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
		if ImageEqual(tc.im1, tc.im2) != tc.equal {
			t.Fatalf("invalid result for test: %s", tc.description)
		}
	}
}

func TestImageError(t *testing.T) {
	err := &ImageError{Message: "test"}
	_ = err.Error()
}

func imageCopy(im *Image) *Image {
	value := *im
	return &value
}
