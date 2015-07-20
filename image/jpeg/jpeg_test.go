package jpeg

import (
	"io/ioutil"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	imageserver_image_test "github.com/pierrre/imageserver/image/_test"
)

var _ imageserver_image.Encoder = &Encoder{}

func TestEncoder(t *testing.T) {
	testEncoder(t, &Encoder{})
}

func TestEncoderDefaultQuality(t *testing.T) {
	enc := &Encoder{
		DefaultQuality: 90,
	}
	testEncoder(t, enc)
}

func TestEncoderQuality(t *testing.T) {
	params := imageserver.Params{
		"quality": 90,
	}
	testEncoderParams(t, &Encoder{}, params)
}

func TestEncoderErrorQuality(t *testing.T) {
	im := imageserver_image_test.NewImage()
	enc := &Encoder{}
	for _, quality := range []interface{}{"foo", -1, 101} {
		err := enc.Encode(ioutil.Discard, im, imageserver.Params{"quality": quality})
		if err == nil {
			t.Fatal("no error")
		}
		errParam, ok := err.(*imageserver.ParamError)
		if !ok {
			t.Fatalf("unexpected error type: %T", err)
		}
		if errParam.Param != "quality" {
			t.Fatalf("unexpected param: %s", errParam.Param)
		}
	}
}

func testEncoder(t *testing.T, enc *Encoder) {
	imageserver_image_test.TestEncoder(t, enc, "jpeg")
}

func testEncoderParams(t *testing.T, enc *Encoder, params imageserver.Params) {
	imageserver_image_test.TestEncoderParams(t, enc, params, "jpeg")
}

func TestEncoderChange(t *testing.T) {
	c := (&Encoder{}).Change(imageserver.Params{})
	if c != false {
		t.Fatal("not false")
	}
}

func TestEncoderChangeQuality(t *testing.T) {
	c := (&Encoder{}).Change(imageserver.Params{"quality": 75})
	if c != true {
		t.Fatal("not true")
	}
}
