// Package list provides a list of HTTP Parser
package list

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

// Parser represents a list of HTTP Parser
type Parser []imageserver_http.Parser

// Parse parses an http Request with sub Parsers in sequential order
func (parser Parser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	for _, subParser := range parser {
		err := subParser.Parse(request, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve resolves the parameter with sub Parsers in sequential order
func (parser Parser) Resolve(parameter string) string {
	for _, subParser := range parser {
		httpParameter := imageserver_http.Resolve(subParser, parameter)
		if httpParameter != "" {
			return httpParameter
		}
	}
	return ""
}
