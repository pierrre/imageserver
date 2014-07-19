package cache

import (
	"fmt"
)

// MissError represents a miss error
type MissError struct {
	Key string
}

func (err *MissError) Error() string {
	return fmt.Sprintf("cache miss for key \"%s\"", err.Key)
}
