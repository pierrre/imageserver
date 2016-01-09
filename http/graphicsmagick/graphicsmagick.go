// Package graphicsmagick provides a imageserver/http.Parser implementation for imageserver/graphicsmagick.Handler.
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

// Parser is a imageserver/http.Parser implementation for imageserver/graphicsmagick.Handler.
//
// It takes the params from the HTTP URL query and stores them in a Params.
// This Params is added to the given Params at the key "graphicsmagick".
//
// See imageserver/graphicsmagick.Handler for params list.
type Parser struct{}

// Parse implements imageserver/http.Parser.
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

// Resolve implements imageserver/http.Parser.
func (parser *Parser) Resolve(param string) string {
	if !strings.HasPrefix(param, globalParam+".") {
		return ""
	}
	return strings.TrimPrefix(param, globalParam+".")
}
