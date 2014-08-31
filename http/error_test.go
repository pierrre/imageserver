package http

import (
	"net/http"
	"testing"
)

func TestErrorInterface(t *testing.T) {
	var _ error = &Error{}
}

func TestError(t *testing.T) {
	err := NewErrorDefaultText(http.StatusTeapot)
	text := "I'm a teapot"
	if err.Text != text {
		t.Fatal("invalid text")
	}
	err.Error()
}
