package processor

import (
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestFuncInterface(t *testing.T) {
	var _ Processor = Func(nil)
}

func TestFunc(t *testing.T) {
	Func(func(image *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
		return image, nil
	}).Process(testdata.Small, imageserver.Params{})
}

func TestListInterface(t *testing.T) {
	var _ Processor = List{}
}
