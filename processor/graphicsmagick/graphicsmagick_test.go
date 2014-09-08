package graphicsmagick

import (
	"testing"
	"time"

	"github.com/pierrre/imageserver"
	imageserver_processor "github.com/pierrre/imageserver/processor"
	"github.com/pierrre/imageserver/testdata"
)

func TestInterface(t *testing.T) {
	var _ imageserver_processor.Processor = &Processor{}
}

func TestProcess(t *testing.T) {
	image := testdata.Medium

	parameters := imageserver.Parameters{
		"graphicsmagick": imageserver.Parameters{
			"width":  100,
			"height": 100,
		},
	}

	processor := &Processor{
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

	processor := &Processor{
		Executable: "gm",
		Timeout:    1 * time.Nanosecond,
	}

	_, err := processor.Process(image, parameters)
	if err == nil {
		t.Fatal("no error")
	}
}
