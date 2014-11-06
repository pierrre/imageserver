package http

import (
	"crypto/sha256"
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestHandlerInterface(t *testing.T) {
	var _ http.Handler = &Handler{}
}

func TestNewParametersHashETagFunc(t *testing.T) {
	f := NewParametersHashETagFunc(sha256.New)
	parameters := imageserver.Parameters{
		"foo": "bar",
	}
	f(parameters)
}
