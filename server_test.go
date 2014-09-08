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
	sf := ServerFunc(func(parameters Parameters) (*Image, error) {
		return testdata.Medium, nil
	})
	sf.Get(Parameters{})
}
