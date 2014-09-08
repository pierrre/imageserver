package provider

import (
	"testing"
)

func TestSourceErrorInterface(t *testing.T) {
	var _ error = &SourceError{}
}

func TestSourceError(t *testing.T) {
	err := &SourceError{Message: "test"}
	err.Error()
}
