package provider_test

import (
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/provider"
	"github.com/pierrre/imageserver/testdata"
)

func TestFuncInterface(t *testing.T) {
	var _ Provider = Func(nil)
}

func TestFunc(t *testing.T) {
	pf := Func(func(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
		return testdata.Small, nil
	})
	pf.Get("foo", make(imageserver.Parameters))
}
