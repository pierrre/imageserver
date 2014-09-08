package provider

import (
	"fmt"
)

// SourceError represents a source error
type SourceError struct {
	Message string
}

func (err *SourceError) Error() string {
	return fmt.Sprintf("source error: %s", err.Message)
}
