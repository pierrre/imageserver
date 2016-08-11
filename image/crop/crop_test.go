package crop

import (
	"image"
	"image/color"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
)

var _ imageserver_image.Processor = &Processor{}

func TestProcess(t *testing.T) {
	prc := &Processor{}
	for _, tc := range []struct {
		name               string
		newImage           func() image.Image
		params             imageserver.Params
		expectedParamError string
		expectedImageError bool
		expectedBounds     image.Rectangle
	}{
		{
			name: "Empty",
			newImage: func() image.Image {
				return image.NewRGBA(image.Rect(0, 0, 100, 100))
			},
			params:         imageserver.Params{},
			expectedBounds: image.Rect(0, 0, 100, 100),
		},
		{
			name: "Crop",
			newImage: func() image.Image {
				return image.NewRGBA(image.Rect(0, 0, 100, 100))
			},
			params: imageserver.Params{param: imageserver.Params{
				"min_x": 20,
				"min_y": 20,
				"max_x": 50,
				"max_y": 50,
			}},
			expectedBounds: image.Rect(20, 20, 50, 50),
		},
		{
			name: "Invalid",
			newImage: func() image.Image {
				return image.NewRGBA(image.Rect(0, 0, 100, 100))
			},
			expectedParamError: "crop",
			params:             imageserver.Params{param: "invalid"},
		},
		{
			name: "InvalidMinX",
			newImage: func() image.Image {
				return image.NewRGBA(image.Rect(0, 0, 100, 100))
			},
			params: imageserver.Params{param: imageserver.Params{
				"min_x": "invalid",
				"min_y": 20,
				"max_x": 50,
				"max_y": 50,
			}},
			expectedParamError: "crop.min_x",
		},
		{
			name: "InvalidMinY",
			newImage: func() image.Image {
				return image.NewRGBA(image.Rect(0, 0, 100, 100))
			},
			params: imageserver.Params{param: imageserver.Params{
				"min_x": 20,
				"min_y": "invalid",
				"max_x": 50,
				"max_y": 50,
			}},
			expectedParamError: "crop.min_y",
		},
		{
			name: "InvalidMaxX",
			newImage: func() image.Image {
				return image.NewRGBA(image.Rect(0, 0, 100, 100))
			},
			params: imageserver.Params{param: imageserver.Params{
				"min_x": 20,
				"min_y": 20,
				"max_x": "invalid",
				"max_y": 50,
			}},
			expectedParamError: "crop.max_x",
		},
		{
			name: "InvalidMaxY",
			newImage: func() image.Image {
				return image.NewRGBA(image.Rect(0, 0, 100, 100))
			},
			params: imageserver.Params{param: imageserver.Params{
				"min_x": 20,
				"min_y": 20,
				"max_x": 50,
				"max_y": "invalid",
			}},
			expectedParamError: "crop.max_y",
		},
		{
			name: "ImageError",
			newImage: func() image.Image {
				return image.NewUniform(color.White)
			},
			params: imageserver.Params{param: imageserver.Params{
				"min_x": 20,
				"min_y": 20,
				"max_x": 50,
				"max_y": 50,
			}},
			expectedImageError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			im, err := prc.Process(tc.newImage(), tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && tc.expectedParamError == err.Param {
					return
				}
				if _, ok := err.(*imageserver.ImageError); ok && tc.expectedImageError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedParamError != "" {
				t.Fatalf("no param error, expected: %s", tc.expectedParamError)
			}
			if tc.expectedImageError {
				t.Fatal("no image error")
			}
			if im.Bounds() != tc.expectedBounds {
				t.Fatalf("unexpected bounds: got %#v, want %#v", im.Bounds(), tc.expectedBounds)
			}
		})
	}
}

func TestChange(t *testing.T) {
	prc := &Processor{}
	for _, tc := range []struct {
		name     string
		params   imageserver.Params
		expected bool
	}{
		{
			name:     "Empty",
			params:   imageserver.Params{},
			expected: false,
		},
		{
			name: "Crop",
			params: imageserver.Params{
				param: imageserver.Params{
					"min_x": 1,
					"min_y": 2,
					"max_x": 3,
					"max_y": 4,
				},
			},
			expected: true,
		},
		{
			name: "Invalid",
			params: imageserver.Params{
				param: "invalid",
			},
			expected: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			change := prc.Change(tc.params)
			if change != tc.expected {
				t.Fatalf("unexpected result for %s: got %t, want %t", tc.params, change, tc.expected)
			}
		})
	}
}
