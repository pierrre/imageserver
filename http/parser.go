package http

import (
	"net/http"

	"github.com/pierrre/imageserver"
)

// Parser represents a http parser
//
// It parses the Tequest and fill Parameters.
type Parser interface {
	Parse(*http.Request, imageserver.Parameters) error
}
