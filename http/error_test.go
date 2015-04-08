package http

import (
	"net/http"
	"testing"
)

var _ error = &Error{}

func TestError(t *testing.T) {
	err := NewErrorDefaultText(http.StatusTeapot)
	text := "I'm a teapot"
	if err.Text != text {
		t.Fatal("invalid text")
	}
	_ = err.Error()
}
