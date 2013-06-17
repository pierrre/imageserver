package http

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

type Parser interface {
	Parse(request *http.Request, parameters imageserver.Parameters) error
}
