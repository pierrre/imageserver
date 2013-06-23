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

func (parser *GraphicsMagickParser) parseDimension(query url.Values, parameters imageserver.Parameters, dimensionName string) error {
	dimensionString := query.Get(dimensionName)
	if len(dimensionString) > 0 {
		dimension, err := strconv.Atoi(dimensionString)
		if err != nil {
			return err
		}
		if dimension <= 0 {
			return fmt.Errorf("Invalid %s", dimensionName)
		}
		parser.setParameter(parameters, dimensionName, dimension)
	}
	return nil
}

func (parser *GraphicsMagickParser) parseString(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	parameter := query.Get(parameterName)
	if len(parameter) > 0 {
		parser.setParameter(parameters, parameterName, parameter)
	}
	return nil
}

func (parser *GraphicsMagickParser) parseBool(query url.Values, parameters imageserver.Parameters, parameterName string) error {
	parameterString := query.Get(parameterName)
	if len(parameterString) > 0 {
		var parameter bool
		switch parameterString {
		case "0":
			parameter = false
		case "1":
			parameter = true
		default:
			return fmt.Errorf("Invalid %s", parameterName)
		}
		parser.setParameter(parameters, parameterName, parameter)
	}
	return nil
}

func (parser *GraphicsMagickParser) setParameter(parameters imageserver.Parameters, key string, value interface{}) {
	parameters.Set("gm."+key, value)
}
