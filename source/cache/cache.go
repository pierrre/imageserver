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

func (source *CacheSource) Get(sourceId string) (image *imageserver.Image, err error) {
	cacheKey := source.getCacheKey(sourceId)

	image, e := source.Cache.Get(cacheKey)
	if e == nil {
		return
	}

	image, err = source.Source.Get(sourceId)
	if err != nil {
		return
	}

	go func() {
		_ = source.Cache.Set(cacheKey, image)
	}()

	return
}

func (server *CacheSource) getCacheKey(key string) string {
	hash := sha256.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}
