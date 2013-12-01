package testdata

import (
	"github.com/pierrre/imageserver"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
	Small    = loadImage("small.jpg", "jpeg")
	Medium   = loadImage("medium.jpg", "jpeg")
	Large    = loadImage("large.jpg", "jpeg")
	Huge     = loadImage("huge.jpg", "jpeg")
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
