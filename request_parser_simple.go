package imageproxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SimpleRequestParser struct {
}

func (parser *SimpleRequestParser) ParseRequest(request *http.Request) (source *url.URL, parameters *Parameters, err error) {
	if request.Method != "GET" {
		err = errors.New("Invalid request method")
		return
	}

	query := request.URL.Query()

	sourceString := query.Get("source")
	if len(sourceString) == 0 {
		err = errors.New("Missing source parameter")
		return
	}
	source, err = url.ParseRequestURI(sourceString)
	if err != nil {
		err = fmt.Errorf("Invalid source parameter (%s)", err)
		return
	}

	parameters = &Parameters{}

	widthString := query.Get("width")
	if len(widthString) > 0 {
		parameters.Width, err = strconv.Atoi(widthString)
		if err != nil {
			return
		}
	}

	heightString := query.Get("height")
	if len(heightString) > 0 {
		parameters.Height, err = strconv.Atoi(heightString)
		if err != nil {
			return
		}
	}

	return
}
