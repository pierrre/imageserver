package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pierrre/imageserver"
)

// Parser represent an HTTP Request parser.
type Parser interface {
	// Parse parses a Request and fill Params.
	Parse(*http.Request, imageserver.Params) error

	// Resolve resolves a param to a HTTP param.
	// It returns the resolved HTTP param, or an empty string.
	Resolve(param string) (httpParam string)
}

// ListParser represents a list of HTTP Parser
type ListParser []Parser

// Parse parses an http Request with sub Parsers in sequential order
func (lp ListParser) Parse(request *http.Request, params imageserver.Params) error {
	for _, subParser := range lp {
		err := subParser.Parse(request, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve resolves the param with sub Parsers in sequential order
func (lp ListParser) Resolve(param string) string {
	for _, subParser := range lp {
		httpParam := subParser.Resolve(param)
		if httpParam != "" {
			return httpParam
		}
	}
	return ""
}

// SourceParser represents an http Parser that takes the "source" param from query
type SourceParser struct{}

// Parse takes the "source" param from query
func (parser *SourceParser) Parse(request *http.Request, params imageserver.Params) error {
	query := request.URL.Query()
	source := query.Get("source")
	if source != "" {
		params.Set("source", source)
	}
	return nil
}

// Resolve resolves the "source" param
func (parser *SourceParser) Resolve(param string) string {
	if param != "source" {
		return ""
	}
	return "source"
}

// SourcePathParser represents an HTTP Parser that takes the "source" param from the path
type SourcePathParser struct {
}

// Parse takes the "source" param from the path
func (parser *SourcePathParser) Parse(request *http.Request, params imageserver.Params) error {
	params.Set("source", request.URL.Path)
	return nil
}

// Resolve resolves the "source" param
func (parser *SourcePathParser) Resolve(param string) string {
	return ""
}

// SourceURLParser is a Parser that takes the "source" from the sub Parser and adds it to the Base URL.
type SourceURLParser struct {
	Parser
	Base *url.URL
}

// Parse implements Parser
func (parser *SourceURLParser) Parse(request *http.Request, params imageserver.Params) error {
	err := parser.Parser.Parse(request, params)
	if err != nil {
		return err
	}
	source, err := params.Get("source")
	if err != nil {
		return &imageserver.ParamError{Param: "source", Message: "missing"}
	}

	u := copyURL(parser.Base)
	u.Path += fmt.Sprint(source)
	params.Set("source", u)

	return nil
}

// Resolve implements Parser
func (parser *SourceURLParser) Resolve(param string) string {
	return parser.Parser.Resolve(param)
}

// ParseQueryString takes the param from the query string and add it to params.
func ParseQueryString(param string, request *http.Request, params imageserver.Params) {
	s := request.URL.Query().Get(param)
	if s != "" {
		params.Set(param, s)
	}
}

// ParseQueryInt takes the param from the query string, parse it as an int and add it to params.
func ParseQueryInt(param string, request *http.Request, params imageserver.Params) error {
	s := request.URL.Query().Get(param)
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return newParseTypeParamError(param, "string", err)
	}
	params.Set(param, i)
	return nil
}

// ParseQueryFloat takes the param from the query string, parse it as a float64 and add it to params.
func ParseQueryFloat(param string, request *http.Request, params imageserver.Params) error {
	s := request.URL.Query().Get(param)
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return newParseTypeParamError(param, "float", err)
	}
	params.Set(param, f)
	return nil
}

// ParseQueryBool takes the param from the query string, parse it as an bool and add it to params.
func ParseQueryBool(param string, request *http.Request, params imageserver.Params) error {
	s := request.URL.Query().Get(param)
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return newParseTypeParamError(param, "bool", err)
	}
	params.Set(param, b)
	return nil
}

func newParseTypeParamError(param string, parseType string, parseErr error) *imageserver.ParamError {
	return &imageserver.ParamError{
		Param:   param,
		Message: fmt.Sprintf("parse %s: %s", parseType, parseErr.Error()),
	}
}
