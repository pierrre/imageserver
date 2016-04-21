package gift

import (
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

const (
	resizeParam = "gift_resize"
)

// ResizeParser is a imageserver/http.Parser implementation for imageserver/image/gift.ResizeProcessor.
//
// It takes the params from the HTTP URL query and stores them in a Params.
// This Params is added to the given Params at the key "gift_resize".
//
// See imageserver/image/gift.ResizeProcessor for params list.
type ResizeParser struct{}

// Parse implements imageserver/http.Parser.
func (prs *ResizeParser) Parse(req *http.Request, params imageserver.Params) error {
	p := imageserver.Params{}
	err := prs.parse(req, p)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = resizeParam + "." + err.Param
		}
		return err
	}
	if !p.Empty() {
		params.Set(resizeParam, p)
	}
	return nil
}

func (prs *ResizeParser) parse(req *http.Request, params imageserver.Params) error {
	if err := imageserver_http.ParseQueryInt("width", req, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryInt("height", req, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("resampling", req, params)
	imageserver_http.ParseQueryString("mode", req, params)
	return nil
}

// Resolve implements imageserver/http.Parser.
func (prs *ResizeParser) Resolve(param string) string {
	if !strings.HasPrefix(param, resizeParam+".") {
		return ""
	}
	return strings.TrimPrefix(param, resizeParam+".")
}
