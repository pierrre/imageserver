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
	type TC struct {
		processor          *ResizeProcessor
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
			params:         imageserver.Params{resizeParam: imageserver.Params{}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  0,
				"height": 0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		// with size
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": 100,
			}},
			expectedHeight: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		// mode
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "fit",
			}},
			expectedWidth:  100,
			expectedHeight: 80,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "fill",
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
				"mode":  "fill",
			}},
			expectedWidth:  100,
			expectedHeight: 80,
		},
		// resampling
		{
			processor: &ResizeProcessor{DefaultResampling: gift.NearestNeighborResampling},
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "nearest_neighbor",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "box",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "linear",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "cubic",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "lanczos",
			}},
			expectedWidth: 100,
		},
		// error
		{
			params:             imageserver.Params{resizeParam: "invalid"},
			expectedParamError: resizeParam,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": "invalid",
			}},
			expectedParamError: resizeParam + ".width",
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": "invalid",
			}},
			expectedParamError: resizeParam + ".height",
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": -1,
			}},
			expectedParamError: resizeParam + ".width",
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": -1,
			}},
			expectedParamError: resizeParam + ".height",
		},
		{
			processor: &ResizeProcessor{MaxWidth: 500},
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 1000,
			}},
			expectedParamError: resizeParam + ".width",
		},
		{
			processor: &ResizeProcessor{MaxHeight: 500},
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": 1000,
			}},
			expectedParamError: resizeParam + ".height",
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   666,
			}},
			expectedParamError: resizeParam + ".mode",
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "invalid",
			}},
			expectedParamError: resizeParam + ".mode",
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": 666,
			}},
			expectedParamError: resizeParam + ".resampling",
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width":      100,
				"resampling": "invalid",
			}},
			expectedParamError: resizeParam + ".resampling",
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
		}()
	}
}

func TestResizeProcessorChange(t *testing.T) {
	prc := &ResizeProcessor{}
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
			params:   imageserver.Params{resizeParam: "invalid"},
			expected: true,
		},
		{
			params:   imageserver.Params{resizeParam: imageserver.Params{}},
			expected: false,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"width": 100,
			}},
			expected: true,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
				"height": 100,
			}},
			expected: true,
		},
		{
			params: imageserver.Params{resizeParam: imageserver.Params{
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
			result := prc.Change(tc.params)
			if result != tc.expected {
				t.Fatalf("unexpected result: got %t, want %t", result, tc.expected)
			}
		}()
	}
}
