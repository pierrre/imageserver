package image

import (
	"context"
	"fmt"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
)

var _ Provider = ProviderFunc(nil)

func TestProviderFunc(t *testing.T) {
	called := false
	prv := ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
		called = true
		return nil, nil
	})
	_, _ = prv.Get(context.Background(), imageserver.Params{})
	if !called {
		t.Fatal("not called")
	}
}

var _ Provider = &ProcessorProvider{}

func TestProcessorProvider(t *testing.T) {
	providerCalled := false
	processorCalled := false
	prv := &ProcessorProvider{
		Provider: ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
			providerCalled = true
			return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
		}),
		Processor: ProcessorFunc(func(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
			processorCalled = true
			return nim, nil
		}),
	}
	_, err := prv.Get(context.Background(), imageserver.Params{})
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
		Provider: ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := prv.Get(context.Background(), imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestProcessorProviderErrorProcessor(t *testing.T) {
	prv := &ProcessorProvider{
		Provider: ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
			return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
		}),
		Processor: ProcessorFunc(func(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := prv.Get(context.Background(), imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}
