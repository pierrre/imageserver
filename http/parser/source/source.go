// Package source provides an HTTP Parser that takes the "source" parameter from query
package source

import (
	"net/http"

	"github.com/pierrre/imageserver"
)

// Parser represents an http Parser that takes the "source" parameter from query
type Parser struct{}

// Parse takes the "source" parameter from query
func (parser *Parser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	query := request.URL.Query()
	source := query.Get("source")
	if source != "" {
		parameters.Set("source", source)
	}
	return nil
}

// Resolve resolves the "source" parameter
func (parser *Parser) Resolve(parameter string) string {
	if parameter != "source" {
		return ""
	}
	return "source"
}
