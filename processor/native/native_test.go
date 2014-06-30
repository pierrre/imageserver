package native

import (
	"testing"

	"github.com/pierrre/imageserver"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestNativeProcessorInterface(t *testing.T) {
	var _ imageserver.Processor = &NativeProcessor{}
}

func TestDecoderFuncInterface(t *testing.T) {
	var _ Decoder = DecoderFunc(nil)
}

func TestProcessorFuncInterface(t *testing.T) {
	var _ Processor = ProcessorFunc(nil)
}

func TestEncoderFuncInterface(t *testing.T) {
	var _ Encoder = EncoderFunc(nil)
}
