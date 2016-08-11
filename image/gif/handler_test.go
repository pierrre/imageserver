package gif

import (
	"bytes"
	"fmt"
	"image/gif"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Handler = &Handler{}

func TestHandler(t *testing.T) {
	hdr := &Handler{
		Processor: ProcessorFunc(func(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
			return g, nil
		}),
	}
	im, err := hdr.Handle(testdata.Animated, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if im.Format != "gif" {
		t.Fatalf("unexpected Format: got %s, want %s", im.Format, "gif")
	}
	_, err = gif.DecodeConfig(bytes.NewReader(im.Data))
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerNoChange(t *testing.T) {
	hdr := &Handler{
		Processor: testProcessorChange(false),
	}
	im, err := hdr.Handle(testdata.Animated, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if im != testdata.Animated {
		t.Fatal("not equal")
	}
}

func TestHandlerErrorFormat(t *testing.T) {
	hdr := &Handler{}
	_, err := hdr.Handle(testdata.Medium, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: got %T, want %T", err, &imageserver.ImageError{})
	}
}

func TestHandlerErrorDecode(t *testing.T) {
	hdr := &Handler{
		Processor: ProcessorFunc(func(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
			return g, nil
		}),
	}
	im := &imageserver.Image{
		Format: "gif",
	}
	_, err := hdr.Handle(im, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: got %T, want %T", err, &imageserver.ImageError{})
	}
}

func TestHandlerErrorProcessor(t *testing.T) {
	errPrc := fmt.Errorf("error")
	hdr := &Handler{
		Processor: ProcessorFunc(func(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
			return nil, errPrc
		}),
	}
	_, err := hdr.Handle(testdata.Animated, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if err != errPrc {
		t.Fatalf("unexpected error: got %#v, want %#v", err, errPrc)
	}
}

func TestHandlerErrorEncode(t *testing.T) {
	hdr := &Handler{
		Processor: ProcessorFunc(func(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
			return &gif.GIF{}, nil
		}),
	}
	_, err := hdr.Handle(testdata.Animated, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: got %T, want %T", err, &imageserver.ImageError{})
	}
}

var _ imageserver.Handler = &FallbackHandler{}

func TestFallbackHandler(t *testing.T) {
	var hc bool
	hdr := &FallbackHandler{
		Handler: &Handler{
			Processor: ProcessorFunc(func(g *gif.GIF, params imageserver.Params) (*gif.GIF, error) {
				hc = true
				return g, nil
			}),
		},
		Fallback: imageserver.HandlerFunc(func(im *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
			hc = false
			if !params.Has("format") {
				return im, nil
			}
			format, err := params.GetString("format")
			if err != nil {
				return nil, err
			}
			im = &imageserver.Image{
				Format: format,
			}
			return im, nil
		}),
	}
	for _, tc := range []struct {
		name               string
		image              *imageserver.Image
		params             imageserver.Params
		expectedHandler    bool
		expectedFormat     string
		expectedParamError string
	}{
		{
			name:            "JPEGDefaultFallbackHandler",
			image:           testdata.Medium,
			params:          imageserver.Params{},
			expectedHandler: false,
			expectedFormat:  "jpeg",
		},
		{
			name:            "GifDefaultHandler",
			image:           testdata.Animated,
			params:          imageserver.Params{},
			expectedHandler: true,
			expectedFormat:  "gif",
		},
		{
			name:               "InvalidFormat",
			image:              testdata.Animated,
			params:             imageserver.Params{"format": 666},
			expectedParamError: "format",
		},
		{
			name:            "GifParamFallbackHandler",
			image:           testdata.Animated,
			params:          imageserver.Params{"format": "jpeg"},
			expectedHandler: false,
			expectedFormat:  "jpeg",
		},
		{
			name:            "GifParamHandler",
			image:           testdata.Animated,
			params:          imageserver.Params{"format": "gif"},
			expectedHandler: true,
			expectedFormat:  "gif",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			im, err := hdr.Handle(tc.image, tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && err.Param == tc.expectedParamError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedParamError != "" {
				t.Fatal("no error")
			}
			if im.Format != tc.expectedFormat {
				t.Fatalf("unexpected format: got %s, want %s", im.Format, tc.expectedFormat)
			}
			if hc != tc.expectedHandler {
				t.Fatalf("wrong Handler called: got %t, want %t", hc, tc.expectedHandler)
			}
		})
	}
}
