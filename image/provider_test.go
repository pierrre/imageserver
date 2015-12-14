package image

import (
	"fmt"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
)

var _ Provider = ProviderFunc(nil)

func TestProviderFunc(t *testing.T) {
	called := false
	prv := ProviderFunc(func(params imageserver.Params) (image.Image, error) {
		called = true
		return nil, nil
	})
	prv.Get(imageserver.Params{})
	if !called {
		t.Fatal("not called")
	}
}

var _ Provider = &StaticProvider{}

func TestStaticProvider(t *testing.T) {
	srv := &StaticProvider{
		Image: image.NewRGBA(image.Rect(0, 0, 1, 1)),
		Error: nil,
	}
	_, err := srv.Get(imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

var _ Provider = &ProcessorProvider{}

func TestProcessorProvider(t *testing.T) {
	providerCalled := false
	processorCalled := false
	prv := &ProcessorProvider{
		Provider: ProviderFunc(func(params imageserver.Params) (image.Image, error) {
			providerCalled = true
			return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
		}),
		Processor: ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			processorCalled = true
			return nim, nil
		}),
	}
	_, err := prv.Get(imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !providerCalled {
		t.Fatal("provider not called")
	}
	if !processorCalled {
		t.Fatal("processor not called")
	}
}

func TestProcessorProviderErrorProvider(t *testing.T) {
	prv := &ProcessorProvider{
		Provider: &StaticProvider{
			Error: fmt.Errorf("error"),
		},
	}
	_, err := prv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestProcessorProviderErrorProcessor(t *testing.T) {
	prv := &ProcessorProvider{
		Provider: &StaticProvider{
			Image: image.NewRGBA(image.Rect(0, 0, 1, 1)),
		},
		Processor: ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := prv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

var _ imageserver.Server = &ProviderServer{}

func TestProviderServer(t *testing.T) {
	srv := &ProviderServer{
		Provider: &StaticProvider{
			Image: image.NewRGBA(image.Rect(0, 0, 1, 1)),
		},
	}
	im, err := srv.Get(imageserver.Params{"format": "jpeg"})
	if err != nil {
		t.Fatal(err)
	}
	if im.Format != "jpeg" {
		t.Fatalf("unexpected format: got %s, want %s", im.Format, "jpeg")
	}
	_, err = Decode(im)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProviderServerErrorFormatNotSet(t *testing.T) {
	srv := &ProviderServer{}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestProviderServerErrorFormatUnknown(t *testing.T) {
	srv := &ProviderServer{}
	_, err := srv.Get(imageserver.Params{"format": "unknown"})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestProviderServerErrorProvider(t *testing.T) {
	srv := &ProviderServer{
		Provider: &StaticProvider{
			Error: fmt.Errorf("error"),
		},
	}
	_, err := srv.Get(imageserver.Params{"format": "jpeg"})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestProviderServerErrorEncode(t *testing.T) {
	srv := &ProviderServer{
		Provider: &StaticProvider{
			Image: image.NewRGBA(image.Rect(0, 0, 1, 1)),
		},
	}
	_, err := srv.Get(imageserver.Params{"format": "jpeg", "quality": 9001})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}
