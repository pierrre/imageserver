// Package crop provides a imageserver/http.Parser implementation for imageserver/image/crop.Processor.
package crop

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
)

const param = "crop"

// Parser is a imageserver/http.Parser implementation for imageserver/image/crop.Processor.
//
// It uses the "crop" param in the query string, with the following format: min_x,min_y|max_x,max_y
type Parser struct{}

// Parse implements imageserver/http.Parser.
func (prs *Parser) Parse(req *http.Request, params imageserver.Params) error {
	crop := req.URL.Query().Get(param)
	if crop == "" {
		return nil
	}
	var minX, minY, maxX, maxY int
	_, err := fmt.Sscanf(crop, "%d,%d|%d,%d", &minX, &minY, &maxX, &maxY)
	if err != nil {
		return &imageserver.ParamError{
			Param:   param,
			Message: fmt.Sprintf("expected format '<int>,<int>|<int>,<int>': %s", err),
		}
	}
	params.Set(param, imageserver.Params{
		"min_x": minX,
		"min_y": minY,
		"max_x": maxX,
		"max_y": maxY,
	})
	return nil
}

// Resolve implements imageserver/http.Parser.
func (prs *Parser) Resolve(p string) (httpParam string) {
	if p == param || strings.HasPrefix(p, param+".") {
		return param
	}
	return ""
}
