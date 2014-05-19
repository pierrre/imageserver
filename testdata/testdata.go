// Package testdata provides test images
package testdata

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/pierrre/imageserver"
)

var (
	Dir = initDir()

	// Images contains all images by filename
	Images = make(map[string]*imageserver.Image)

	// Small is a small Image
	SmallFileName = "small.jpg"
	Small         = loadImage(SmallFileName, "jpeg")
	// Medium is a medium Image
	MediumFileName = "medium.jpg"
	Medium         = loadImage(MediumFileName, "jpeg")
	// Large is a large image
	LargeFileName = "large.jpg"
	Large         = loadImage(LargeFileName, "jpeg")
	// Huge is a huge image
	HugeFileName = "huge.jpg"
	Huge         = loadImage(HugeFileName, "jpeg")
	// Animated is an animated GIF Image
	AnimatedFileName = "animated.gif"
	Animated         = loadImage(AnimatedFileName, "gif")

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
		return nil, imageserver.NewError(fmt.Sprintf("source is not a string: %#v", source))
	}

	image, ok := Images[name]
	if !ok {
		return nil, imageserver.NewError(fmt.Sprintf("source is unknown: %s", name))
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
