package graphicsmagick

import (
	"context"
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
	_, err := hdr.Handle(context.Background(), testdata.Medium, params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleErrorTimeout(t *testing.T) {
	testCheckAvailable(t)
	hdr := &Handler{
		Executable: testExecutable,
	}
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	params := imageserver.Params{
		param: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := hdr.Handle(ctx, testdata.Medium, params)
	if err == nil {
		t.Fatal("no error")
	}
	if err != context.DeadlineExceeded {
		t.Fatalf("unexpected error: got %#v, want %#v", err, context.DeadlineExceeded)
	}
}

func testCheckAvailable(tb testing.TB) {
	_, err := exec.LookPath(testExecutable)
	if err != nil {
		tb.Skipf("GraphicsMagick is not available: %s", err)
	}
}
