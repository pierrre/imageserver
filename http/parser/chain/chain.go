// Package chain provides a chained http Parser
package chain

import (
	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	"net/http"
)

// ChainParser represents a chained http Parser
type ChainParser []imageserver_http.Parser

// Parse parses an http Request with sub Parsers in sequential order
func (parser ChainParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	for _, subParser := range parser {
		err := subParser.Parse(request, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}
