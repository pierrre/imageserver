// Package imageserver provides an image server
package imageserver

// ImageServer represents an Image server
type ImageServer struct {
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
func (imageServer *ImageServer) Get(parameters Parameters) (*Image, error) {
	var cacheKey string
	if imageServer.Cache != nil && imageServer.CacheKeyFunc != nil {
		cacheKey = imageServer.CacheKeyFunc(parameters)

		image, err := imageServer.Cache.Get(cacheKey, parameters)

		if err == nil {
			return image, nil
		}
	}

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

	if imageServer.Cache != nil {
		imageServer.Cache.Set(cacheKey, image, parameters)
		// TODO handle errors properly
	}

	return image, nil
}
