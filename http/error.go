package http

import (
	"net/http"
)

// Error represent an http error
type Error struct {
	Code int
	Text string
}

// NewError creates an Error with the message associated with the code
func NewError(code int) *Error {
	return NewErrorWithText(code, http.StatusText(code))
}

// NewErrorWithText creates an Error
func NewErrorWithText(code int, text string) *Error {
	return &Error{
		Code: code,
		Text: text,
	}
}

func (err *Error) Error() string {
	return err.Text
}
