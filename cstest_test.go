package imageserver

import (
	"testing"

	"github.com/pierrre/cstest"
)

func TestCSGofmt(t *testing.T) {
	cstest.RunGofmt(t)
}

func TestCSGoToolVet(t *testing.T) {
	cstest.RunGoToolVet(t)
}
