package imageserver

import (
	"testing"
)

func TestError(t *testing.T) {
	text := "foo"
	err := NewError(text)
	if err.Text != text {
		t.Fatal("Invalid text")
	}
	if err.Error() != text {
		t.Fatal("Invalid error message")
	}
}
