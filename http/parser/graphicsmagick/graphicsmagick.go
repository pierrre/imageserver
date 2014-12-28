// Package graphicsmagick provides a GraphicsMagick HTTP Parser
package graphicsmagick

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pierrre/imageserver"
)

// Parser represents a GraphicsMagick HTTP Parser
type Parser struct{}

// Parse parses an http Request for GraphicsMagickProcessor
//
// See Processor for params list.
func (parser *Parser) Parse(request *http.Request, params imageserver.Params) error {
	p := make(imageserver.Params)
	params.Set("graphicsmagick", p)
	params = p

	query := request.URL.Query()
	if err := parser.parseInt(query, params, "width"); err != nil {
		return err
	}
	if err := parser.parseInt(query, params, "height"); err != nil {
		return err
	}
	if err := parser.parseBool(query, params, "fill"); err != nil {
		return err
	}
	if err := parser.parseBool(query, params, "ignore_ratio"); err != nil {
		return err
	}
	if err := parser.parseBool(query, params, "only_shrink_larger"); err != nil {
		return err
	}
	if err := parser.parseBool(query, params, "only_enlarge_smaller"); err != nil {
		return err
	}
	parser.parseString(query, params, "background")
	if err := parser.parseBool(query, params, "extent"); err != nil {
		return err
	}
	parser.parseString(query, params, "format")
	if err := parser.parseInt(query, params, "quality"); err != nil {
		return err
	}
	return nil
}

func (parser *Parser) parseString(query url.Values, params imageserver.Params, paramName string) {
	param := query.Get(paramName)
	if param == "" {
		return
	}
	params[paramName] = param
	return
}

func (parser *Parser) parseInt(query url.Values, params imageserver.Params, paramName string) error {
	paramString := query.Get(paramName)
	if paramString == "" {
		return nil
	}
	param, err := strconv.Atoi(paramString)
	if err != nil {
		return parser.newParseError(paramName, "int")
	}
	params[paramName] = param
	return nil
}

func (parser *Parser) parseBool(query url.Values, params imageserver.Params, paramName string) error {
	paramString := query.Get(paramName)
	if paramString == "" {
		return nil
	}
	param, err := strconv.ParseBool(paramString)
	if err != nil {
		return parser.newParseError(paramName, "bool")
	}
	params[paramName] = param
	return nil
}

func (parser *Parser) newParseError(paramName string, parseType string) *imageserver.ParamError {
	return &imageserver.ParamError{
		Param:   fmt.Sprintf("graphicsmagick.%s", paramName),
		Message: fmt.Sprintf("parse %s error", parseType),
	}
}

// Resolve resolves GraphicsMagick's params
func (parser *Parser) Resolve(param string) string {
	if !strings.HasPrefix(param, "graphicsmagick.") {
		return ""
	}
	return strings.TrimPrefix(param, "graphicsmagick.")
}
