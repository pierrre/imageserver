package http

import (
	"net/http"

	"github.com/pierrre/imageserver"
)

// Parser parses a Request and fill Parameters.
type Parser interface {
	Parse(*http.Request, imageserver.Parameters) error
}

// ParserFunc is a Parser func
type ParserFunc func(*http.Request, imageserver.Parameters) error

// Parse calls the func
func (f ParserFunc) Parse(request *http.Request, parameters imageserver.Parameters) error {
	return f(request, parameters)
}
