// Package nfntresize provides a nfnt resize HTTP Parser.
package nfntresize

import (
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

const (
	globalParam = "nfntresize"
)

// Parser is a nfnt resize HTTP Parser.
//
// See Processor for params list.
type Parser struct{}

// Parse implements Parser.
func (parser *Parser) Parse(request *http.Request, params imageserver.Params) error {
	p := make(imageserver.Params)
	err := parser.parse(request, p)
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

func (parser *Parser) parse(request *http.Request, params imageserver.Params) error {
	if err := imageserver_http.ParseQueryInt("width", request, params); err != nil {
		return err
	}
	if err := imageserver_http.ParseQueryInt("height", request, params); err != nil {
		return err
	}
	imageserver_http.ParseQueryString("interpolation", request, params)
	imageserver_http.ParseQueryString("mode", request, params)
	return nil
}

// Resolve implements Parser.
func (parser *Parser) Resolve(param string) string {
	if !strings.HasPrefix(param, globalParam+".") {
		return ""
	}
	return strings.TrimPrefix(param, globalParam+".")
}
