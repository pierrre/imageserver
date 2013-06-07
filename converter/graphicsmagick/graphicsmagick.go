package graphicsmagick

import (
	"fmt"
	"github.com/pierrre/imageproxy"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type GraphicsMagickConverter struct {
	Executable       string
	TempDir          string
	AllowedFormats  []string
	DefaultQualities map[string]string
}

func (converter *GraphicsMagickConverter) Convert(sourceImage *imageproxy.Image, parameters imageproxy.Parameters) (image *imageproxy.Image, err error) {
	var arguments []string

	arguments = append(arguments, "mogrify")

	arguments, _, _, err = converter.buildArgumentsSize(arguments, parameters)
	if err != nil {
		return
	}

	arguments, format, hasFileExtension, err := converter.buildArgumentsFormat(arguments, parameters, sourceImage)
	if err != nil {
		return
	}

	arguments, _, err = converter.buildArgumentsQuality(arguments, parameters, format)
	if err != nil {
		return
	}

	tempDir, err := ioutil.TempDir(converter.TempDir, "imageproxy_")
	if err != nil {
		return
	}
	defer os.RemoveAll(tempDir)

	file := filepath.Join(tempDir, "image")
	arguments = append(arguments, file)
	err = ioutil.WriteFile(file, sourceImage.Data, os.FileMode(0600))
	if err != nil {
		return
	}

	cmd := exec.Command(converter.Executable, arguments...)
	err = cmd.Run()
	if err != nil {
		return
	}

	if hasFileExtension {
		file = fmt.Sprintf("%s.%s", file, format)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	image = &imageproxy.Image{}
	image.Data = data
	image.Type = format

	return image, nil
}

func (converter *GraphicsMagickConverter) buildArgumentsSize(in []string, parameters imageproxy.Parameters) (arguments []string, width int, height int, err error) {
	arguments = in

	width, _ = parameters.GetInt("gm.width")
	if width < 0 {
		err = fmt.Errorf("Invalid width")
		return
	}

	height, _ = parameters.GetInt("gm.height")
	if height < 0 {
		err = fmt.Errorf("Invalid height")
		return
	}

	if width != 0 || height != 0 {
		widthString := ""
		if width != 0 {
			widthString = strconv.Itoa(width)
		}
		heightString := ""
		if height != 0 {
			heightString = strconv.Itoa(height)
		}
		arguments = append(arguments, "-resize", fmt.Sprintf("%sx%s", widthString, heightString))
	}

	return
}

func (converter *GraphicsMagickConverter) buildArgumentsFormat(in []string, parameters imageproxy.Parameters, sourceImage *imageproxy.Image) (arguments []string, format string, hasFileExtension bool, err error) {
	arguments = in

	format, _ = parameters.GetString("gm.format")

	formatSpecified := true
	if len(format) == 0 {
		format = sourceImage.Type
		formatSpecified = false
	}

	err = converter.validateFormat(format)
	if err != nil {
		return
	}

	if formatSpecified {
		arguments = append(arguments, "-format", format)
	}

	hasFileExtension = formatSpecified

	return
}

func (converter *GraphicsMagickConverter) validateFormat(format string) (err error) {
	if converter.AllowedFormats != nil {
		ok := false
		for _, f := range converter.AllowedFormats {
			if f == format {
				ok = true
				break
			}
		}
		if !ok {
			err = fmt.Errorf("Invalid format")
			return
		}
	}
	return
}

func (converter *GraphicsMagickConverter) buildArgumentsQuality(in []string, parameters imageproxy.Parameters, format string) (arguments []string, quality string, err error) {
	arguments = in

	quality, _ = parameters.GetString("gm.quality")

	if len(quality) == 0 && converter.DefaultQualities != nil {
		if q, ok := converter.DefaultQualities[format]; ok {
			quality = q
		}
	}

	if len(quality) > 0 {
		err = converter.validateQuality(quality, format)
		if err != nil {
			return
		}
		arguments = append(arguments, "-quality", quality)
	}

	return
}

func (converter *GraphicsMagickConverter) validateQuality(quality string, format string) (err error) {
	qualityInt, err := strconv.Atoi(quality)
	if err != nil {
		return
	}

	if qualityInt < 0 {
		err = fmt.Errorf("Invalid quality")
		return
	}

	if format == "jpeg" {
		if qualityInt < 0 || qualityInt > 100 {
			err = fmt.Errorf("Invalid quality")
			return
		}
	}

	return
}
