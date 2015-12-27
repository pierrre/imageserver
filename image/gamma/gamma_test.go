package gamma

import (
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
		_, err = prc.Process(nim, imageserver.Params{})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestProcessorChange(t *testing.T) {
	prc := NewProcessor(2.2, true)
	c := prc.Change(imageserver.Params{})
	if c != true {
		t.Fatal("not true")
	}
}

var _ imageserver_image.Processor = &CorrectionProcessor{}

func TestCorrectionProcessor(t *testing.T) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		t.Fatal(err)
	}
	simplePrc := imageserver_image.ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
		return nim, nil
	})
	errPrc := imageserver_image.ProcessorFunc(func(nim image.Image, params imageserver.Params) (image.Image, error) {
		return nil, fmt.Errorf("error")
	})
	type TC struct {
		processor     imageserver_image.Processor
		enabled       bool
		params        imageserver.Params
		errorExpected bool
	}
	for _, tc := range []TC{
		{
			processor: simplePrc,
			enabled:   true,
		},
		{
			processor: simplePrc,
			enabled:   false,
		},
		{
			processor: simplePrc,
			enabled:   true,
			params: imageserver.Params{
				"gamma_correction": false,
			},
		},
		{
			processor: simplePrc,
			enabled:   false,
			params: imageserver.Params{
				"gamma_correction": true,
			},
		},
		{
			processor:     errPrc,
			enabled:       true,
			errorExpected: true,
		},
		{
			processor: simplePrc,
			enabled:   true,
			params: imageserver.Params{
				"gamma_correction": "invalid",
			},
			errorExpected: true,
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			prc := NewCorrectionProcessor(tc.processor, tc.enabled)
			params := tc.params
			if params == nil {
				params = imageserver.Params{}
			}
			_, err = prc.Process(nim, params)
			if tc.errorExpected && err == nil {
				t.Fatal("no error")
			} else if !tc.errorExpected && err != nil {
				t.Fatal(err)
			}
		}()

	}
}

func TestIsHighQuality(t *testing.T) {
	r := image.Rect(0, 0, 1, 1)
	type TC struct {
		p        image.Image
		expected bool
	}
	for _, tc := range []TC{
		{
			p:        image.NewRGBA64(r),
			expected: true,
		},
		{
			p:        image.NewNRGBA64(r),
			expected: true,
		},
		{
			p:        image.NewRGBA(r),
			expected: false,
		},
	} {
		res := isHighQuality(tc.p)
		if res != tc.expected {
			t.Fatalf("unexpected result for %T: got %t, want %t", tc.p, res, tc.expected)
		}
	}
}
