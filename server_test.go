package imageserver

import (
	"testing"
)

func TestServerGetSourceErrorMissingSource(t *testing.T) {
	server := &Server{}
	parameters := make(Parameters)
	_, err := server.getSource(parameters)
	if err == nil {
		t.Fatal("No error")
	}
}
