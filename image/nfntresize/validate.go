package nfntresize

import (
	"fmt"

	"github.com/pierrre/imageserver"
)

// ValidateParamsServer is an Image Server that validates Params for the Processor.
type ValidateParamsServer struct {
	imageserver.Server
	WidthMax  uint
	HeightMax uint
}

// Get implements Server.
func (srv *ValidateParamsServer) Get(params imageserver.Params) (*imageserver.Image, error) {
	if params.Has(Param) {
		params, err := params.GetParams(Param)
		if err != nil {
			return nil, err
		}
		err = srv.validate(params)
		if err != nil {
			return nil, wrapParamError(err)
		}
	}
	return srv.Server.Get(params)
}

func (srv *ValidateParamsServer) validate(params imageserver.Params) error {
	if params.Empty() {
		return nil
	}
	err := validateDimension("width", srv.WidthMax, params)
	if err != nil {
		return err
	}
	err = validateDimension("height", srv.HeightMax, params)
	if err != nil {
		return err
	}
	return nil
}

func validateDimension(name string, max uint, params imageserver.Params) error {
	if max == 0 {
		return nil
	}
	d, err := getDimension(name, params)
	if err != nil {
		return err
	}
	if d > max {
		return &imageserver.ParamError{Param: name, Message: fmt.Sprintf("must be less than or equal to %d", max)}
	}
	return nil
}
