package simple

import (
	"fmt"
	"github.com/pierrre/imageproxy"
	"net/http"
	"net/url"
	"strconv"
)

type SimpleRequestParser struct {
}

func (parser *SimpleRequestParser) ParseRequest(request *http.Request) (parameters imageproxy.Parameters, err error) {
	if request.Method != "GET" {
		err = fmt.Errorf("Invalid request method")
		return
	}

	parameters = make(imageproxy.Parameters)

	query := request.URL.Query()

	err = parser.parseSource(query, parameters)
	if err != nil {
		return
	}

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

func (parser *SimpleRequestParser) parseSource(query url.Values, parameters imageproxy.Parameters) error {
	source := query.Get("source")
	if len(source) > 0 {
		parameters.Set("source", source)
	}
	return nil
}

func (parser *SimpleRequestParser) parseWidth(query url.Values, parameters imageproxy.Parameters) error {
	widthString := query.Get("width")
	if len(widthString) > 0 {
		width, err := strconv.Atoi(widthString)
		if err != nil {
			return err
		}
		if width <= 0 {
			return fmt.Errorf("Invalid width parameter")
		}
		parameters.Set("width", width)
	}
	return nil
}

func (parser *SimpleRequestParser) parseHeight(query url.Values, parameters imageproxy.Parameters) error {
	heightString := query.Get("height")
	if len(heightString) > 0 {
		height, err := strconv.Atoi(heightString)
		if err != nil {
			return err
		}
		if height <= 0 {
			return fmt.Errorf("Invalid height parameter")
		}
		parameters.Set("height", height)
	}
	return nil
}

func (parser *SimpleRequestParser) parseFormat(query url.Values, parameters imageproxy.Parameters) error {
	format := query.Get("format")
	if len(format) > 0 {
		parameters.Set("format", format)
	}
	return nil
}

func (parser *SimpleRequestParser) parseQuality(query url.Values, parameters imageproxy.Parameters) error {
	quality := query.Get("quality")
	if len(quality) > 0 {
		parameters.Set("quality", quality)
	}
	return nil
}
