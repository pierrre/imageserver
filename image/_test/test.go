// Package _test provides utilities for imageserver/image.Encoder testing.
package _test

import (
	"bytes"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

// TestEncoder is a helper to test imageserver/image.Encoder.
func TestEncoder(t *testing.T, enc imageserver_image.Encoder, expectedFormat string) {
	TestEncoderParams(t, enc, imageserver.Params{}, expectedFormat)
}

// TestEncoderParams is a helper to test imageserver/image.Encoder with Params.
func TestEncoderParams(t *testing.T, enc imageserver_image.Encoder, params imageserver.Params, expectedFormat string) {
	buf := new(bytes.Buffer)
	nim := NewImage()
	err := enc.Encode(buf, nim, params)
	if err != nil {
		t.Fatal(err)
	}
	_, format, err := image.Decode(buf)
	if err != nil {
		t.Fatal(err)
	}
	if format != expectedFormat {
		t.Fatalf("unexpected format: got %s, want %s", format, expectedFormat)
	}
}

// NewImage creates a new test Image.
func NewImage() image.Image {
	return image.NewRGBA(image.Rect(0, 0, 64, 64))
}
