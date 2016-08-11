package file

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestServerGet(t *testing.T) {
	srv := &Server{
		Root: testdata.Dir,
	}
	for _, tc := range []struct {
		name               string
		params             imageserver.Params
		expectedParamError string
		expectedImage      *imageserver.Image
	}{
		{
			name: "Normal",
			params: imageserver.Params{
				imageserver_source.Param: testdata.MediumFileName,
			},
			expectedImage: testdata.Medium,
		},
		{
			name: "CleanSource",
			params: imageserver.Params{
				imageserver_source.Param: "../" + testdata.MediumFileName,
			},
			expectedImage: testdata.Medium,
		},
		{
			name:               "ErrorNoParam",
			params:             imageserver.Params{},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "ErrorNotFound",
			params: imageserver.Params{
				imageserver_source.Param: "invalid",
			},
			expectedParamError: imageserver_source.Param,
		},
		{
			name: "ErrorIdentify",
			params: imageserver.Params{
				imageserver_source.Param: "testdata.go",
			},
			expectedParamError: imageserver_source.Param,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			im, err := srv.Get(tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && err.Param == tc.expectedParamError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedParamError != "" {
				t.Fatal("no error")
			}
			if im == nil {
				t.Fatal("no image")
			}
			if im.Format != tc.expectedImage.Format {
				t.Fatalf("unexpected image format: got \"%s\", want \"%s\"", im.Format, tc.expectedImage.Format)
			}
			if !bytes.Equal(im.Data, tc.expectedImage.Data) {
				t.Fatal("data not equal")
			}
		})
	}
}

func TestServerGetPath(t *testing.T) {
	srv := &Server{
		Root: "root",
	}
	for _, tc := range []struct {
		name     string
		source   string
		expected string
	}{
		{"Normal", "file", "root/file"},
		{"SlashBefore", "/file", "root/file"},
		{"SlashAfter", "file/", "root/file"},
		{"MultipleSlash", "///file", "root/file"},
		{"Up", "../file", "root/file"},
		{"UpDown", "../dir/file", "root/dir/file"},
		{"DownUp", "dir/../file", "root/file"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pth, err := srv.getPath(imageserver.Params{
				imageserver_source.Param: tc.source,
			})
			if err != nil {
				t.Fatal(err)
			}
			if pth != tc.expected {
				t.Fatalf("unexpected result: got \"%s\", want \"%s\"", pth, tc.expected)
			}
		})
	}
}

func TestIdentifyMime(t *testing.T) {
	for _, tc := range []struct {
		name           string
		filename       string
		data           []byte
		expectedFormat string
		expectedError  bool
	}{
		{
			name:           "JPEG",
			filename:       testdata.MediumFileName,
			data:           testdata.Medium.Data,
			expectedFormat: testdata.Medium.Format,
			expectedError:  false,
		},
		{
			name:           "PNG",
			filename:       testdata.RandomFileName,
			data:           testdata.Random.Data,
			expectedFormat: testdata.Random.Format,
			expectedError:  false,
		},
		{
			name:           "GIF",
			filename:       testdata.AnimatedFileName,
			data:           testdata.Animated.Data,
			expectedFormat: testdata.Animated.Format,
			expectedError:  false,
		},
		{
			name:          "ErrorNoExtension",
			filename:      "invalid",
			data:          testdata.Medium.Data,
			expectedError: true,
		},
		{
			name:          "ErrorUnknownExtension",
			filename:      "invalid.invalid",
			data:          testdata.Medium.Data,
			expectedError: true,
		},
		{
			name:          "ErrorUnsupportedType",
			filename:      "invalid.txt",
			data:          testdata.Medium.Data,
			expectedError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			format, err := IdentifyMime(filepath.Join(testdata.Dir, tc.filename), tc.data)
			if err != nil {
				if tc.expectedError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedError {
				t.Fatal("no error")
			}
			if format != tc.expectedFormat {
				t.Fatalf("unexpected format: got %s, want %s", format, tc.expectedFormat)
			}
		})
	}
}
