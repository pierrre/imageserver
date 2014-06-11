package provider

import (
	"fmt"

	"github.com/pierrre/imageserver"
)

// ProviderImageServer is a Provider ImageServer
type ProviderImageServer struct {
	Provider Provider
}

// Get get an Image from the Provider using the "source" parameter
func (pis *ProviderImageServer) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	source, err := parameters.Get("source")
	if err != nil {
		return nil, newSourceParameterError("missing")
	}

	image, err := pis.Provider.Get(source, parameters)
	if err != nil {
		if err, ok := err.(*SourceError); ok {
			return nil, newSourceParameterError(err.Message)
		}
		return nil, err
	}

	return image, nil
}

func newSourceParameterError(message string) error {
	return &imageserver.ParameterError{Parameter: "source", Message: message}
}

// SourceError represents a source error
type SourceError struct {
	Message string
}

func (err *SourceError) Error() string {
	return fmt.Sprintf("source error: %s", err.Message)
}

// NewSourceError creates a new SourceError
func NewSourceError(message string) *SourceError {
	return &SourceError{Message: message}
}
