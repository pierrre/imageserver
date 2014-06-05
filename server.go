// Package imageserver provides an image server
package imageserver

// ImageServer represents an Image server
type ImageServer struct {
	Provider
	Processor // optional
}

// Get returns an Image for given Parameters
//
// The "source" parameter is required.
func (imageServer *ImageServer) Get(parameters Parameters) (*Image, error) {
	source, err := parameters.Get("source")
	if err != nil {
		return nil, NewError("Missing source parameter")
	}

	image, err := imageServer.Provider.Get(source, parameters)
	if err != nil {
		return nil, err
	}

	if imageServer.Processor != nil {
		image, err = imageServer.Processor.Process(image, parameters)
		if err != nil {
			return nil, err
		}
	}

	return image, nil
}

// ImageServerFunc is a ImageServer func
type ImageServerFunc func(parameters Parameters) (*Image, error)

// Get calls the func
func (f ImageServerFunc) Get(parameters Parameters) (*Image, error) {
	return f(parameters)
}

// ImageServerInterface represents an interface for an Image server
type ImageServerInterface interface {
	Get(Parameters) (*Image, error)
}
