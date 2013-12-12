// Package source provides an http Parser that takes the "source" parameter from query
package source

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

// SourceParser represents an http Parser that takes the "source" parameter from query
type SourceParser struct {
}

// Parse parses the http Request and takes the "source" parameter from query
func (parser *SourceParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	query := request.URL.Query()
	source := query.Get("source")
	if len(source) > 0 {
		parameters.Set("source", source)
	}
	return nil
}
