package nfntresize

import (
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &ValidateParamsServer{}

func TestValidateParamsServer(t *testing.T) {
	type TC struct {
		widthMax           uint
		heightMax          uint
		params             imageserver.Params
		expectedParamError string
	}
	for _, tc := range []TC{
		{
			params: imageserver.Params{},
		},
		{
			params: imageserver.Params{Param: imageserver.Params{}},
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"width": 1024,
			}},
		},
		{
			params: imageserver.Params{Param: imageserver.Params{
				"height": 768,
			}},
		},
		{
			widthMax: 1024,
			params: imageserver.Params{Param: imageserver.Params{
				"width": 1024,
			}},
		},
		{
			heightMax: 768,
			params: imageserver.Params{Param: imageserver.Params{
				"height": 768,
			}},
		},
		{
			widthMax: 1024,
			params: imageserver.Params{Param: imageserver.Params{
				"width": 9001,
			}},
			expectedParamError: Param + ".width",
		},
		{
			heightMax: 768,
			params: imageserver.Params{Param: imageserver.Params{
				"height": 9001,
			}},
			expectedParamError: Param + ".height",
		},
		{
			params:             imageserver.Params{Param: "invalid"},
			expectedParamError: Param,
		},
		{
			widthMax: 1024,
			params: imageserver.Params{Param: imageserver.Params{
				"width": "invalid",
			}},
			expectedParamError: Param + ".width",
		},
		{
			heightMax: 768,
			params: imageserver.Params{Param: imageserver.Params{
				"height": "invalid",
			}},
			expectedParamError: Param + ".height",
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			srv := &ValidateParamsServer{
				Server: &imageserver.StaticServer{
					Image: testdata.Medium,
				},
			}

			if tc.widthMax != 0 {
				srv.WidthMax = tc.widthMax
			}
			if tc.heightMax != 0 {
				srv.HeightMax = tc.heightMax
			}
			_, err := srv.Get(tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && err.Param == tc.expectedParamError {
					return
				}
				t.Fatal(err)
			}
		}()
	}
}
