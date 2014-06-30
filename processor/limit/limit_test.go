package limit

import (
	"testing"

	"github.com/pierrre/imageserver"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestInterface(t *testing.T) {
	var _ imageserver.Processor = &LimitProcessor{}
}
