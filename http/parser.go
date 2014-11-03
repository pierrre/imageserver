package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pierrre/imageserver"
)

// Parser represent an HTTP Request parser.
type Parser interface {
	// Parse parses a Request and fill Parameters.
	Parse(*http.Request, imageserver.Parameters) error

	// Resolve resolves a parameter to a HTTP parameter.
	// It returns the resolved HTTP parameter, or an empty string.
	Resolve(parameter string) (httpParameter string)
}

// ListParser represents a list of HTTP Parser
type ListParser []Parser

// Parse parses an http Request with sub Parsers in sequential order
func (lp ListParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	for _, subParser := range lp {
		err := subParser.Parse(request, parameters)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve resolves the parameter with sub Parsers in sequential order
func (lp ListParser) Resolve(parameter string) string {
	for _, subParser := range lp {
		httpParameter := subParser.Resolve(parameter)
		if httpParameter != "" {
			return httpParameter
		}
	}
	return ""
}

// SourceParser represents an http Parser that takes the "source" parameter from query
type SourceParser struct{}

// Parse takes the "source" parameter from query
func (parser *SourceParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	query := request.URL.Query()
	source := query.Get("source")
	if source != "" {
		parameters.Set("source", source)
	}
	return nil
}

// Resolve resolves the "source" parameter
func (parser *SourceParser) Resolve(parameter string) string {
	if parameter != "source" {
		return ""
	}
	return "source"
}

// SourcePathParser represents an HTTP Parser that takes the "source" parameter from the path
type SourcePathParser struct {
}

// Parse takes the "source" parameter from the path
func (parser *SourcePathParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	parameters.Set("source", request.URL.Path)
	return nil
}

// Resolve resolves the "source" parameter
func (parser *SourcePathParser) Resolve(parameter string) string {
	return ""
}

// SourceURLParser is a Parser that takes the "source" from the sub Parser and adds it to the Base URL.
type SourceURLParser struct {
	Parser
	Base *url.URL
}

// Parse implements Parser
func (parser *SourceURLParser) Parse(request *http.Request, parameters imageserver.Parameters) error {
	err := parser.Parser.Parse(request, parameters)
	if err != nil {
		return err
	}
	source, err := parameters.Get("source")
	if err != nil {
		return &imageserver.ParameterError{Parameter: "source", Message: "missing"}
	}

	u := copyURL(parser.Base)
	u.Path += fmt.Sprint(source)
	parameters.Set("source", u)

	return nil
}

// Resolve implements Parser
func (parser *SourceURLParser) Resolve(parameter string) string {
	return parser.Parser.Resolve(parameter)
}
