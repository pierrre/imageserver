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
	type TC struct {
		processor          *Processor
		params             imageserver.Params
		expectedWidth      int
		expectedHeight     int
		expectedParamError string
	}
	for _, tc := range []TC{
		// no size
		{
			params:         imageserver.Params{},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			params:         imageserver.Params{param: imageserver.Params{}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":  0,
				"height": 0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		// with size
		{
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"height": 100,
			}},
			expectedHeight: 100,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":  100,
				"height": 100,
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		// mode
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "resize",
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		{
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
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "nearest_neighbor",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "bilinear",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "bicubic",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "mitchell_netravali",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos2",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos3",
			}},
			expectedWidth: 100,
		},
		// error
		{
			params:             imageserver.Params{param: "invalid"},
			expectedParamError: param,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width": "invalid",
			}},
			expectedParamError: param + ".width",
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"height": "invalid",
			}},
			expectedParamError: param + ".height",
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width": -1,
			}},
			expectedParamError: param + ".width",
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"height": -1,
			}},
			expectedParamError: param + ".height",
		},
		{
			processor: &Processor{MaxWidth: 500},
			params: imageserver.Params{param: imageserver.Params{
				"width": 1000,
			}},
			expectedParamError: param + ".width",
		},
		{
			processor: &Processor{MaxHeight: 500},
			params: imageserver.Params{param: imageserver.Params{
				"height": 1000,
			}},
			expectedParamError: param + ".height",
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": false,
			}},
			expectedParamError: param + ".interpolation",
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width":         100,
				"interpolation": "invalid",
			}},
			expectedParamError: param + ".interpolation",
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
				"mode":  false,
			}},
			expectedParamError: param + ".mode",
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
				"mode":  "invalid",
			}},
			expectedParamError: param + ".mode",
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
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
		}()
	}
}

func TestProcessorChange(t *testing.T) {
	prc := &Processor{}
	type TC struct {
		params   imageserver.Params
		expected bool
	}
	for _, tc := range []TC{
		{
			params:   imageserver.Params{},
			expected: false,
		},
		{
			params:   imageserver.Params{param: imageserver.Params{}},
			expected: false,
		},
		{
			params:   imageserver.Params{param: 666},
			expected: true,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"width": 100,
			}},
			expected: true,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"height": 100,
			}},
			expected: true,
		},
		{
			params: imageserver.Params{param: imageserver.Params{
				"foo": "bar",
			}},
			expected: false,
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			c := prc.Change(tc.params)
			if c != tc.expected {
				t.Fatalf("unexpected result: got %t, want %t", c, tc.expected)
			}
		}()
	}
}
