package graphicsmagick

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"net/http"
	"net/url"
	"strconv"
)

type GraphicsMagickParser struct {
}

func (parser *GraphicsMagickParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	query := request.URL.Query()
	if err := parser.parseDimension(query, parameters, "width"); err != nil {
		return err
	}
	if err := parser.parseDimension(query, parameters, "height"); err != nil {
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
	if err := parser.parseString(query, parameters, "background"); err != nil {
		return err
	}
	if err := parser.parseBool(query, parameters, "extent"); err != nil {
		return err
	}
	if err := parser.parseString(query, parameters, "format"); err != nil {
		return err
	}
	if err := parser.parseString(query, parameters, "quality"); err != nil {
		return err
	}
	return nil
}

func (parser *GraphicsMagickParser) parseDimension(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	dimensionString := query.Get(parameterName)
	if len(dimensionString) == 0 {
		return nil
	}
	dimension, err := strconv.Atoi(dimensionString)
	if err != nil {
		return parser.createParseError(parameterName, "int")
	}
	if dimension <= 0 {
		return parser.createError(parameterName, "lower than or equal to zero")
	}
	parser.setParameter(parameters, parameterName, dimension)
	return nil
}

func (parser *GraphicsMagickParser) parseString(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	parameter := query.Get(parameterName)
	if len(parameter) == 0 {
		return nil
	}
	parser.setParameter(parameters, parameterName, parameter)
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
	parser.setParameter(parameters, parameterName, parameter)
	return nil
}

func (parser *GraphicsMagickParser) setParameter(parameters imageserver.Parameters, key string, value interface{}) {
	parameters.Set("gm."+key, value)
}

func (parser *GraphicsMagickParser) createError(parameterName string, cause string) *imageserver.Error {
	return imageserver.NewError(fmt.Sprintf("Invalid %s parameter (%s)", parameterName, cause))
}

func (parser *GraphicsMagickParser) createParseError(parameterName string, parseType string) *imageserver.Error {
	return parser.createError(parameterName, fmt.Sprintf("parse %s error", parseType))
}
