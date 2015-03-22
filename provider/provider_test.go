package provider_test

import (
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/provider"
	"github.com/pierrre/imageserver/testdata"
)

var _ Provider = Func(nil)

func TestFunc(t *testing.T) {
	Func(func(source interface{}, params imageserver.Params) (*imageserver.Image, error) {
		return testdata.Small, nil
	}).Get("foo", imageserver.Params{})
}
