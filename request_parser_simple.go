package imageproxy

import (
	"errors"
	"net/http"
)

type SimpleRequestParser struct {
}

func (parser *SimpleRequestParser) ParseRequest(request *http.Request) (parameters *Parameters, err error) {
	if request.Method != "GET" {
		err = errors.New("Invalid request method")
		return
	}

	parameters = &Parameters{}

	query := request.URL.Query()

	err = parameters.ParseSource(query.Get("source"))
	if err != nil {
		return
	}

	err = parameters.ParseWidth(query.Get("width"))
	if err != nil {
		return
	}

	err = parameters.ParseHeight(query.Get("height"))
	if err != nil {
		return
	}

	return
}
