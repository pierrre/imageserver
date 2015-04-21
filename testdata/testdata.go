// Package testdata provides test images
package testdata

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/pierrre/imageserver"
)

var (
	// Dir is the path to the directory containing the test data
	Dir = initDir()

	// Images contains all images by filename
	Images = make(map[string]*imageserver.Image)

	// SmallFileName is the file name of Small
	SmallFileName = "small.jpg"
	// Small is a small Image
	Small = loadImage(SmallFileName, "jpeg")
	// MediumFileName is the file name of Medium
	MediumFileName = "medium.jpg"
	// Medium is a medium Image
	Medium = loadImage(MediumFileName, "jpeg")
	// LargeFileName is the file name of Large
	LargeFileName = "large.jpg"
	// Large is a large image
	Large = loadImage(LargeFileName, "jpeg")
	// HugeFileName is the file name of Huge
	HugeFileName = "huge.jpg"
	// Huge is a huge image
	Huge = loadImage(HugeFileName, "jpeg")
	// AnimatedFileName is the file name of Animated
	AnimatedFileName = "animated.gif"
	// Animated is an animated GIF Image
	Animated = loadImage(AnimatedFileName, "gif")
	// DalaiGammaFileName is the file name of DalaiGamma
	DalaiGammaFileName = "dalai_gamma.jpg"
	// DalaiGamma is a gamma test Image (from http://www.4p8.com/eric.brasseur/gamma.html)
	DalaiGamma = loadImage(DalaiGammaFileName, "jpeg")
	// InvalidFileName is the file name of Invalid
	InvalidFileName = "invalid.jpg"
	// Invalid is an invalid Image
	Invalid = loadImage(InvalidFileName, "jpeg")

	// Provider is an Image Provider that uses filename as source
	Server imageserver.Server = imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
		source, err := params.Get(imageserver.SourceParam)
		if err != nil {
			return nil, err
		}
		name, ok := source.(string)
		if !ok {
			return nil, &imageserver.ParamError{Param: imageserver.SourceParam, Message: "not a string"}
		}
		im, ok := Images[name]
		if !ok {
			return nil, &imageserver.ParamError{Param: imageserver.SourceParam, Message: "unknown image"}
		}
		return im, nil
	})
)

func initDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}

func loadImage(filename string, format string) *imageserver.Image {
	filePath := filepath.Join(Dir, filename)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	im := &imageserver.Image{
		Format: format,
		Data:   data,
	}
	Images[filename] = im
	return im
}
