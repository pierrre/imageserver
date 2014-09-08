// Package testdata provides test images
package testdata

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/pierrre/imageserver"
	imageserver_provider "github.com/pierrre/imageserver/provider"
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

	// Provider is an Image Provider that uses filename as source
	Provider = new(testDataProvider)
)

func initDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}

type testDataProvider struct{}

func (provider *testDataProvider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	name, ok := source.(string)
	if !ok {
		return nil, &imageserver_provider.SourceError{Message: "not a string"}
	}

	image, ok := Images[name]
	if !ok {
		return nil, &imageserver_provider.SourceError{Message: "unknown image"}
	}

	return image, nil
}

func loadImage(filename string, format string) *imageserver.Image {
	filePath := filepath.Join(Dir, filename)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	image := &imageserver.Image{
		Format: format,
		Data:   data,
	}

	Images[filename] = image

	return image
}
