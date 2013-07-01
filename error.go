package imageserver

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
