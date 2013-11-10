package http

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

// Parser represents a http parser
//
// It parses the Tequest and fill Parameters.
type Parser interface {
	Parse(*http.Request, imageserver.Parameters) error
}
