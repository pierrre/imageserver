package image_test

import (
	"image"
	"io"
	"io/ioutil"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/jpeg"
	"github.com/pierrre/imageserver/testdata"
)

var _ Encoder = EncoderFunc(nil)

func TestEncoderFunc(t *testing.T) {
	called := false
	f := EncoderFunc(func(w io.Writer, nim image.Image, params imageserver.Params) error {
		called = true
		return nil
	})
	nim := image.NewRGBA(image.Rect(0, 0, 1, 1))
	err := f.Encode(ioutil.Discard, nim, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("not called")
	}
	if f.Change(imageserver.Params{}) != true {
		t.Fatal("not true")
	}
}

func TestDecode(t *testing.T) {
	nim, err := Decode(testdata.Medium)
	if err != nil {
		t.Fatal(err)
	}
	if nim == nil {
		t.Fatal("image nil")
	}
}

func TestDecodeErrorInvalid(t *testing.T) {
	_, err := Decode(testdata.Invalid)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestDecodeErrorFormat(t *testing.T) {
	im := &imageserver.Image{Format: "error", Data: testdata.Medium.Data}
	_, err := Decode(im)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}
