package native

import (
	"testing"

	imageserver_processor "github.com/pierrre/imageserver/processor"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestNativeProcessorInterface(t *testing.T) {
	var _ imageserver_processor.Processor = &NativeProcessor{}
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
