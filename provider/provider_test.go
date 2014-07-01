package provider

import (
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestProviderFuncInterface(t *testing.T) {
	var _ Provider = ProviderFunc(nil)
}

func TestProviderFunc(t *testing.T) {
	pf := ProviderFunc(func(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
		return testdata.Small, nil
	})
	pf.Get("foo", make(imageserver.Parameters))
}
