package native

import (
	"testing"

	imageserver_processor "github.com/pierrre/imageserver/processor"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestProcessorInterface(t *testing.T) {
	var _ imageserver_processor.Processor = &Processor{}
}

func TestDecoderFuncInterface(t *testing.T) {
	var _ Decoder = DecoderFunc(nil)
}

func TestProcessorNativeFuncInterface(t *testing.T) {
	var _ ProcessorNative = ProcessorNativeFunc(nil)
}

func TestEncoderFuncInterface(t *testing.T) {
	var _ Encoder = EncoderFunc(nil)
}
