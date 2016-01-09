// Package image provides imageserver/http.Parser implementations for imageserver/image.
package image

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

// FormatParser is an http Parser that takes the "format" param from query.
type FormatParser struct{}

// Parse implements Parser.
func (parser *FormatParser) Parse(req *http.Request, params imageserver.Params) error {
	imageserver_http.ParseQueryString("format", req, params)
	if !params.Has("format") {
		return nil
	}
	format, err := params.GetString("format")
	if err != nil {
		return err
	}
	if format == "jpg" {
		format = "jpeg"
	}
	params.Set("format", format)
	return nil
}

// Resolve implements Parser.
func (parser *FormatParser) Resolve(param string) string {
	if param == "format" {
		return "format"
	}
	return ""
}

// QualityParser is an http Parser that takes the "quality" param from query.
type QualityParser struct{}

// Parse implements Parser.
func (parser *QualityParser) Parse(req *http.Request, params imageserver.Params) error {
	return imageserver_http.ParseQueryInt("quality", req, params)
}

// Resolve implements Parser.
func (parser *QualityParser) Resolve(param string) string {
	if param == "quality" {
		return "quality"
	}
	return ""
}
