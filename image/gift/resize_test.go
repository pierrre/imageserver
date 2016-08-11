package gift

import (
	"testing"

	"github.com/disintegration/gift"
	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/jpeg"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

var _ imageserver_image.Processor = &ResizeProcessor{}

func TestResizeProcessorProcess(t *testing.T) {
	nim, err := imageserver_image.Decode(imageserver_testdata.Medium)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name               string
		processor          *ResizeProcessor
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
			name:           "EmptyParam",
			params:         imageserver.Params{resizeParam: imageserver.Params{}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			name: "SizeZero",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  0,
				"height": 0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		// with size
		{
			name: "Width",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			name: "Height",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": 100,
			}},
			expectedHeight: 100,
		},
		{
			name: "WidthHeight",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		// mode
		{
			name: "ModeFit",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "fit",
			}},
			expectedWidth:  100,
			expectedHeight: 80,
		},
		{
			name: "ModeFillWidthHeight",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "fill",
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		{
			name: "ModeFill",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
				"mode":  "fill",
			}},
			expectedWidth:  100,
			expectedHeight: 80,
		},
		// resampling
		{
			name:      "ResamplingDefault",
			processor: &ResizeProcessor{DefaultResampling: gift.NearestNeighborResampling},
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			name: "ResamplingNearestNeighbor",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "nearest_neighbor",
			}},
			expectedWidth: 100,
		},
		{
			name: "ResamplingBox",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "box",
			}},
			expectedWidth: 100,
		},
		{
			name: "ResamplingLinear",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "linear",
			}},
			expectedWidth: 100,
		},
		{
			name: "ResamplingCubic",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "cubic",
			}},
			expectedWidth: 100,
		},
		{
			name: "ResamplingLanczos",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "lanczos",
			}},
			expectedWidth: 100,
		},
		// error
		{
			name:               "ParamInvalid",
			params:             imageserver.Params{resizeParam: "invalid"},
			expectedParamError: resizeParam,
		},
		{
			name: "WidthInvalidType",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": "invalid",
			}},
			expectedParamError: resizeParam + ".width",
		},
		{
			name: "HeightInvalidInvalid",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": "invalid",
			}},
			expectedParamError: resizeParam + ".height",
		},
		{
			name: "WidthInvalidNegative",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": -1,
			}},
			expectedParamError: resizeParam + ".width",
		},
		{
			name: "HeightInvalidNegative",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": -1,
			}},
			expectedParamError: resizeParam + ".height",
		},
		{
			name:      "WidthInvalidTooLarge",
			processor: &ResizeProcessor{MaxWidth: 500},
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 1000,
			}},
			expectedParamError: resizeParam + ".width",
		},
		{
			name:      "HeightInvalidTooLarge",
			processor: &ResizeProcessor{MaxHeight: 500},
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": 1000,
			}},
			expectedParamError: resizeParam + ".height",
		},
		{
			name: "ModeInvalidType",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   666,
			}},
			expectedParamError: resizeParam + ".mode",
		},
		{
			name: "ModeInvalidUnknown",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "invalid",
			}},
			expectedParamError: resizeParam + ".mode",
		},
		{
			name: "ResamplingInvalidType",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": 666,
			}},
			expectedParamError: resizeParam + ".resampling",
		},
		{
			name: "ResamplingInvalidUnknown",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "invalid",
			}},
			expectedParamError: resizeParam + ".resampling",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prc := tc.processor
			if prc == nil {
				prc = &ResizeProcessor{}
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

func TestResizeProcessorChange(t *testing.T) {
	prc := &ResizeProcessor{}
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
			name:     "ParamInvalid",
			params:   imageserver.Params{resizeParam: "invalid"},
			expected: true,
		},
		{
			name:     "EmptyParam",
			params:   imageserver.Params{resizeParam: imageserver.Params{}},
			expected: false,
		},
		{
			name: "Width",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
			expected: true,
		},
		{
			name: "Invalid",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": 100,
			}},
			expected: true,
		},
		{
			name: "ParamUnknown",
			params: imageserver.Params{resizeParam: imageserver.Params{
				"foo": "bar",
			}},
			expected: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := prc.Change(tc.params)
			if result != tc.expected {
				t.Fatalf("unexpected result: got %t, want %t", result, tc.expected)
			}
		})
	}
}
