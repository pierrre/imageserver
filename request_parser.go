package imageserver

import (
	"net/http"
)

type RequestParser interface {
	ParseRequest(request *http.Request) (parameters Parameters, err error)
}
