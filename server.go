package imageserver

type Server struct {
	Cache     Cache
	Source    Source
	Converter Converter
}

func (server *Server) GetImage(parameters Parameters) (image *Image, err error) {
	cacheKey := parameters.Hash()

	if server.Cache != nil {
		image, err = server.Cache.Get(cacheKey)
		if err == nil {
			return
		}
	}

	sourceImage, err := server.getSourceImage(parameters)
	if err != nil {
		return
	}

	image, err = server.convertImage(sourceImage, parameters)
	if err != nil {
		return
	}

	if server.Cache != nil {
		go func() {
			_ = server.Cache.Set(cacheKey, image)
		}()
	}

	return
}

func (server *Server) getSourceImage(parameters Parameters) (image *Image, err error) {
	sourceId, err := parameters.GetString("source")
	if err != nil {
		return
	}

	image, err = server.Source.Get(sourceId)
	if err != nil {
		return
	}

	return
}

func (server *Server) convertImage(sourceImage *Image, parameters Parameters) (image *Image, err error) {
	if server.Converter != nil {
		image, err = server.Converter.Convert(sourceImage, parameters)
	} else {
		image = sourceImage
	}

	return
}
