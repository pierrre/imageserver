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
// See Processor for parameters list.
func (parser *Parser) Parse(request *http.Request, parameters imageserver.Parameters) error {
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

func (parser *Parser) parseString(query url.Values, parameters imageserver.Parameters, parameterName string) {
	parameter := query.Get(parameterName)
	if parameter == "" {
		return
	}
	parameters[parameterName] = parameter
	return
}

func (parser *Parser) parseInt(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	parameterString := query.Get(parameterName)
	if parameterString == "" {
		return nil
	}
	parameter, err := strconv.Atoi(parameterString)
	if err != nil {
		return parser.newParseError(parameterName, "int")
	}
	parameters[parameterName] = parameter
	return nil
}

func (parser *Parser) parseBool(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	parameterString := query.Get(parameterName)
	if parameterString == "" {
		return nil
	}
	parameter, err := strconv.ParseBool(parameterString)
	if err != nil {
		return parser.newParseError(parameterName, "bool")
	}
	parameters[parameterName] = parameter
	return nil
}

func (parser *Parser) newParseError(parameterName string, parseType string) *imageserver.ParameterError {
	return &imageserver.ParameterError{
		Parameter: fmt.Sprintf("graphicsmagick.%s", parameterName),
		Message:   fmt.Sprintf("parse %s error", parseType),
	}
}

// Resolve resolves GraphicsMagick's parameters
func (parser *Parser) Resolve(parameter string) string {
	if !strings.HasPrefix(parameter, "graphicsmagick.") {
		return ""
	}
	return strings.TrimPrefix(parameter, "graphicsmagick.")
}
