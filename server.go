package imageserver

import (
	"fmt"
)

type Server struct {
	Cache     Cache
	Provider  Provider
	Processor Processor
}

func (server *Server) Get(parameters Parameters) (image *Image, err error) {
	cacheKey := parameters.Hash()

	if server.Cache != nil {
		image, err = server.Cache.Get(cacheKey, parameters)
		if err == nil {
			return
		}
	}

	sourceImage, err := server.getSource(parameters)
	if err != nil {
		return
	}

	image, err = server.process(sourceImage, parameters)
	if err != nil {
		return
	}

	if server.Cache != nil {
		go func() {
			_ = server.Cache.Set(cacheKey, image, parameters)
		}()
	}

	return
}

func (server *Server) getSource(parameters Parameters) (image *Image, err error) {
	sourceId, _ := parameters.GetString("source")
	if len(sourceId) == 0 {
		err = fmt.Errorf("Source is missing")
		return
	}

	image, err = server.Provider.Get(sourceId, parameters)
	if err != nil {
		return
	}

	return
}

func (server *Server) process(sourceImage *Image, parameters Parameters) (image *Image, err error) {
	if server.Processor != nil {
		image, err = server.Processor.Process(sourceImage, parameters)
	} else {
		image = sourceImage
	}

	return
}
