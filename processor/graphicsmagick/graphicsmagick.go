// GraphicsMagick processor
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

const tempDirPrefix = "imageserver_"

// Processes an image with GraphicsMagick command line (mogrify command)
//
// All parameters are prefixed with "gm." and are optionals.
//
// See GraphicsMagick documentation for more information about arguments.
//
// Parameters
//
// - width / height: sizes for "-resize" argument (both optionals)
//
// - fill: "^" for "-resize" argument
//
// - ignore_ratio: "!" for "-resize" argument
//
// - only_shrink_larger: ">" for "-resize" argument
//
// - only_enlarge_smaller: "<" for "-resize" argument
//
// - background: color for "-background" argument, 3/4/6/8 lower case hexadecimal characters
//
// - extent: "-extent" parameter, uses width/height parameters and add "-gravity center" argument
//
// - format: "-format" parameter
//
// - quality: "-quality" parameter
type GraphicsMagickProcessor struct {
	Executable string // path to "gm" executable, usually "/usr/bin/gm"

	TempDir          string            // temp directory for image files, optional
	AllowedFormats   []string          // allowed format list, optional
	DefaultQualities map[string]string // default qualities by format, optional
}

func (processor *GraphicsMagickProcessor) Process(sourceImage *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	var arguments []string

	arguments = append(arguments, "mogrify")

	arguments, width, height, err := processor.buildArgumentsResize(arguments, parameters)
	if err != nil {
		return nil, err
	}

	arguments, err = processor.buildArgumentsBackground(arguments, parameters)
	if err != nil {
		return nil, err
	}

	arguments, err = processor.buildArgumentsExtent(arguments, parameters, width, height)
	if err != nil {
		return nil, err
	}

	arguments, format, hasFileExtension, err := processor.buildArgumentsFormat(arguments, parameters, sourceImage)
	if err != nil {
		return nil, err
	}

	arguments, err = processor.buildArgumentsQuality(arguments, parameters, format)
	if err != nil {
		return nil, err
	}

	if len(arguments) == 1 {
		return sourceImage, nil
	}

	tempDir, err := ioutil.TempDir(processor.TempDir, tempDirPrefix)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	file := filepath.Join(tempDir, "image")
	arguments = append(arguments, file)
	err = ioutil.WriteFile(file, sourceImage.Data, os.FileMode(0600))
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(processor.Executable, arguments...)
	err = cmd.Run()
	if err != nil {
		return nil, imageserver.NewError("Error during execution of GraphicsMagick")
	}

	if hasFileExtension {
		file = fmt.Sprintf("%s.%s", file, format)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	image := &imageserver.Image{}
	image.Data = data
	image.Type = format

	return image, nil
}

func (processor *GraphicsMagickProcessor) buildArgumentsResize(in []string, parameters imageserver.Parameters) (arguments []string, width int, height int, err error) {
	arguments = in

	width, _ = parameters.GetInt("gm.width")
	if width < 0 {
		err = imageserver.NewError("Invalid width parameter")
		return
	}

	height, _ = parameters.GetInt("gm.height")
	if height < 0 {
		err = imageserver.NewError("Invalid height parameter")
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

func (processor *GraphicsMagickProcessor) buildArgumentsBackground(arguments []string, parameters imageserver.Parameters) ([]string, error) {
	background, _ := parameters.GetString("gm.background")

	if backgroundLength := len(background); backgroundLength > 0 {
		if backgroundLength != 6 && backgroundLength != 8 && backgroundLength != 3 && backgroundLength != 4 {
			return nil, imageserver.NewError("Invalid background parameter")
		}

		for _, r := range background {
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
				return nil, imageserver.NewError("Invalid background parameter")
			}
		}

		arguments = append(arguments, "-background", fmt.Sprintf("#%s", background))
	}

	return arguments, nil
}

func (processor *GraphicsMagickProcessor) buildArgumentsExtent(arguments []string, parameters imageserver.Parameters, width int, height int) ([]string, error) {
	if width != 0 && height != 0 {
		if extent, _ := parameters.GetBool("gm.extent"); extent {
			arguments = append(arguments, "-gravity", "center")
			arguments = append(arguments, "-extent", fmt.Sprintf("%dx%d", width, height))
		}
	}

	return arguments, nil
}

func (processor *GraphicsMagickProcessor) buildArgumentsFormat(in []string, parameters imageserver.Parameters, sourceImage *imageserver.Image) (arguments []string, format string, hasFileExtension bool, err error) {
	arguments = in

	format, _ = parameters.GetString("gm.format")

	formatSpecified := true
	if len(format) == 0 {
		format = sourceImage.Type
		formatSpecified = false
	}

	if processor.AllowedFormats != nil {
		ok := false
		for _, f := range processor.AllowedFormats {
			if f == format {
				ok = true
				break
			}
		}
		if !ok {
			err = imageserver.NewError("Invalid format parameter")
			return
		}
	}

	if formatSpecified {
		arguments = append(arguments, "-format", format)
	}

	hasFileExtension = formatSpecified

	return
}

func (processor *GraphicsMagickProcessor) buildArgumentsQuality(arguments []string, parameters imageserver.Parameters, format string) ([]string, error) {
	quality, _ := parameters.GetString("gm.quality")

	if len(quality) == 0 && len(arguments) == 1 {
		return arguments, nil
	}

	if len(quality) == 0 && processor.DefaultQualities != nil {
		if q, ok := processor.DefaultQualities[format]; ok {
			quality = q
		}
	}

	if len(quality) > 0 {
		qualityInt, err := strconv.Atoi(quality)
		if err != nil {
			return nil, imageserver.NewError("Invalid quality parameter (parse int error)")
		}

		if qualityInt < 0 {
			return nil, imageserver.NewError("Invalid quality parameter (less than 0)")
		}

		if format == "jpeg" {
			if qualityInt < 0 || qualityInt > 100 {
				return nil, imageserver.NewError("Invalid quality parameter (must be between 0 and 100)")
			}
		}

		arguments = append(arguments, "-quality", quality)
	}

	return arguments, nil
}
