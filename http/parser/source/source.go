// Package source provides an http Parser that takes the "source" parameter from query
package source

import (
	"net/http"

	"github.com/pierrre/imageserver"
)

// SourceParser represents an http Parser that takes the "source" parameter from query
type SourceParser struct{}

// Parse takes the "source" parameter from query
func (parser *SourceParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	query := request.URL.Query()
	source := query.Get("source")
	if len(source) > 0 {
		parameters.Set("source", source)
	}
	return nil
}
