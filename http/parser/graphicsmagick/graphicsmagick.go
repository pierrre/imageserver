// Package graphicsmagick provides a GraphicsMagick HTTP Parser.
package graphicsmagick

import (
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

// Parser represents a GraphicsMagick HTTP Parser.
//
// See Server for params list.
type Parser struct{}

// Parse implements Parser.
func (parser *Parser) Parse(request *http.Request, params imageserver.Params) error {
	p := make(imageserver.Params)
	err := parser.parse(request, p)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = "graphicsmagick." + err.Param
		}
		return err
	}
	if !p.Empty() {
		params.Set("graphicsmagick", p)
	}
	return nil
}

func (parser *Parser) parse(request *http.Request, params imageserver.Params) error {
	if err := imageserver_http.ParseQueryInt("width", request, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryInt("height", request, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("fill", request, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("ignore_ratio", request, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("only_shrink_larger", request, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryBool("only_enlarge_smaller", request, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("background", request, params)
	if err := imageserver_http.ParseQueryBool("extent", request, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("format", request, params)
	if err := imageserver_http.ParseQueryInt("quality", request, params); err != nil {
		return err
	}
	return nil
}

// Resolve implements Parser.
func (parser *Parser) Resolve(param string) string {
	if !strings.HasPrefix(param, "graphicsmagick.") {
		return ""
	}
	return strings.TrimPrefix(param, "graphicsmagick.")
}
