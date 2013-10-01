package imageserver

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	text := "foo"
	previous := errors.New("bar")
	err := NewErrorWithPrevious(text, previous)
	if err.Text != text {
		t.Fatal("Invalid text")
	}
	if err.Previous != previous {
		t.Fatal("Invalid previous")
	}
	if err.Error() != text {
		t.Fatal("Invalid error message")
	}
}
