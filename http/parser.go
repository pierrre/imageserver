package http

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

// Http parser interface
//
// Parses the request and fills the parameters.
type Parser interface {
	Parse(*http.Request, imageserver.Parameters) error
}
