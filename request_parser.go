package imageproxy

import (
	"net/http"
	"net/url"
)

type RequestParser interface {
	ParseRequest(request *http.Request) (source *url.URL, parameters *Parameters, err error)
}
