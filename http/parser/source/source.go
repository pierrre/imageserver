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
	if source != "" {
		parameters.Set("source", source)
	}
	return nil
}

// Resolve resolves the "source" parameter
func (parser *SourceParser) Resolve(parameter string) string {
	if parameter != "source" {
		return ""
	}
	return "source"
}
