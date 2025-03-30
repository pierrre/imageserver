package hmac

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pierrre/imageserver"
)

const param = "expiration"

// Parser is a imageserver/http.Parser implementation for imageserver/image/expiration.Processor.
type Parser struct{}

// Parse implements imageserver/http.Parser.
func (prs *Parser) Parse(req *http.Request, params imageserver.Params) error {
	exp := req.URL.Query().Get(param)
	if exp == "" {
		return nil
	}

	expi, err := strconv.ParseInt(exp, 10, 0)
	if err != nil {
		return &imageserver.ParamError{
			Param:   param,
			Message: fmt.Sprintf("expected format '<int>': %s", err),
		}
	}
	params.Set(param, expi)
	return nil
}

// Resolve implements imageserver/http.Parser.
func (prs *Parser) Resolve(p string) (httpParam string) {
	if p == param || strings.HasPrefix(p, param+".") {
		return param
	}
	return ""
}
