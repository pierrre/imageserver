package gift

import (
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

const (
	rotateParam = "gift_rotate"
)

// RotateParser is a imageserver/http.Parser implementation for imageserver/image/gift.RotateProcessor.
//
// It takes the params from the HTTP URL query and stores them in a Params.
// This Params is added to the given Params at the key "gift_rotate".
//
// See imageserver/image/gift.RotateProcessor for params list.
type RotateParser struct{}

// Parse implements imageserver/http.Parser.
func (prs *RotateParser) Parse(req *http.Request, params imageserver.Params) error {
	p := imageserver.Params{}
	err := prs.parse(req, p)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = rotateParam + "." + err.Param
		}
		return err
	}
	if !p.Empty() {
		params.Set(rotateParam, p)
	}
	return nil
}

func (prs *RotateParser) parse(req *http.Request, params imageserver.Params) error {
	if err := imageserver_http.ParseQueryFloat("rotation", req, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("background", req, params)
	imageserver_http.ParseQueryString("interpolation", req, params)
	return nil
}

// Resolve implements imageserver/http.Parser.
func (prs *RotateParser) Resolve(param string) string {
	if !strings.HasPrefix(param, rotateParam+".") {
		return ""
	}
	return strings.TrimPrefix(param, rotateParam+".")
}
