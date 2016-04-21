// Package graphicsmagick provides a GraphicsMagick imageserver.Handler implementation.
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
	param         = "graphicsmagick"
	tempDirPrefix = "imageserver_"
)

// Handler is a GraphicsMagick imageserver.Handler implementation.
//
// It processes the Image with the GraphicsMagick command line (mogrify command).
//
// All params are extracted from the "graphicsmagick" node param and are optionals.
//
// Params (see GraphicsMagick documentation for more information about arguments):
//  - width / height: sizes for "-resize" argument (both optionals)
//  - fill: "^" for "-resize" argument
//  - ignore_ratio: "!" for "-resize" argument
//  - only_shrink_larger: ">" for "-resize" argument
//  - only_enlarge_smaller: "<" for "-resize" argument
//  - background: color for "-background" argument, 3/4/6/8 lower case hexadecimal characters
//  - extent: "-extent" param, uses width/height params and add "-gravity center" argument
//  - format: "-format" param
//  - quality: "-quality" param
type Handler struct {
	// Executable is the path to "gm" executable, usually "/usr/bin/gm".
	Executable string

	// Timeoput is an optional timeout for process.
	Timeout time.Duration

	// TempDir is an optional temp directory for image files.
	TempDir string

	// AllowedFormats is an optional list of allowed formats.
	AllowedFormats []string
}

// Handle implements imageserver.Handler.
func (hdr *Handler) Handle(im *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	if !params.Has(param) {
		return im, nil
	}
	params, err := params.GetParams(param)
	if err != nil {
		return nil, err
	}
	if params.Empty() {
		return im, nil
	}
	im, err = hdr.handle(im, params)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = param + "." + err.Param
		}
		return nil, err
	}
	return im, nil
}

func (hdr *Handler) handle(im *imageserver.Image, params imageserver.Params) (*imageserver.Image, error) {
	arguments := list.New()

	width, height, err := hdr.buildArgumentsResize(arguments, params)
	if err != nil {
		return nil, err
	}

	err = hdr.buildArgumentsBackground(arguments, params)
	if err != nil {
		return nil, err
	}

	err = hdr.buildArgumentsExtent(arguments, params, width, height)
	if err != nil {
		return nil, err
	}

	format, formatSpecified, err := hdr.buildArgumentsFormat(arguments, params, im)
	if err != nil {
		return nil, err
	}

	err = hdr.buildArgumentsQuality(arguments, params, format)
	if err != nil {
		return nil, err
	}

	if arguments.Len() == 0 {
		return im, nil
	}

	arguments.PushFront("mogrify")

	tempDir, err := ioutil.TempDir(hdr.TempDir, tempDirPrefix)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	file := filepath.Join(tempDir, "image")
	arguments.PushBack(file)
	err = ioutil.WriteFile(file, im.Data, os.FileMode(0600))
	if err != nil {
		return nil, err
	}

	argumentSlice := convertArgumentsToSlice(arguments)
	cmd := exec.Command(hdr.Executable, argumentSlice...)
	err = hdr.runCommand(cmd)
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

	im = &imageserver.Image{
		Format: format,
		Data:   data,
	}
	return im, nil
}

func (hdr *Handler) buildArgumentsResize(arguments *list.List, params imageserver.Params) (width int, height int, err error) {
	width, err = getDimension("width", params)
	if err != nil {
		return 0, 0, err
	}
	height, err = getDimension("height", params)
	if err != nil {
		return 0, 0, err
	}
	if width == 0 && height == 0 {
		return 0, 0, nil
	}
	widthString := ""
	if width != 0 {
		widthString = strconv.Itoa(width)
	}
	heightString := ""
	if height != 0 {
		heightString = strconv.Itoa(height)
	}
	resize := fmt.Sprintf("%sx%s", widthString, heightString)
	if params.Has("fill") {
		fill, err := params.GetBool("fill")
		if err != nil {
			return 0, 0, err
		}
		if fill {
			resize = resize + "^"
		}
	}
	if params.Has("ignore_ratio") {
		ignoreRatio, err := params.GetBool("ignore_ratio")
		if err != nil {
			return 0, 0, err
		}
		if ignoreRatio {
			resize = resize + "!"
		}
	}
	if params.Has("only_shrink_larger") {
		onlyShrinkLarger, err := params.GetBool("only_shrink_larger")
		if err != nil {
			return 0, 0, err
		}
		if onlyShrinkLarger {
			resize = resize + ">"
		}
	}
	if params.Has("only_enlarge_smaller") {
		onlyEnlargeSmaller, err := params.GetBool("only_enlarge_smaller")
		if err != nil {
			return 0, 0, err
		}
		if onlyEnlargeSmaller {
			resize = resize + "<"
		}
	}
	arguments.PushBack("-resize")
	arguments.PushBack(resize)
	return width, height, nil
}

