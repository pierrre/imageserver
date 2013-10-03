package http

import (
	"net/http"
)

// Http error
type Error struct {
	Code int
	Text string
}

func NewError(code int) *Error {
	return NewErrorWithText(code, http.StatusText(code))
}

func NewErrorWithText(code int, text string) *Error {
	return &Error{
		Code: code,
		Text: text,
	}
}

func (err *Error) Error() string {
	return err.Text
}
