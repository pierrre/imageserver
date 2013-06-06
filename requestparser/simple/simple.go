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

	err = parseSource(parameters, query)
	if err != nil {
		return
	}

	err = parseWidth(parameters, query)
	if err != nil {
		return
	}

	err = parseHeight(parameters, query)
	if err != nil {
		return
	}

	return
}

func parseSource(parameters imageproxy.Parameters, query url.Values) error {
	source := query.Get("source")
	if len(source) == 0 {
		return fmt.Errorf("Source parameter is missing")
	}
	parameters.Set("source", source)
	return nil
}

func parseWidth(parameters imageproxy.Parameters, query url.Values) error {
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

func parseHeight(parameters imageproxy.Parameters, query url.Values) error {
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
