// Package graphicsmagick provides a GraphicsMagick Image Processor
package graphicsmagick

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pierrre/imageserver"
)

const (
	globalParameterName = "graphicsmagick"
	tempDirPrefix       = "imageserver_"
)

// Processor represents a GraphicsMagick Image Processor
type Processor struct {
	Executable string // path to "gm" executable, usually "/usr/bin/gm"

	Timeout        time.Duration // timeout for process, optional
	TempDir        string        // temp directory for image files, optional
	AllowedFormats []string      // allowed format list, optional
}

// Process processes Image with the GraphicsMagick command line (mogrify command)
//
// All parameters are extracted from the "graphicsmagick" node parameter and are optionals.
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
func (processor *Processor) Process(image *imageserver.Image, parameters imageserver.Parameters) (*imageserver.Image, error) {
	parameters, err := processor.getParameters(parameters)
	if err != nil {
		return nil, err
	}
	if parameters == nil || parameters.Empty() {
		return image, nil
	}

	arguments := list.New()

	width, height, err := processor.buildArgumentsResize(arguments, parameters)
	if err != nil {
		return nil, err
	}

	err = processor.buildArgumentsBackground(arguments, parameters)
	if err != nil {
		return nil, err
	}

	err = processor.buildArgumentsExtent(arguments, parameters, width, height)
	if err != nil {
		return nil, err
	}

	format, formatSpecified, err := processor.buildArgumentsFormat(arguments, parameters, image)
	if err != nil {
		return nil, err
	}

	err = processor.buildArgumentsQuality(arguments, parameters, format)
	if err != nil {
		return nil, err
	}

	if arguments.Len() == 0 {
		return image, nil
	}

	arguments.PushFront("mogrify")

	tempDir, err := ioutil.TempDir(processor.TempDir, tempDirPrefix)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	file := filepath.Join(tempDir, "image")
	arguments.PushBack(file)
	err = ioutil.WriteFile(file, image.Data, os.FileMode(0600))
	if err != nil {
		return nil, err
	}

	argumentSlice := processor.convertArgumentsToSlice(arguments)

	cmd := exec.Command(processor.Executable, argumentSlice...)

	err = processor.runCommand(cmd)
	if err != nil {
		return nil, err
	}

	if formatSpecified {
		file = fmt.Sprintf("%s.%s", file, format)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	resultImage := &imageserver.Image{
		Format: format,
		Data:   data,
	}

	return resultImage, nil
}

func (processor *Processor) getParameters(parameters imageserver.Parameters) (imageserver.Parameters, error) {
	if !parameters.Has(globalParameterName) {
		return nil, nil
	}

	return parameters.GetParameters(globalParameterName)
}

func (processor *Processor) buildArgumentsResize(arguments *list.List, parameters imageserver.Parameters) (width int, height int, err error) {
	width, _ = parameters.GetInt("width")
	if width < 0 {
		return 0, 0, newParameterError("width", "must be greater than or equal to 0")
	}

	height, _ = parameters.GetInt("height")
	if height < 0 {
		return 0, 0, newParameterError("height", "must be greater than or equal to 0")
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

		if fill, _ := parameters.GetBool("fill"); fill {
			resize = resize + "^"
		}

		if ignoreRatio, _ := parameters.GetBool("ignore_ratio"); ignoreRatio {
			resize = resize + "!"
		}

		if onlyShrinkLarger, _ := parameters.GetBool("only_shrink_larger"); onlyShrinkLarger {
			resize = resize + ">"
		}

		if onlyEnlargeSmaller, _ := parameters.GetBool("only_enlarge_smaller"); onlyEnlargeSmaller {
			resize = resize + "<"
		}

		arguments.PushBack("-resize")
		arguments.PushBack(resize)
	}

	return width, height, nil
}

func (processor *Processor) buildArgumentsBackground(arguments *list.List, parameters imageserver.Parameters) error {
	background, _ := parameters.GetString("background")

	if backgroundLength := len(background); backgroundLength > 0 {
		if backgroundLength != 6 && backgroundLength != 8 && backgroundLength != 3 && backgroundLength != 4 {
			return newParameterError("background", "length must be equal to 3, 4, 6 or 8")
		}

		for _, r := range background {
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
				return newParameterError("background", "must only contain characters in 0-9a-f")
			}
		}

		arguments.PushBack("-background")
		arguments.PushBack(fmt.Sprintf("#%s", background))
	}

	return nil
}

func (processor *Processor) buildArgumentsExtent(arguments *list.List, parameters imageserver.Parameters, width int, height int) error {
	if width != 0 && height != 0 {
		if extent, _ := parameters.GetBool("extent"); extent {
			arguments.PushBack("-gravity")
			arguments.PushBack("center")

			arguments.PushBack("-extent")
			arguments.PushBack(fmt.Sprintf("%dx%d", width, height))
		}
	}

	return nil
}

func (processor *Processor) buildArgumentsFormat(arguments *list.List, parameters imageserver.Parameters, sourceImage *imageserver.Image) (format string, formatSpecified bool, err error) {
	format, _ = parameters.GetString("format")

	formatSpecified = true
	if format == "" {
		format = sourceImage.Format
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
			return "", false, newParameterError("format", "not allowed")
		}
	}

	if formatSpecified {
		arguments.PushBack("-format")
		arguments.PushBack(format)
	}

	return format, formatSpecified, nil
}

func (processor *Processor) buildArgumentsQuality(arguments *list.List, parameters imageserver.Parameters, format string) error {
	if !parameters.Has("quality") {
		return nil
	}

	quality, err := parameters.GetInt("quality")
	if err != nil {
		return err
	}

	if quality < 0 {
		return newParameterError("quality", "must be greater than or equal to 0")
	}

	if format == "jpeg" {
		if quality < 0 || quality > 100 {
			return newParameterError("quality", "must be between 0 and 100")
		}
	}

	arguments.PushBack("-quality")
	arguments.PushBack(strconv.Itoa(quality))

	return nil
}

func (processor *Processor) convertArgumentsToSlice(arguments *list.List) []string {
	argumentSlice := make([]string, 0, arguments.Len())
	for e := arguments.Front(); e != nil; e = e.Next() {
		argumentSlice = append(argumentSlice, e.Value.(string))
	}
	return argumentSlice
}

func (processor *Processor) runCommand(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return err
	}

	cmdChan := make(chan error)
	go func() {
		cmdChan <- cmd.Wait()
	}()

	var timeoutChan <-chan time.Time
	if processor.Timeout != 0 {
		timeoutChan = time.After(processor.Timeout)
	}

	select {
	case err = <-cmdChan:
		if err != nil {
			return err
		}

		return nil
	case <-timeoutChan:
		err = cmd.Process.Kill()
		if err != nil {
			return err
		}

		<-cmdChan

		return fmt.Errorf("command timeout after %s: %+v", processor.Timeout, cmd)
	}
}

func newParameterError(parameter string, message string) *imageserver.ParameterError {
	return &imageserver.ParameterError{
		Parameter: fmt.Sprintf("graphicsmagick.%s", parameter),
		Message:   message,
	}
}
