package cache

import (
	"fmt"
	"testing"
)

func TestNewCacheMissError(t *testing.T) {
	key := "foobar"
	previousErr := fmt.Errorf("not found")

	err := NewMissError(key, previousErr)
	err.Error()
}
