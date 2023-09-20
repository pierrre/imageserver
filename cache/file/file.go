// Package file provides a disk based cache.
package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pierrre/imageserver"
)

// Cache is a implementation of disk based cache system.
type Cache struct {
	Path string
	mu   sync.RWMutex
}

// Get implements imageserver/cache.Cache.
func (cache *Cache) Get(key string, params imageserver.Params) (*imageserver.Image, error) {
	if cache.Path == "" {
		return nil, fmt.Errorf("file cache path is not")
	}
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	data, err := cache.getData(key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	im := new(imageserver.Image)
	if err := im.UnmarshalBinaryNoCopy(data); err != nil {
		return nil, err
	}
	return im, nil
}

func (cache *Cache) getData(key string) ([]byte, error) {
	item, err := os.ReadFile(filepath.Join(cache.Path, key))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Set implements imageserver/cache.Cache.
func (cache *Cache) Set(key string, im *imageserver.Image, params imageserver.Params) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	data, err := im.MarshalBinary()
	if err != nil {
		return err
	}
	return cache.setData(key, data)
}

func (cache *Cache) setData(key string, data []byte) error {
	return os.WriteFile(filepath.Join(cache.Path, key), data, 0644)
}
