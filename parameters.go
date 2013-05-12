package imageproxy

import (
	"errors"
	"net/url"
	"strconv"
)

type Parameters struct {
	Source *url.URL
	Width  int
	Height int
}

func (parameters *Parameters) ParseSource(source string) (err error) {
	parameters.Source, err = url.ParseRequestURI(source)
	return
}

func (parameters *Parameters) ParseWidth(width string) (err error) {
	if len(width) != 0 {
		parameters.Width, err = strconv.Atoi(width)
	} else {
		parameters.Width = 0
	}

	return
}

func (parameters *Parameters) ParseHeight(height string) (err error) {
	if len(height) != 0 {
		parameters.Height, err = strconv.Atoi(height)
	} else {
		parameters.Height = 0
	}

	return
}

func (parameters *Parameters) Validate() error {
	if parameters.Source == nil {
		return errors.New("Invalid source")
	}
	if !parameters.Source.IsAbs() {
		return errors.New("Invalid source")
	}

	if parameters.Width < 0 {
		return errors.New("Invalid width")
	}

	if parameters.Height < 0 {
		return errors.New("Invalid height")
	}

	return nil
}
