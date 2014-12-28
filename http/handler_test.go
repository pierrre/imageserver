package http

import (
	"crypto/sha256"
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestHandlerInterface(t *testing.T) {
	var _ http.Handler = &Handler{}
}

func TestNewParamsHashETagFunc(t *testing.T) {
	NewParamsHashETagFunc(sha256.New)(imageserver.Params{
		"foo": "bar",
	})
}
