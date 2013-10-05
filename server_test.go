package imageserver

import (
	"testing"
)

func TestServerErrorMissingSource(t *testing.T) {
	server := &Server{}
	parameters := make(Parameters)
	_, err := server.Get(parameters)
	if err == nil {
		t.Fatal("No error")
	}
}
