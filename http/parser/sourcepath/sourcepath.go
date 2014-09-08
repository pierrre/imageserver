// Package sourcepath provides an http Parser that takes the "source" parameter from the path
package sourcepath

import (
	"net/http"
	"net/url"

	"github.com/pierrre/imageserver"
)

// Parser represents an HTTP Parser that takes the "source" parameter from the path
type Parser struct {
	Base *url.URL
}

// Parse takes the "source" parameter from the path
func (parser *Parser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	s := *parser.Base
	source := &s
	source.Path += request.URL.Path
	parameters.Set("source", source)
	return nil
}