func getDimension(name string, params imageserver.Params) (int, error) {
	if !params.Has(name) {
		return 0, nil
	}
	dimension, err := params.GetInt(name)
	if err != nil {
		return 0, err
	}
	if dimension < 0 {
		return 0, &imageserver.ParamError{Param: name, Message: "must be greater than or equal to 0"}
	}
	return dimension, nil
}

func (hdr *Handler) buildArgumentsBackground(arguments *list.List, params imageserver.Params) error {
	if !params.Has("background") {
		return nil
	}
	background, err := params.GetString("background")
	if err != nil {
		return err
	}
	switch len(background) {
	case 3, 4, 6, 8:
	default:
		return &imageserver.ParamError{Param: "background", Message: "length must be equal to 3, 4, 6 or 8"}
	}
	for _, r := range background {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return &imageserver.ParamError{Param: "background", Message: "must only contain characters in 0-9a-f"}
		}
	}
	arguments.PushBack("-background")
	arguments.PushBack(fmt.Sprintf("#%s", background))
	return nil
}

func (hdr *Handler) buildArgumentsExtent(arguments *list.List, params imageserver.Params, width int, height int) error {
	if width == 0 || height == 0 {
		return nil
	}
	if !params.Has("extent") {
		return nil
	}
	extent, err := params.GetBool("extent")
	if err != nil {
		return err
	}
	if extent {
		arguments.PushBack("-gravity")
		arguments.PushBack("center")
		arguments.PushBack("-extent")
		arguments.PushBack(fmt.Sprintf("%dx%d", width, height))
	}
	return nil
}

func (hdr *Handler) buildArgumentsFormat(arguments *list.List, params imageserver.Params, sourceImage *imageserver.Image) (format string, formatSpecified bool, err error) {
	if !params.Has("format") {
		return sourceImage.Format, false, nil
	}
	format, err = params.GetString("format")
	if err != nil {
		return "", false, err
	}
	if hdr.AllowedFormats != nil {
		ok := false
		for _, f := range hdr.AllowedFormats {
			if f == format {
				ok = true
				break
			}
		}
		if !ok {
			return "", false, &imageserver.ParamError{Param: "format", Message: "not allowed"}
		}
	}
	arguments.PushBack("-format")
	arguments.PushBack(format)
	return format, true, nil
}

func (hdr *Handler) buildArgumentsQuality(arguments *list.List, params imageserver.Params, format string) error {
	if !params.Has("quality") {
		return nil
	}
	quality, err := params.GetInt("quality")
	if err != nil {
		return err
	}
	if quality < 0 {
		return &imageserver.ParamError{Param: "quality", Message: "must be greater than or equal to 0"}
	}
	if format == "jpeg" {
		if quality < 0 || quality > 100 {
			return &imageserver.ParamError{Param: "quality", Message: "must be between 0 and 100"}
		}
	}
	arguments.PushBack("-quality")
	arguments.PushBack(strconv.Itoa(quality))
	return nil
}

func convertArgumentsToSlice(arguments *list.List) []string {
	argumentSlice := make([]string, 0, arguments.Len())
	for e := arguments.Front(); e != nil; e = e.Next() {
		argumentSlice = append(argumentSlice, e.Value.(string))
	}
	return argumentSlice
}

func (hdr *Handler) runCommand(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmdChan := make(chan error, 1)
	go func() {
		cmdChan <- cmd.Wait()
	}()
	var timeoutChan <-chan time.Time
	if hdr.Timeout != 0 {
		timeoutChan = time.After(hdr.Timeout)
	}
	select {
	case err = <-cmdChan:
	case <-timeoutChan:
		cmd.Process.Kill()
		err = fmt.Errorf("timeout after %s", hdr.Timeout)
	}
	if err != nil {
		return &imageserver.ImageError{Message: fmt.Sprintf("GraphicsMagick command: %s", err)}
	}
	return nil
}
