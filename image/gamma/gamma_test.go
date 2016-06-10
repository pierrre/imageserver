package gamma

import (
	"context"
	"fmt"
	"image"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/jpeg"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver_image.Processor = &Processor{}

func TestProcessorProcess(t *testing.T) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		t.Fatal(err)
	}
	for _, highQuality := range []bool{false, true} {
		prc := NewProcessor(2.2, highQuality)
		_, err = prc.Process(context.Background(), nim, imageserver.Params{})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestProcessorChange(t *testing.T) {
	prc := NewProcessor(2.2, true)
	c := prc.Change(imageserver.Params{})
	if !c {
		t.Fatal("not true")
	}
}

var _ imageserver_image.Processor = &CorrectionProcessor{}

func TestCorrectionProcessor(t *testing.T) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		t.Fatal(err)
	}
	simplePrc := imageserver_image.ProcessorFunc(func(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
		return nim, nil
	})
	errPrc := imageserver_image.ProcessorFunc(func(ctx context.Context, nim image.Image, params imageserver.Params) (image.Image, error) {
		return nil, fmt.Errorf("error")
	})
	for _, tc := range []struct {
		name          string
		processor     imageserver_image.Processor
		enabled       bool
		params        imageserver.Params
		errorExpected bool
	}{
		{
			name:      "DefaultEnabled",
			processor: simplePrc,
			enabled:   true,
		},
		{
			name:      "DefaultDisabled",
			processor: simplePrc,
			enabled:   false,
		},
		{
			name:      "ParamDisabled",
			processor: simplePrc,
			enabled:   true,
			params: imageserver.Params{
				"gamma_correction": false,
			},
		},
		{
			name:      "ParamEnabled",
			processor: simplePrc,
			enabled:   false,
			params: imageserver.Params{
				"gamma_correction": true,
			},
		},
		{
			name:          "Error",
			processor:     errPrc,
			enabled:       true,
			errorExpected: true,
		},
		{
			name:      "Invalid",
			processor: simplePrc,
			enabled:   true,
			params: imageserver.Params{
				"gamma_correction": "invalid",
			},
			errorExpected: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prc := NewCorrectionProcessor(tc.processor, tc.enabled)
			params := tc.params
			if params == nil {
				params = imageserver.Params{}
			}
			_, err = prc.Process(context.Background(), nim, params)
			if tc.errorExpected && err == nil {
				t.Fatal("no error")
			} else if !tc.errorExpected && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIsHighQuality(t *testing.T) {
	r := image.Rect(0, 0, 1, 1)
	for _, tc := range []struct {
		name     string
		p        image.Image
		expected bool
	}{
		{
			name:     "RGBA64",
			p:        image.NewRGBA64(r),
			expected: true,
		},
		{
			name:     "NRGBA64",
			p:        image.NewNRGBA64(r),
			expected: true,
		},
		{
			name:     "RGBA",
			p:        image.NewRGBA(r),
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := isHighQuality(tc.p)
			if res != tc.expected {
				t.Fatalf("unexpected result for %T: got %t, want %t", tc.p, res, tc.expected)
			}
		})
	}
}
