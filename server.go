package imageserver

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

type Server struct {
	Cache       Cache
	SourceCache Cache
	Source      Source
	Converter   Converter
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

	cacheKey := server.hashCacheKey(sourceId) //TODO cache source provider

	if server.SourceCache != nil {
		image, _ = server.SourceCache.Get(cacheKey)
		if image != nil {
			return
		}
	}

	image, err = server.Source.Get(sourceId)
	if err != nil {
		return
	}

	if server.SourceCache != nil {
		go func() {
			_ = server.SourceCache.Set(cacheKey, image)
		}()
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

func (server *Server) hashCacheKey(key string) string {
	hash := md5.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	hashedKey := hex.EncodeToString(data)
	return hashedKey
}
