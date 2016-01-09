package http

import (
	"fmt"
	"net/http"
)

// Error is a HTTP error.
//
// It is supported by Handler.
type Error struct {
	Code int
	Text string
}

// NewErrorDefaultText creates an Error with the default message associated with the code.
func NewErrorDefaultText(code int) *Error {
	return &Error{Code: code, Text: http.StatusText(code)}
}

func (err *Error) Error() string {
	return fmt.Sprintf("http error %d: %s", err.Code, err.Text)
}
