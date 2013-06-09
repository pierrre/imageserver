package graphicsmagick

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type GraphicsMagickConverter struct {
	Executable       string
	TempDir          string
	AllowedFormats   []string
	DefaultQualities map[string]string
}

func (converter *GraphicsMagickConverter) Convert(sourceImage *imageserver.Image, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	var arguments []string

	arguments = append(arguments, "mogrify")

	arguments, width, height, err := converter.buildArgumentsResize(arguments, parameters)
	if err != nil {
		return
	}

	arguments, err = converter.buildArgumentsBackground(arguments, parameters)
	if err != nil {
		return
	}

	arguments, err = converter.buildArgumentsExtent(arguments, parameters, width, height)
	if err != nil {
		return
	}

	arguments, format, hasFileExtension, err := converter.buildArgumentsFormat(arguments, parameters, sourceImage)
	if err != nil {
		return
	}

	arguments, err = converter.buildArgumentsQuality(arguments, parameters, format)
	if err != nil {
		return
	}

	tempDir, err := ioutil.TempDir(converter.TempDir, "imageserver_")
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

	image = &imageserver.Image{}
	image.Data = data
	image.Type = format

	return image, nil
}

func (converter *GraphicsMagickConverter) buildArgumentsResize(in []string, parameters imageserver.Parameters) (arguments []string, width int, height int, err error) {
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
		resize := fmt.Sprintf("%sx%s", widthString, heightString)

		if fill, _ := parameters.GetBool("gm.fill"); fill {
			resize = resize + "^"
		}

		if ignoreRatio, _ := parameters.GetBool("gm.ignore_ratio"); ignoreRatio {
			resize = resize + "!"
		}

		if onlyShrinkLarger, _ := parameters.GetBool("gm.only_shrink_larger"); onlyShrinkLarger {
			resize = resize + ">"
		}

		if onlyEnlargeSmaller, _ := parameters.GetBool("gm.only_enlarge_smaller"); onlyEnlargeSmaller {
			resize = resize + "<"
		}

		arguments = append(arguments, "-resize", resize)
	}

	return
}

func (converter *GraphicsMagickConverter) buildArgumentsBackground(in []string, parameters imageserver.Parameters) (arguments []string, err error) {
	arguments = in

	background, _ := parameters.GetString("gm.background")

	if backgroundLength := len(background); backgroundLength > 0 {
		if backgroundLength != 6 && backgroundLength != 8 && backgroundLength != 3 && backgroundLength != 4 {
			err = fmt.Errorf("Invalid background")
			return
		}

		for _, r := range background {
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
				err = fmt.Errorf("Invalid background")
				return
			}
		}

		arguments = append(arguments, "-background", fmt.Sprintf("#%s", background))
	}

	return
}

func (converter *GraphicsMagickConverter) buildArgumentsExtent(in []string, parameters imageserver.Parameters, width int, height int) (arguments []string, err error) {
	arguments = in

	if width != 0 && height != 0 {
		if extent, _ := parameters.GetBool("gm.extent"); extent {
			arguments = append(arguments, "-gravity", "center")
			arguments = append(arguments, "-extent", fmt.Sprintf("%dx%d", width, height))
		}
	}

	return
}

func (converter *GraphicsMagickConverter) buildArgumentsFormat(in []string, parameters imageserver.Parameters, sourceImage *imageserver.Image) (arguments []string, format string, hasFileExtension bool, err error) {
	arguments = in

	format, _ = parameters.GetString("gm.format")

	formatSpecified := true
	if len(format) == 0 {
		format = sourceImage.Type
		formatSpecified = false
	}

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

	if formatSpecified {
		arguments = append(arguments, "-format", format)
	}

	hasFileExtension = formatSpecified

	return
}

func (converter *GraphicsMagickConverter) buildArgumentsQuality(in []string, parameters imageserver.Parameters, format string) (arguments []string, err error) {
	arguments = in

	quality, _ := parameters.GetString("gm.quality")

	if len(quality) == 0 && converter.DefaultQualities != nil {
		if q, ok := converter.DefaultQualities[format]; ok {
			quality = q
		}
	}

	if len(quality) > 0 {
		qualityInt, e := strconv.Atoi(quality)
		if e != nil {
			err = e
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
		arguments = append(arguments, "-quality", quality)
	}

	return
}
