package provider

import (
	"testing"
)

var _ error = &SourceError{}

func TestSourceError(t *testing.T) {
	err := &SourceError{Message: "test"}
	_ = err.Error()
}
