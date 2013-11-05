// Chain http parser
package chain

import (
	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	"net/http"
)

// Merges multiple parsers
type ChainParser []imageserver_http.Parser

// Calls sequentially each parser
func (parser ChainParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	for _, subParser := range parser {
		err := subParser.Parse(request, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}
