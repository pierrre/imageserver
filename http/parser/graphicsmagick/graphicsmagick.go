// Package graphicsmagick provides a GraphicsMagick http Parser
package graphicsmagick

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pierrre/imageserver"
)

// GraphicsMagickParser represents a GraphicsMagick http Parser
type GraphicsMagickParser struct{}

// Parse parses an http Request for GraphicsMagickProcessor
//
// See GraphicsMagickProcessor for parameters list.
func (parser *GraphicsMagickParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	p := make(imageserver.Parameters)
	parameters.Set("graphicsmagick", p)
	parameters = p

	query := request.URL.Query()
	if err := parser.parseInt(query, parameters, "width"); err != nil {
		return err
	}
	if err := parser.parseInt(query, parameters, "height"); err != nil {
		return err
	}
	if err := parser.parseBool(query, parameters, "fill"); err != nil {
		return err
	}
	if err := parser.parseBool(query, parameters, "ignore_ratio"); err != nil {
		return err
	}
	if err := parser.parseBool(query, parameters, "only_shrink_larger"); err != nil {
		return err
	}
	if err := parser.parseBool(query, parameters, "only_enlarge_smaller"); err != nil {
		return err
	}
	parser.parseString(query, parameters, "background")
	if err := parser.parseBool(query, parameters, "extent"); err != nil {
		return err
	}
	parser.parseString(query, parameters, "format")
	if err := parser.parseInt(query, parameters, "quality"); err != nil {
		return err
	}
	return nil
}

func (parser *GraphicsMagickParser) parseString(query url.Values, parameters imageserver.Parameters, parameterName string) {
	parameter := query.Get(parameterName)
	if len(parameter) == 0 {
		return
	}
	parameters[parameterName] = parameter
	return
}

func (parser *GraphicsMagickParser) parseInt(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	parameterString := query.Get(parameterName)
	if len(parameterString) == 0 {
		return nil
	}
	parameter, err := strconv.Atoi(parameterString)
	if err != nil {
		return parser.createParseError(parameterName, "int")
	}
	parameters[parameterName] = parameter
	return nil
}

func (parser *GraphicsMagickParser) parseBool(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	parameterString := query.Get(parameterName)
	if len(parameterString) == 0 {
		return nil
	}
	parameter, err := strconv.ParseBool(parameterString)
	if err != nil {
		return parser.createParseError(parameterName, "bool")
	}
	parameters[parameterName] = parameter
	return nil
}

func (parser *GraphicsMagickParser) createError(parameterName string, cause string) *imageserver.Error {
	return imageserver.NewError(fmt.Sprintf("invalid \"%s\" parameter: %s", parameterName, cause))
}

func (parser *GraphicsMagickParser) createParseError(parameterName string, parseType string) *imageserver.Error {
	return parser.createError(parameterName, fmt.Sprintf("parse %s error", parseType))
}
