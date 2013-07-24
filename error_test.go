package imageserver

import (
	"testing"
)

func TestError(t *testing.T) {
	text := "foo"
	err := NewError(text)
	if err.Error() != text {
		t.FailNow()
	}
}
