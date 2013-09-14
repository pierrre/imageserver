package http

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

type Parser interface {
	Parse(*http.Request, imageserver.Parameters) error
}
