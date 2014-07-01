package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestImageServerFuncInterface(t *testing.T) {
	var _ ImageServer = ImageServerFunc(nil)
}

func TestImageServerFunc(t *testing.T) {
	isf := ImageServerFunc(func(parameters Parameters) (*Image, error) {
		return testdata.Medium, nil
	})
	isf.Get(Parameters{})
}
