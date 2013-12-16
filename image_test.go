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
	_, err := NewImageUnmarshalBinary(nil)
	if err == nil {
		t.Fatal("no error")
	}
}
