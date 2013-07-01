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

func (parser *GraphicsMagickParser) Parse(request *http.Request, parameters imageserver.Parameters) (err error) {
	query := request.URL.Query()
	if err = parser.parseDimension(query, parameters, "width"); err != nil {
		return
	}
	if err = parser.parseDimension(query, parameters, "height"); err != nil {
		return
	}
	if err = parser.parseBool(query, parameters, "fill"); err != nil {
		return
	}
	if err = parser.parseBool(query, parameters, "ignore_ratio"); err != nil {
		return
	}
	if err = parser.parseBool(query, parameters, "only_shrink_larger"); err != nil {
		return
	}
	if err = parser.parseBool(query, parameters, "only_enlarge_smaller"); err != nil {
		return
	}
	if err = parser.parseString(query, parameters, "background"); err != nil {
		return
	}
	if err = parser.parseBool(query, parameters, "extent"); err != nil {
		return
	}
	if err = parser.parseString(query, parameters, "format"); err != nil {
		return
	}
	if err = parser.parseString(query, parameters, "quality"); err != nil {
		return
	}
	return
}

func (parser *GraphicsMagickParser) parseDimension(query url.Values, parameters imageserver.Parameters, parameterName string) (err error) {
	dimensionString := query.Get(parameterName)
	if len(dimensionString) == 0 {
		return
	}
	dimension, err := strconv.Atoi(dimensionString)
	if err != nil {
		err = parser.createParseError(parameterName, "int")
		return
	}
	if dimension <= 0 {
		err = parser.createError(parameterName, "lower than or equal to zero")
		return
	}
	parser.setParameter(parameters, parameterName, dimension)
	return
}

func (parser *GraphicsMagickParser) parseString(query url.Values, parameters imageserver.Parameters, parameterName string) (err error) {
	parameter := query.Get(parameterName)
	if len(parameter) == 0 {
		return
	}
	parser.setParameter(parameters, parameterName, parameter)
	return
}

func (parser *GraphicsMagickParser) parseBool(query url.Values, parameters imageserver.Parameters, parameterName string) (err error) {
	parameterString := query.Get(parameterName)
	if len(parameterString) == 0 {
		return
	}
	parameter, err := strconv.ParseBool(parameterString)
	if err != nil {
		err = parser.createParseError(parameterName, "bool")
		return
	}
	parser.setParameter(parameters, parameterName, parameter)
	return
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
