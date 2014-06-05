package async

import (
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestSet(t *testing.T) {
	funcCache := &cachetest.FuncCache{}
	funcCache.GetFunc = func(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
		return nil, imageserver_cache.NewCacheMissError(key, funcCache, nil)
	}
	setCallCh := make(chan struct{})
	funcCache.SetFunc = func(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
		setCallCh <- struct{}{}
		return fmt.Errorf("error")
	}

	asyncCache := &AsyncCache{
		Cache: funcCache,
	}
	errFuncCallCh := make(chan struct{})
	asyncCache.ErrFunc = func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters) {
		errFuncCallCh <- struct{}{}
	}

	asyncCache.Set("foo", testdata.Small, cachetest.ParametersEmpty)
	<-setCallCh
	<-errFuncCallCh
}
