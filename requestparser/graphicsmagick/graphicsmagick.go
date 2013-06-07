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
	/*
		TODO
		fill
		ignore_ratio
		only_shrink_larger
		only_enlarge_smaller
		extent
		gravity
	*/
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

	err = parser.parseFormat(query, parameters)
	if err != nil {
		return
	}

	err = parser.parseQuality(query, parameters)
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

func (parser *GraphicsMagickRequestParser) parseFormat(query url.Values, parameters imageproxy.Parameters) error {
	format := query.Get("format")
	if len(format) > 0 {
		parser.setParameter(parameters, "format", format)
	}
	return nil
}

func (parser *GraphicsMagickRequestParser) parseQuality(query url.Values, parameters imageproxy.Parameters) error {
	quality := query.Get("quality")
	if len(quality) > 0 {
		parser.setParameter(parameters, "quality", quality)
	}
	return nil
}

func (parser *GraphicsMagickRequestParser) setParameter(parameters imageproxy.Parameters, key string, value interface{}) {
	parameters.Set("gm."+key, value)
}
