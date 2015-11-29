// Package graphicsmagick provides a GraphicsMagick HTTP Parser.
package graphicsmagick

import (
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

const (
	globalParam = "graphicsmagick"
)

// Parser is a GraphicsMagick HTTP Parser.
//
// See Server for params list.
type Parser struct{}

// Parse implements Parser.
func (parser *Parser) Parse(req *http.Request, params imageserver.Params) error {
	p := imageserver.Params{}
	err := parser.parse(req, p)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = globalParam + "." + err.Param
		}
		return err
	}
	if !p.Empty() {
		params.Set(globalParam, p)
	}
	return nil
}

func (parser *Parser) parse(req *http.Request, params imageserver.Params) error {
	// parse bool params
	boolParams := []string{
		"ignore_ratio",
		"fill",
		"only_shrink_larger",
		"only_enlarge_smaller",
		"extent",
		"monochrome",
		"grey",
		"no_strip",
		"trim",
		"no_interlace",
		"flip",
		"flop",
	}
	for _, bp := range boolParams {
		err := imageserver_http.ParseQueryBool(bp, req, params)
		if err != nil {
			return err
		}
	}

	// parse integer params
	intParams := []string{
		"w",
		"h",
		"rotate",
		"q",
	}
	for _, ip := range intParams {
		err := imageserver_http.ParseQueryInt(ip, req, params)
		if err != nil {
			return err
		}
	}

	// parse string params
	stringParams := []string{
		"bg",
		"gravity",
		"crop",
		"format",
	}
	for _, sp := range stringParams {
		imageserver_http.ParseQueryString(sp, req, params)
	}
	return nil
}

// Resolve implements Parser.
func (parser *Parser) Resolve(param string) string {
	if !strings.HasPrefix(param, globalParam+".") {
		return ""
	}
	return strings.TrimPrefix(param, globalParam+".")
}
