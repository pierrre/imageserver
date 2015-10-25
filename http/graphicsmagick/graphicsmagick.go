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
	if err := imageserver_http.ParseQueryInt("width", req, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryInt("height", req, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("fill", req, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("ignore_ratio", req, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("only_shrink_larger", req, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("only_enlarge_smaller", req, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("background", req, params)
	if err := imageserver_http.ParseQueryBool("extent", req, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("format", req, params)
	if err := imageserver_http.ParseQueryInt("quality", req, params); err != nil {
		return err
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
