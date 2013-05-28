package imageproxy

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

type Cache interface {
	Get(key string) (*Image, error)
	Set(key string, image *Image) error
}

func HashCacheKey(key string) string {
	hash := md5.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	hashedKey := hex.EncodeToString(data)
	return hashedKey
}
