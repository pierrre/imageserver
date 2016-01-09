// Package gamma provides a imageserver/http.Parser implementation for imageserver/image/gamma.CorrectionProcessor.
package gamma

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
)

// CorrectionParser is a imageserver/http.Parser implementation for imageserver/image/gamma.CorrectionProcessor.
//
// It takes the boolean "gamma_correction" param from the HTTP URL query.
type CorrectionParser struct{}

// Parse implements imageserver/http.Parser.
func (parser *CorrectionParser) Parse(req *http.Request, params imageserver.Params) error {
	return imageserver_http.ParseQueryBool("gamma_correction", req, params)
}

// Resolve implements imageserver/http.Parser.
func (parser *CorrectionParser) Resolve(param string) string {
	if param == "gamma_correction" {
		return "gamma_correction"
	}
	return ""
}
