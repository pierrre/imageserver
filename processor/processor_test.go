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
	pf := Func(func(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
		return image, nil
	})
	pf.Process(testdata.Small, make(imageserver.Parameters))
}

func TestListInterface(t *testing.T) {
	var _ Processor = List{}
}

func TestLimitInterface(t *testing.T) {
	var _ Processor = &Limit{}
}
