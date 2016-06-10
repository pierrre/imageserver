package image

import (
	"context"
	"fmt"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
)

var _ imageserver.Server = &Server{}

func TestServer(t *testing.T) {
	srv := &Server{
		Provider: ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
			return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
		}),
	}
	im, err := srv.Get(context.Background(), imageserver.Params{"format": "jpeg"})
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

func TestServerDefaultFormat(t *testing.T) {
	srv := &Server{
		Provider: ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
			return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
		}),
		DefaultFormat: "jpeg",
	}
	im, err := srv.Get(context.Background(), imageserver.Params{})
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

func TestServerErrorFormatNotSet(t *testing.T) {
	srv := &Server{}
	_, err := srv.Get(context.Background(), imageserver.Params{})
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

func TestServerErrorFormatUnknown(t *testing.T) {
	srv := &Server{}
	_, err := srv.Get(context.Background(), imageserver.Params{"format": "unknown"})
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

func TestServerErrorProvider(t *testing.T) {
	srv := &Server{
		Provider: ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	_, err := srv.Get(context.Background(), imageserver.Params{"format": "jpeg"})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorEncode(t *testing.T) {
	srv := &Server{
		Provider: ProviderFunc(func(ctx context.Context, params imageserver.Params) (image.Image, error) {
			return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
		}),
	}
	_, err := srv.Get(context.Background(), imageserver.Params{"format": "jpeg", "quality": 9001})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}
