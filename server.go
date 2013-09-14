package imageserver

type Server struct {
	Cache     Cache
	Provider  Provider
	Processor Processor
}

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
