// Package gift provides a GIFT HTTP Parser
package gift

import (
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

const (
	// Param is the sub-param used by this package.
	Param = "gift"
)

// Parser is a GIFT HTTP Parser.
//
// See Processor for params list.
type Parser struct{}

// Parse implements Parser.
func (parser *Parser) Parse(req *http.Request, params imageserver.Params) error {
	p := imageserver.Params{}
	err := parse(req, p)
	if err != nil {
		if err, ok := err.(*imageserver.ParamError); ok {
			err.Param = Param + "." + err.Param
		}
		return err
	}
	if !p.Empty() {
		params.Set(Param, p)
	}
	return nil
}

func parse(req *http.Request, params imageserver.Params) error {
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

// Resolve implements Parser.
func (parser *Parser) Resolve(param string) string {
	if !strings.HasPrefix(param, Param+".") {
		return ""
	}
	return strings.TrimPrefix(param, Param+".")
}
