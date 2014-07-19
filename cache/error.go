package cache

import (
	"fmt"
)

// MissError represents a miss error
type MissError struct {
	Key      string
	Previous error
}

// NewMissError creates a new MissError
func NewMissError(key string, previous error) *MissError {
	return &MissError{
		Key:      key,
		Previous: previous,
	}
}

func (err *MissError) Error() string {
	s := fmt.Sprintf("cache miss for key \"%s\"", err.Key)
	if err.Previous != nil {
		s = fmt.Sprintf("%s (%s)", s, err.Previous)
	}
	return s
}
