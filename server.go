// Image server
package imageserver

// Image server
type Server struct {
	Cache     Cache // optional
	Provider  Provider
	Processor Processor // optional
}

// Get an image
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
	cacheKey := parameters.Hash()

	if server.Cache != nil {
		if image, err := server.Cache.Get(cacheKey, parameters); err == nil {
			return image, nil
		}
	}

	sourceImage, err := server.getSource(parameters)
	if err != nil {
		return nil, err
	}

	image, err := server.process(sourceImage, parameters)
	if err != nil {
		return nil, err
	}

	if server.Cache != nil {
		go func() {
			_ = server.Cache.Set(cacheKey, image, parameters)
		}()
	}

	return image, nil
}

func (server *Server) getSource(parameters Parameters) (*Image, error) {
	source, err := parameters.Get("source")
	if err != nil {
		err = NewError("Missing source parameter")
		return nil, err
	}

	image, err := server.Provider.Get(source, parameters)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (server *Server) process(sourceImage *Image, parameters Parameters) (*Image, error) {
	if server.Processor == nil {
		return sourceImage, nil
	}

	image, err := server.Processor.Process(sourceImage, parameters)

	return image, err
}
