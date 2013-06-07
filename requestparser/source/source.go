package source

import (
	"github.com/pierrre/imageproxy"
	"net/http"
)

type SourceRequestParser struct {
}

func (parser *SourceRequestParser) ParseRequest(request *http.Request) (parameters imageproxy.Parameters, err error) {
	parameters = make(imageproxy.Parameters)

	query := request.URL.Query()

	source := query.Get("source")
	if len(source) > 0 {
		parameters.Set("source", source)
	}

	return
}
