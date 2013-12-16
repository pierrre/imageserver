// Package testdata provides test images
package testdata

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
	// Images contains all images by filename
	Images = make(map[string]*imageserver.Image)

	// Small is a small Image
	Small = loadImage("small.jpg", "jpeg")
	// Medium is a medium Image
	Medium = loadImage("medium.jpg", "jpeg")
	// Large is a large image
	Large = loadImage("large.jpg", "jpeg")
	// Huge is a huge image
	Huge = loadImage("huge.jpg", "jpeg")
	// Animated is an animated GIF Image
	Animated = loadImage("animated.gif", "gif")

	// Provider is an Image Provider that uses filename as source
	Provider = new(testDataProvider)
)

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
	_, currentFile, _, _ := runtime.Caller(0)
	dataPath := filepath.Dir(currentFile)
	filePath := filepath.Join(dataPath, filename)

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
