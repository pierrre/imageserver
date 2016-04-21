package graphicsmagick

import (
	"os/exec"
	"testing"
	"time"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

const testExecutable = "gm"

var _ imageserver.Handler = &Handler{}

func TestHandle(t *testing.T) {
	testCheckAvailable(t)
	hdr := &Handler{
		Executable: testExecutable,
	}
	params := imageserver.Params{
		param: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := hdr.Handle(testdata.Medium, params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleErrorTimeout(t *testing.T) {
	testCheckAvailable(t)
	hdr := &Handler{
		Executable: testExecutable,
		Timeout:    1 * time.Nanosecond,
	}
	params := imageserver.Params{
		param: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := hdr.Handle(testdata.Medium, params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func testCheckAvailable(tb testing.TB) {
	_, err := exec.LookPath(testExecutable)
	if err != nil {
		tb.Skipf("GraphicsMagick is not available: %s", err)
	}
}
