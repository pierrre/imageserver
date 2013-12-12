// Package testdata provides test images
package testdata

import (
	"github.com/pierrre/imageserver"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
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
)

func loadImage(filename string, format string) *imageserver.Image {
	_, currentFile, _, _ := runtime.Caller(0)
	dataPath := filepath.Dir(currentFile)
	data, err := ioutil.ReadFile(filepath.Join(dataPath, filename))
	if err != nil {
		panic(err)
	}
	return &imageserver.Image{
		Format: format,
		Data:   data,
	}
}
