// Package testdata provides test images.
package testdata

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
)

var (
	// Dir is the path to the directory containing the test data.
	Dir = initDir()

	// Images contains all images by filename.
	Images = make(map[string]*imageserver.Image)

	// SmallFileName is the file name of Small.
	SmallFileName = "small.jpg"
	// Small is a small Image (from https://www.flickr.com/photos/maradentro/3600833235/).
	Small = loadImage(SmallFileName, "jpeg")

	// MediumFileName is the file name of Medium.
	MediumFileName = "medium.jpg"
	// Medium is a medium Image (from https://www.flickr.com/photos/doug88888/5793867021/).
	Medium = loadImage(MediumFileName, "jpeg")

	// LargeFileName is the file name of Large.
	LargeFileName = "large.jpg"
	// Large is a large image (from https://www.flickr.com/photos/doug88888/4130990745/).
	Large = loadImage(LargeFileName, "jpeg")

	// HugeFileName is the file name of Huge.
	HugeFileName = "huge.jpg"
	// Huge is a huge image.
	Huge = loadImage(HugeFileName, "jpeg")

	// AnimatedFileName is the file name of Animated.
	AnimatedFileName = "animated.gif"
	// Animated is an animated GIF Image.
	Animated = loadImage(AnimatedFileName, "gif")

	// SpaceshipFileName is the file name of Spaceship.
	SpaceshipFileName = "spaceship.gif"
	// Spaceship is an animated spaceship GIF Image.
	Spaceship = loadImage(SpaceshipFileName, "gif")

	// DalaiGammaFileName is the file name of DalaiGamma.
	DalaiGammaFileName = "dalai_gamma.jpg"
	// DalaiGamma is a gamma test Image (from http://www.ericbrasseur.org/gamma.html).
	DalaiGamma = loadImage(DalaiGammaFileName, "jpeg")

	// GraySquaresFileName is the file name of GraySquares.
	GraySquaresFileName = "gray_squares.jpg"
	// GraySquares is a gamma test Image (from http://www.ericbrasseur.org/gamma.html)
	GraySquares = loadImage(GraySquaresFileName, "jpeg")

	// RulesSucksFileName is the file name of RulesSucks.
	RulesSucksFileName = "rules_sucks.png"
	// RulesSucks is a gamma test Image (from http://www.ericbrasseur.org/gamma.html)
	RulesSucks = loadImage(RulesSucksFileName, "png")

	// RingsFileName is the file name of Rings.
	RingsFileName = "rings.png"
	// Rings is a moiré test Image
	Rings = loadImage(RingsFileName, "png")

	// RandomFileName is the file name of Random.
	RandomFileName = "random.png"
	// Random is a random Image.
	Random = loadImage(RandomFileName, "png")

	// InvalidFileName is the file name of Invalid.
	InvalidFileName = "invalid.jpg"
	// Invalid is an invalid Image.
	Invalid = loadImage(InvalidFileName, "invalid")

	// Server is an Image Server that uses filename as source.
	Server = imageserver.Server(imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
		source, err := params.GetString(imageserver_source.Param)
		if err != nil {
			return nil, err
		}
		im, err := Get(source)
		if err != nil {
			return nil, &imageserver.ParamError{Param: imageserver_source.Param, Message: err.Error()}
		}
		return im, nil
	}))
)

// Get returns an Image for a name.
func Get(name string) (*imageserver.Image, error) {
	im, ok := Images[name]
	if !ok {
		return nil, fmt.Errorf("unknown image \"%s\"", name)
	}
	return im, nil
}

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
