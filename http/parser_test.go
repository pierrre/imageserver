package http

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestParserFuncInterface(t *testing.T) {
	var _ Parser = ParserFunc(nil)
}

func TestParserFunc(t *testing.T) {
	pf := ParserFunc(func(request *http.Request, parameters imageserver.Parameters) error {
		return nil
	})
	pf.Parse(&http.Request{}, imageserver.Parameters{})
}
