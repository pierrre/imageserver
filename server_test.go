package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestServerFuncInterface(t *testing.T) {
	var _ Server = ServerFunc(nil)
}

func TestServerFunc(t *testing.T) {
	ServerFunc(func(params Params) (*Image, error) {
		return testdata.Medium, nil
	}).Get(Params{})
}
