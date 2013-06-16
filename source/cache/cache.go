package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/pierrre/imageserver"
	"io"
)

type CacheSource struct {
	Cache  imageserver.Cache
	Source imageserver.Source
}

func (source *CacheSource) Get(sourceId string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	cacheKey := source.getCacheKey(sourceId)

	image, e := source.Cache.Get(cacheKey, parameters)
	if e == nil {
		return
	}

	image, err = source.Source.Get(sourceId, parameters)
	if err != nil {
		return
	}

	go func() {
		_ = source.Cache.Set(cacheKey, image, parameters)
	}()

	return
}

func (server *CacheSource) getCacheKey(key string) string {
	hash := sha256.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}
