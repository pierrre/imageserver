package gift

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

var _ imageserver_image.Processor = &RotateProcessor{}

func TestRotateProcessorProcess(t *testing.T) {
	nim, err := imageserver_image.Decode(imageserver_testdata.Medium)
	if err != nil {
		t.Fatal(err)
	}
	prc := &RotateProcessor{}
	for _, tc := range []struct {
		name               string
		params             imageserver.Params
		expectedWidth      int
		expectedHeight     int
		expectedParamError string
	}{
		// no rotation
		{
			name:           "Empty",
			params:         imageserver.Params{},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			name:           "EmptyParam",
			params:         imageserver.Params{rotateParam: imageserver.Params{}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		// with rotation
		{
			name: "Rotation0",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 0.0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			name: "Rotation360",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 360.0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			name: "Rotation90",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 90.0,
			}},
			expectedWidth:  819,
			expectedHeight: 1024,
		},
		{
			name: "Rotation180",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 180.0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			name: "Rotation270",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 270.0,
			}},
			expectedWidth:  819,
			expectedHeight: 1024,
		},
		{
			name: "Rotation45",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 45.0,
			}},
			expectedWidth:  1304,
			expectedHeight: 1304,
		},
		// background
		{
			name: "Background",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":   45.0,
				"background": "FF0000",
			}},
		},
		// interpolation
		{
			name: "InterpolationNearestNeighbor",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":      45.0,
				"interpolation": "nearest_neighbor",
			}},
		},
		{
			name: "InterpolationLinear",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":      45.0,
				"interpolation": "linear",
			}},
		},
		{
			name: "InterpolationCubic",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":      45.0,
				"interpolation": "cubic",
			}},
		},
		// error
		{
			name:               "ParamInvalid",
			params:             imageserver.Params{rotateParam: "invalid"},
			expectedParamError: rotateParam,
		},
		{
			name: "RotationInvalid",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": "invalid",
			}},
			expectedParamError: rotateParam + ".rotation",
		},
		{
			name: "BackgroundInvalidColor",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":   45.0,
				"background": "invalid",
			}},
			expectedParamError: rotateParam + ".background",
		},
		{
			name: "BackgroundInvalidType",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":   45.0,
				"background": 666,
			}},
			expectedParamError: rotateParam + ".background",
		},
		{
			name: "InterpolationInvalidUnknown",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":      45.0,
				"interpolation": "invalid",
			}},
			expectedParamError: rotateParam + ".interpolation",
		},
		{
			name: "InteprolationInvalidType",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation":      45.0,
				"interpolation": 666,
			}},
			expectedParamError: rotateParam + ".interpolation",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
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

func TestRotateProcessorGetRotation(t *testing.T) {
	prc := &RotateProcessor{}
	for _, tc := range []struct {
		name     string
		params   imageserver.Params
		expected float32
	}{
		{
			name:     "Empty",
			params:   imageserver.Params{},
			expected: 0,
		},
		{
			name: "0",
			params: imageserver.Params{
				"rotation": 0.0,
			},
			expected: 0,
		},
		{
			name: "10",
			params: imageserver.Params{
				"rotation": 10.0,
			},
			expected: 10,
		},
		{
			name: "-10",
			params: imageserver.Params{
				"rotation": -10.0,
			},
			expected: 350,
		},
		{
			name: "370",
			params: imageserver.Params{
				"rotation": 370.0,
			},
			expected: 10,
		},
		{
			name: "360",
			params: imageserver.Params{
				"rotation": 360.0,
			},
			expected: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, err := prc.getRotation(tc.params)
			if err != nil {
				t.Fatal(err)
			}
			if res != tc.expected {
				t.Fatalf("unexpected result for %#v: got %f, want %f", tc.params, res, tc.expected)
			}
		})
	}
}

func TestRotateProcessorGetRotationError(t *testing.T) {
	_, err := (&RotateProcessor{}).getRotation(imageserver.Params{"rotation": "invalid"})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestRotateProcessorChange(t *testing.T) {
	prc := &RotateProcessor{}
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
			params:   imageserver.Params{rotateParam: "invalid"},
			expected: true,
		},
		{
			name:     "ParamEmpty",
			params:   imageserver.Params{rotateParam: imageserver.Params{}},
			expected: false,
		},
		{
			name: "Rotation",
			params: imageserver.Params{rotateParam: imageserver.Params{
				"rotation": 45.0,
			}},
			expected: true,
		},
		{
			name: "ParamUnknown",
			params: imageserver.Params{rotateParam: imageserver.Params{
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

func TestParseHexColor(t *testing.T) {
	for _, tc := range []struct {
		hex      string
		expected color.Color
	}{
		{
			hex:      "F84",
			expected: color.NRGBA{R: 0xff, G: 0x88, B: 0x44, A: 0xff},
		},
		{
			hex:      "F842",
			expected: color.NRGBA{R: 0x88, G: 0x44, B: 0x22, A: 0xff},
		},
		{
			hex:      "FC8642",
			expected: color.NRGBA{R: 0xfc, G: 0x86, B: 0x42, A: 0xff},
		},
		{
			hex:      "FC864210",
			expected: color.NRGBA{R: 0x86, G: 0x42, B: 0x10, A: 0xfc},
		},
	} {
		t.Run(tc.hex, func(t *testing.T) {
			res, err := parseHexColor(tc.hex)
			if err != nil {
				t.Fatal(err)
			}
			if res != tc.expected {
				t.Fatalf("unexpected result for \"%s\": got %#v, want %#v", tc.hex, res, tc.expected)
			}
		})
	}
}

func TestParseHexColorError(t *testing.T) {
	for _, hex := range []string{
		"0000000000",
		"zzz",
		"0000000",
	} {
		t.Run(hex, func(t *testing.T) {
			_, err := parseHexColor(hex)
			if err == nil {
				t.Fatalf("no error for \"%s\"", hex)
			}
		})
	}
}

func TestHexStringToInts(t *testing.T) {
	for _, tc := range []struct {
		hex      string
		expected []uint8
	}{
		{
			hex:      "",
			expected: nil,
		},
		{
			hex:      "1",
			expected: []uint8{0x1},
		},
		{
			hex: "0123456789abcdefABCDEF",
			expected: []uint8{
				0x0, 0x1, 0x2, 0x3, 0x4,
				0x5, 0x6, 0x7, 0x8, 0x9,
				0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
				0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
			},
		},
	} {
		t.Run(tc.hex, func(t *testing.T) {
			res, err := hexStringToInts(tc.hex)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(res, tc.expected) {
				t.Fatalf("unexpected result for \"%s\": got %#v, want %#v", tc.hex, res, tc.expected)
			}
		})
	}
}

func TestHexStringToIntsError(t *testing.T) {
	_, err := hexStringToInts("zzz")
	if err == nil {
		t.Fatal("no error")
	}
}
