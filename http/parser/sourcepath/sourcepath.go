package sourcepath

import (
	"github.com/pierrre/imageserver"
	"net/http"
	"net/url"
)

type SourcePathParser struct {
	Base *url.URL
}

func (parser *SourcePathParser) Parse(request *http.Request, parameters imageserver.Parameters) (err error) {
	s := *parser.Base
	source := &s
	source.Path += request.URL.Path
	parameters.Set("source", source)
	return
}
