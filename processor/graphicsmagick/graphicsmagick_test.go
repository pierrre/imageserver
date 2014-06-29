package graphicsmagick

import (
	"testing"
	"time"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func TestInterfaceProcessor(t *testing.T) {
	var _ imageserver.Processor = &GraphicsMagickProcessor{}
}

func TestProcess(t *testing.T) {
	image := testdata.Medium

	parameters := imageserver.Parameters{
		"graphicsmagick": imageserver.Parameters{
			"width":  100,
			"height": 100,
		},
	}

	processor := &GraphicsMagickProcessor{
		Executable: "gm",
	}

	_, err := processor.Process(image, parameters)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProcessErrorTimeout(t *testing.T) {
	image := testdata.Medium

	parameters := imageserver.Parameters{
		"graphicsmagick": imageserver.Parameters{
			"width":  100,
			"height": 100,
		},
	}

	processor := &GraphicsMagickProcessor{
		Executable: "gm",
		Timeout:    1 * time.Nanosecond,
	}

	_, err := processor.Process(image, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}
