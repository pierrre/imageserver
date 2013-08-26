package imageserver

// Error displayable to the end user.
//
// It is used in the http server: the errors of this type are displayed, the others are shown as "internal error".
type Error struct {
	text string
}

func NewError(text string) *Error {
	return &Error{
		text: text,
	}
}

func (err *Error) Error() string {
	return err.text
}
