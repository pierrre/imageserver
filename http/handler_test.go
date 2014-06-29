package http

import (
	"net/http"

	"testing"
)

func TestTODO(t *testing.T) {
	t.Log("TODO")
}

func TestImageHTTPHandlerInterface(t *testing.T) {
	var _ http.Handler = &ImageHTTPHandler{}
}
