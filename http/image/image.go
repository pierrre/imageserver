// Package image provides imageserver/http.Parser implementations for imageserver/image.
package image

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

// FormatParser is a imageserver/http.Parser implementation for imageserver/image.
//
// It takes the string "format" param from the HTTP URL query.
type FormatParser struct{}

// Parse implements imageserver/http.Parser.
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

// Resolve implements imageserver/http.Parser.
func (parser *FormatParser) Resolve(param string) string {
	if param == "format" {
		return "format"
	}
	return ""
}

// QualityParser is a imageserver/http.Parser implementation for imageserver/image.
//
// It takes the integer "quality" param from the HTTP URL query.
type QualityParser struct{}

// Parse implements imageserver/http.Parser.
func (parser *QualityParser) Parse(req *http.Request, params imageserver.Params) error {
	return imageserver_http.ParseQueryInt("quality", req, params)
}

// Resolve implements imageserver/http.Parser.
func (parser *QualityParser) Resolve(param string) string {
	if param == "quality" {
		return "quality"
	}
	return ""
}
