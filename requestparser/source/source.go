package source

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

type SourceRequestParser struct {
}

func (parser *SourceRequestParser) ParseRequest(request *http.Request) (parameters imageserver.Parameters, err error) {
	parameters = make(imageserver.Parameters)

	query := request.URL.Query()

	source := query.Get("source")
	if len(source) > 0 {
		parameters.Set("source", source)
	}

	return
}
