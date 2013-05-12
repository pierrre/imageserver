package imageproxy

import (
	"errors"
)

type Parameters struct {
	Width  int
	Height int
}

func (parameters *Parameters) Validate() error {
	if parameters.Width < 0 {
		return errors.New("Invalid width")
	}

	if parameters.Height < 0 {
		return errors.New("Invalid height parameter")
	}

	return nil
}
