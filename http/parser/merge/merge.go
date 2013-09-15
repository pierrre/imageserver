// Merge http parser
package merge

import (
	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	"net/http"
)

// Merges multiple parsers
type MergeParser []imageserver_http.Parser

// Calls sequentially each parser
func (parser MergeParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	for _, subParser := range parser {
		err := subParser.Parse(request, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}
