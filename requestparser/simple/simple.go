package simple

import (
	"errors"
	"github.com/pierrre/imageproxy"
	"net/http"
)

type SimpleRequestParser struct {
}

func (parser *SimpleRequestParser) ParseRequest(request *http.Request) (parameters *imageproxy.Parameters, err error) {
	if request.Method != "GET" {
		err = errors.New("Invalid request method")
		return
	}

	parameters = &imageproxy.Parameters{}

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
