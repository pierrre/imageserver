package nfntresize

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/jpeg"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

var _ imageserver_image.Processor = &Processor{}

func TestProcessor(t *testing.T) {
	nim, err := imageserver_image.Decode(imageserver_testdata.Medium)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name               string
		processor          *Processor
		params             imageserver.Params
		expectedWidth      int
		expectedHeight     int
		expectedParamError string
	}{
		// no size
		{
			name:           "Empty",
			params:         imageserver.Params{},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			name:           "ParamEmpty",
			params:         imageserver.Params{param: imageserver.Params{}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			name: "SizeZero",
			params: imageserver.Params{param: imageserver.Params{
				"width":  0,
				"height": 0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		// with size
		{
			name: "Width",
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			name: "Height",
			params: imageserver.Params{param: imageserver.Params{
				"height": 100,
			}},
			expectedHeight: 100,
		},
		{
			name: "WidthHeight",
			params: imageserver.Params{param: imageserver.Params{
				"width":  100,
				"height": 100,
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		// mode
		{
			name: "ModeResize",
			params: imageserver.Params{param: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "resize",
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		{
			name: "ModeThumbnail",
			params: imageserver.Params{param: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "thumbnail",
			}},
			expectedWidth:  100,
			expectedHeight: 79,
		},
		// interpolation
		{
			name: "InterpolationNearestNeighbor",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "nearest_neighbor",
			}},
			expectedWidth: 100,
		},
		{
			name: "InterpolationBilinear",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "bilinear",
			}},
			expectedWidth: 100,
		},
		{
			name: "InterpolationBicubic",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "bicubic",
			}},
			expectedWidth: 100,
		},
		{
			name: "InterpolationMitchelNetravali",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "mitchell_netravali",
			}},
			expectedWidth: 100,
		},
		{
			name: "InterpolationLanczos2",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos2",
			}},
			expectedWidth: 100,
		},
		{
			name: "InterpolationLanczos3",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos3",
			}},
			expectedWidth: 100,
		},
		// error
		{
			name:               "ParamInvalid",
			params:             imageserver.Params{param: "invalid"},
			expectedParamError: param,
		},
		{
			name: "WidthInvalidType",
			params: imageserver.Params{param: imageserver.Params{
				"width": "invalid",
			}},
			expectedParamError: param + ".width",
		},
		{
			name: "HeightInvalidType",
			params: imageserver.Params{param: imageserver.Params{
				"height": "invalid",
			}},
			expectedParamError: param + ".height",
		},
		{
			name: "WidthInvalidNegative",
			params: imageserver.Params{param: imageserver.Params{
				"width": -1,
			}},
			expectedParamError: param + ".width",
		},
		{
			name: "HeightInvalidNegative",
			params: imageserver.Params{param: imageserver.Params{
				"height": -1,
			}},
			expectedParamError: param + ".height",
		},
		{
			name:      "WidthInvalidTooLarge",
			processor: &Processor{MaxWidth: 500},
			params: imageserver.Params{param: imageserver.Params{
				"width": 1000,
			}},
			expectedParamError: param + ".width",
		},
		{
			name:      "HeightInvalidTooLarge",
			processor: &Processor{MaxHeight: 500},
			params: imageserver.Params{param: imageserver.Params{
				"height": 1000,
			}},
			expectedParamError: param + ".height",
		},
		{
			name: "InterpolationInvalidType",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": false,
			}},
			expectedParamError: param + ".interpolation",
		},
		{
			name: "InterpolationInvalidUnknown",
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "invalid",
			}},
			expectedParamError: param + ".interpolation",
		},
		{
			name: "ModeInvalidType",
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
				"mode":  false,
			}},
			expectedParamError: param + ".mode",
		},
		{
			name: "ModeInvalidUnknown",
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
				"mode":  "invalid",
			}},
			expectedParamError: param + ".mode",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prc := tc.processor
			if prc == nil {
				prc = &Processor{}
			}
			nim, err := prc.Process(nim, tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && err.Param == tc.expectedParamError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedParamError != "" {
				t.Fatal("no error")
			}
			if tc.expectedWidth != 0 && nim.Bounds().Dx() != tc.expectedWidth {
				t.Fatalf("unexpected width: got %d, want %d", nim.Bounds().Dx(), tc.expectedWidth)
			}
			if tc.expectedHeight != 0 && nim.Bounds().Dy() != tc.expectedHeight {
				t.Fatalf("unexpected height: got %d, want %d", nim.Bounds().Dy(), tc.expectedHeight)
			}
		})
	}
}

func TestProcessorChange(t *testing.T) {
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
			name:     "ParamEmpty",
			params:   imageserver.Params{param: imageserver.Params{}},
			expected: false,
		},
		{
			name:     "ParamInvalidType",
			params:   imageserver.Params{param: 666},
			expected: true,
		},
		{
			name: "Width",
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
			}},
			expected: true,
		},
		{
			name: "Height",
			params: imageserver.Params{param: imageserver.Params{
				"height": 100,
			}},
			expected: true,
		},
		{
			name: "UnknownParam",
			params: imageserver.Params{param: imageserver.Params{
				"foo": "bar",
			}},
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := prc.Change(tc.params)
			if c != tc.expected {
				t.Fatalf("unexpected result: got %t, want %t", c, tc.expected)
			}
		})
	}
}
