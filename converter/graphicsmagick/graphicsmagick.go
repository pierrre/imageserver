package graphicsmagick

import (
	"fmt"
	"github.com/pierrre/imageproxy"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type GraphicsMagickConverter struct {
	Executable string
	TempDir    string
}

func (converter *GraphicsMagickConverter) Convert(sourceImage *imageproxy.Image, parameters imageproxy.Parameters) (image *imageproxy.Image, err error) {
	tempDir, err := ioutil.TempDir(converter.TempDir, "imageproxy_")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "image")
	err = ioutil.WriteFile(filePath, sourceImage.Data, os.FileMode(0600))
	if err != nil {
		return
	}

	var arguments []string
	arguments = append(arguments, "mogrify")
	width, _ := parameters.GetInt("width")
	height, _ := parameters.GetInt("height")
	if width != 0 && height != 0 {
		if width <= 0 {
			err = fmt.Errorf("Invalid width")
		}
		if height <= 0 {
			err = fmt.Errorf("Invalid height")
		}
		arguments = append(arguments, "-resize", fmt.Sprintf("%dx%d", width, height))
	}
	arguments = append(arguments, filePath)

	cmd := exec.Command(converter.Executable, arguments...)
	err = cmd.Run()
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	image = &imageproxy.Image{}
	image.Data = data
	//FIX type

	return image, nil
}
