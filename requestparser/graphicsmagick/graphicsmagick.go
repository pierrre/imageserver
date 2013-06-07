package graphicsmagick

import (
	"fmt"
	"github.com/pierrre/imageproxy"
	"net/http"
	"net/url"
	"strconv"
)

type GraphicsMagickRequestParser struct {
}

func (parser *GraphicsMagickRequestParser) ParseRequest(request *http.Request) (parameters imageproxy.Parameters, err error) {
	parameters = make(imageproxy.Parameters)

	query := request.URL.Query()

	err = parser.parseDimension(query, parameters, "width")
	if err != nil {
		return
	}

	err = parser.parseDimension(query, parameters, "height")
	if err != nil {
		return
	}

	err = parser.parseString(query, parameters, "format")
	if err != nil {
		return
	}

	err = parser.parseString(query, parameters, "quality")
	if err != nil {
		return
	}

	err = parser.parseBool(query, parameters, "fill")
	if err != nil {
		return
	}

	err = parser.parseBool(query, parameters, "ignore_ratio")
	if err != nil {
		return
	}

	err = parser.parseBool(query, parameters, "only_shrink_larger")
	if err != nil {
		return
	}

	err = parser.parseBool(query, parameters, "only_enlarge_smaller")
	if err != nil {
		return
	}

	err = parser.parseBool(query, parameters, "extent")
	if err != nil {
		return
	}

	err = parser.parseString(query, parameters, "background")
	if err != nil {
		return
	}

	return
}

func (parser *GraphicsMagickRequestParser) parseDimension(query url.Values, parameters imageproxy.Parameters, dimensionName string) error {
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

func (parser *GraphicsMagickRequestParser) parseString(query url.Values, parameters imageproxy.Parameters, parameterName string) error {
	parameter := query.Get(parameterName)
	if len(parameter) > 0 {
		parser.setParameter(parameters, parameterName, parameter)
	}
	return nil
}

func (parser *GraphicsMagickRequestParser) parseBool(query url.Values, parameters imageproxy.Parameters, parameterName string) error {
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

func (parser *GraphicsMagickRequestParser) setParameter(parameters imageproxy.Parameters, key string, value interface{}) {
	parameters.Set("gm."+key, value)
}
