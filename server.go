// Package imageserver provides an image server
package imageserver

// Server represents an Image server
type Server struct {
	Cache                                           // optional
	CacheKeyFunc func(parameters Parameters) string // optional
	Provider
	Processor // optional
}

// Get returns an Image for given Parameters
//
// The "source" parameter is required.
//
// Steps
//
// - get image from the cache and return it if available
//
// - get the image from the provider
//
// - process the image
//
// - store the image in the cache
func (server *Server) Get(parameters Parameters) (*Image, error) {
	var cacheKey string
	if server.Cache != nil && server.CacheKeyFunc != nil {
		cacheKey = server.CacheKeyFunc(parameters)

		image, err := server.Cache.Get(cacheKey, parameters)

		if err == nil {
			return image, nil
		}
	}

	source, err := parameters.Get("source")
	if err != nil {
		return nil, NewError("Missing source parameter")
	}

	image, err := server.Provider.Get(source, parameters)
	if err != nil {
		return nil, err
	}

	if server.Processor != nil {
		image, err = server.Processor.Process(image, parameters)
		if err != nil {
			return nil, err
		}
	}

	if server.Cache != nil {
		server.Cache.Set(cacheKey, image, parameters)
		// TODO handle errors properly
	}

	return image, nil
}
