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

	err = parser.parseWidth(query, parameters)
	if err != nil {
		return
	}

	err = parser.parseHeight(query, parameters)
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

func (parser *GraphicsMagickRequestParser) parseWidth(query url.Values, parameters imageproxy.Parameters) error {
	widthString := query.Get("width")
	if len(widthString) > 0 {
		width, err := strconv.Atoi(widthString)
		if err != nil {
			return err
		}
		if width <= 0 {
			return fmt.Errorf("Invalid width parameter")
		}
		parser.setParameter(parameters, "width", width)
	}
	return nil
}

func (parser *GraphicsMagickRequestParser) parseHeight(query url.Values, parameters imageproxy.Parameters) error {
	heightString := query.Get("height")
	if len(heightString) > 0 {
		height, err := strconv.Atoi(heightString)
		if err != nil {
			return err
		}
		if height <= 0 {
			return fmt.Errorf("Invalid height parameter")
		}
		parser.setParameter(parameters, "height", height)
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
