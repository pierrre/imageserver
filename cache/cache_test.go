package cache_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestMissErrorInterface(t *testing.T) {
	var _ error = &MissError{}
}

func TestMissError(t *testing.T) {
	err := &MissError{Key: "foobar"}
	err.Error()
}

func TestListInterface(t *testing.T) {
	var _ Cache = List{}
}

func TestAsyncInterface(t *testing.T) {
	var _ Cache = &Async{}
}

func TestAsyncGetSet(t *testing.T) {
	mapCache := cachetest.NewMapCache()

	setCallCh := make(chan struct{})
	funcCache := &Func{
		GetFunc: func(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
			return mapCache.Get(key, parameters)
		},
		SetFunc: func(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
			setCallCh <- struct{}{}
			return mapCache.Set(key, image, parameters)
		},
	}

	asyncCache := &Async{
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

func TestAsyncSetErrFunc(t *testing.T) {
	funcCache := &Func{
		SetFunc: func(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
			return fmt.Errorf("error")
		},
	}

	errFuncCallCh := make(chan struct{})
	asyncCache := &Async{
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

func TestFuncInterface(t *testing.T) {
	var _ Cache = &Func{}
}
