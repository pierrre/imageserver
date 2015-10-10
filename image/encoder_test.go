package image_test

import (
	"fmt"
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
	nim := NewTestImage()
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

func TestEncode(t *testing.T) {
	nim := NewTestImage()
	im, err := Encode(nim, "jpeg", imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if im == nil {
		t.Fatal("image nil")
	}
	if im.Format != "jpeg" {
		t.Fatalf("unexpected format: got %s, want %s", im.Format, "jpeg")
	}
}

func TestEncodeErrorFormat(t *testing.T) {
	nim := NewTestImage()
	_, err := Encode(nim, "foo", imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestEncodeErrorEncoderParams(t *testing.T) {
	nim := NewTestImage()
	_, err := Encode(nim, "jpeg", imageserver.Params{"quality": 9001})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

var _ imageserver.Server = &DecodeCheckServer{}

func TestDecodeCheckServer(t *testing.T) {
	srv := &DecodeCheckServer{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeCheckServerErrorServer(t *testing.T) {
	srv := &DecodeCheckServer{
		Server: &imageserver.StaticServer{
			Error: fmt.Errorf("error"),
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecodeCheckServerErrorPreDecode(t *testing.T) {
	srv := &DecodeCheckServer{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
		PreDecode: func(im *imageserver.Image, params imageserver.Params) error {
			return fmt.Errorf("error")
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestDecodeCheckServerErrorDecodeConfig(t *testing.T) {
	srv := &DecodeCheckServer{
		Server: &imageserver.StaticServer{
			Image: testdata.Invalid,
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestDecodeCheckServerErrorFormat(t *testing.T) {
	srv := &DecodeCheckServer{
		Server: &imageserver.StaticServer{
			Image: &imageserver.Image{
				Format: "error",
				Data:   testdata.Medium.Data,
			},
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestDecodeCheckServerErrorPostDecode(t *testing.T) {
	srv := &DecodeCheckServer{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
		PostDecode: func(cfg image.Config, format string, params imageserver.Params) error {
			return fmt.Errorf("error")
		},
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}
