// Source path http parser
package sourcepath

import (
	"github.com/pierrre/imageserver"
	"net/http"
	"net/url"
)

// Similar to the SourceParser, but takes the request's path, and concatenates it to the base url
type SourcePathParser struct {
	Base *url.URL
}

func (parser *SourcePathParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	s := *parser.Base
	source := &s
	source.Path += request.URL.Path
	parameters.Set("source", source)
	return nil
}
