package gif

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"reflect"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

var _ Processor = &SimpleProcessor{}

func TestSimpleProcessor(t *testing.T) {
	g1 := newTestImage()
	prc := &SimpleProcessor{
		Processor: imageserver_image.ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return image.NewRGBA(image.Rectangle{Min: nim.Bounds().Min.Div(2), Max: nim.Bounds().Max.Div(2)}), nil
		}),
	}
	g2, err := prc.Process(g1, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if len(g1.Image) != len(g2.Image) {
		t.Fatalf("Image length not equal: %d & %d", len(g1.Image), len(g2.Image))
	}
	if !reflect.DeepEqual(g1.Delay, g2.Delay) {
		t.Fatalf("Delay not equal: %#v & %#v", g1.Delay, g2.Delay)
	}
	if g1.LoopCount != g2.LoopCount {
		t.Fatalf("LoopCount not equal: %d & %d", g1.LoopCount, g2.LoopCount)
	}
	if !reflect.DeepEqual(g1.Config.ColorModel, g2.Config.ColorModel) {
		t.Fatalf("Config.ColorModel not equal: %#v & %#v", g1.Config.ColorModel, g2.Config.ColorModel)
	}
	if g2.Config.Width != 50 {
		t.Fatalf("unexpected Config.Width value: got %d, want %d", g2.Config.Width, 50)
	}
	if g2.Config.Height != 50 {
		t.Fatalf("unexpected Config.Height value: got %d, want %d", g2.Config.Height, 50)
	}
	if g1.BackgroundIndex != g2.BackgroundIndex {
		t.Fatalf("BackgroundIndex not equal: %d & %d", g1.BackgroundIndex, g2.BackgroundIndex)
	}
}

func TestSimpleProcessorError(t *testing.T) {
	g := newTestImage()
	prc := &SimpleProcessor{
		Processor: imageserver_image.ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := prc.Process(g, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

var _ Processor = ProcessorFunc(nil)

func TestProcessorFunc(t *testing.T) {
	g := newTestImage()
	called := false
	prc := ProcessorFunc(func(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
		called = true
		return g, nil
	})
	_, err := prc.Process(g, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("not called")
	}
	if !prc.Change(imageserver.Params{}) {
		t.Fatal("Change() returned false")
	}
}

func newTestImage() *gif.GIF {
	pl := color.Palette{
		color.RGBA{0xff, 0, 0, 0xff},
		color.RGBA{0, 0xff, 0, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}
	g := &gif.GIF{
		Image: []*image.Paletted{
			image.NewPaletted(image.Rect(10, 10, 80, 80), pl),
			image.NewPaletted(image.Rect(0, 0, 100, 100), pl),
		},
		Delay:     []int{0, 1},
		LoopCount: 666,
		Disposal:  []byte{gif.DisposalNone, gif.DisposalNone},
		Config: image.Config{
			Width:  100,
			Height: 100,
		},
		BackgroundIndex: 1,
	}
	return g
}

type testProcessorChange bool

func (prc testProcessorChange) Process(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
	return g, nil
}

func (prc testProcessorChange) Change(params imageserver.Params) bool {
	return bool(prc)
}
