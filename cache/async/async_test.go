package async

import (
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestInterface(t *testing.T) {
	var _ imageserver_cache.Cache = &Cache{}
}

func TestGetSet(t *testing.T) {
	mapCache := cachetest.NewMapCache()

	setCallCh := make(chan struct{})
	funcCache := &cachetest.FuncCache{
		GetFunc: func(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return mapCache.Get(key, parameters)
		},
		SetFunc: func(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
			setCallCh <- struct{}{}
			return mapCache.Set(key, image, parameters)
		},
	}

	asyncCache := &Cache{
		Cache: funcCache,
	}

	err := asyncCache.Set("foo", testdata.Small, cachetest.ParametersEmpty)
	if err != nil {
		panic(err)
	}
	<-setCallCh
	_, err = asyncCache.Get("foo", cachetest.ParametersEmpty)
	if err != nil {
		panic(err)
	}
}

func TestSetErrFunc(t *testing.T) {
	funcCache := &cachetest.FuncCache{
		SetFunc: func(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
			return fmt.Errorf("error")
		},
	}

	errFuncCallCh := make(chan struct{})
	asyncCache := &Cache{
		Cache: funcCache,
		ErrFunc: func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters) {
			errFuncCallCh <- struct{}{}
		},
	}

	err := asyncCache.Set("foo", testdata.Small, cachetest.ParametersEmpty)
	if err != nil {
		panic(err)
	}
	<-errFuncCallCh
}
