// Package list provides a list of http Parser
package list

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

// ListParser represents a list of http Parser
type ListParser []imageserver_http.Parser

// Parse parses an http Request with sub Parsers in sequential order
func (parser ListParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	for _, subParser := range parser {
		err := subParser.Parse(request, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}
