package cache

import (
	"testing"
)

func TestCacheMissError(t *testing.T) {
	err := &MissError{Key: "foobar"}
	err.Error()
}
