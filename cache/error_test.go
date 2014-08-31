package cache

import (
	"testing"
)

func TestMissErrorInterface(t *testing.T) {
	var _ error = &MissError{}
}

func TestMissError(t *testing.T) {
	err := &MissError{Key: "foobar"}
	err.Error()
}
