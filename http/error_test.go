package http

import (
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	err := NewError(http.StatusTeapot)
	text := "I'm a teapot"
	if err.Text != text {
		t.Fatal("invalid text")
	}
	if err.Error() != text {
		t.Fatal("invalid error message")
	}
}
