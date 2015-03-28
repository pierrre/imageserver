package cache_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

var _ error = &MissError{}

func TestMissError(t *testing.T) {
	err := &MissError{Key: "foobar"}
	err.Error()
}

var _ Cache = &Async{}

func TestAsyncGetSet(t *testing.T) {
	mapCache := cachetest.NewMapCache()

	setCallCh := make(chan struct{})
	funcCache := &Func{
		GetFunc: func(key string, params imageserver.Params) (*imageserver.Image, error) {
			return mapCache.Get(key, params)
		},
		SetFunc: func(key string, image *imageserver.Image, params imageserver.Params) error {
			setCallCh <- struct{}{}
			return mapCache.Set(key, image, params)
		},
	}

	asyncCache := &Async{
		Cache: funcCache,
	}

	err := asyncCache.Set("foo", testdata.Small, imageserver.Params{})
	if err != nil {
		panic(err)
	}
	<-setCallCh
	_, err = asyncCache.Get("foo", imageserver.Params{})
	if err != nil {
		panic(err)
	}
}

func TestAsyncSetErrFunc(t *testing.T) {
	funcCache := &Func{
		SetFunc: func(key string, image *imageserver.Image, params imageserver.Params) error {
			return fmt.Errorf("error")
		},
	}

	errFuncCallCh := make(chan struct{})
	asyncCache := &Async{
		Cache: funcCache,
		ErrFunc: func(err error, key string, image *imageserver.Image, params imageserver.Params) {
			errFuncCallCh <- struct{}{}
		},
	}

	err := asyncCache.Set("foo", testdata.Small, imageserver.Params{})
	if err != nil {
		panic(err)
	}
	<-errFuncCallCh
}

var _ Cache = &Func{}

var _ Cache = &Prefix{}

func TestPrefixSet(t *testing.T) {
	c := cachetest.NewMapCache()
	pc := &Prefix{Cache: c, Prefix: "foo"}

	err := pc.Set("bar", testdata.Medium, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Get("foobar", imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPrefixGet(t *testing.T) {
	c := cachetest.NewMapCache()
	pc := &Prefix{Cache: c, Prefix: "foo"}

	err := c.Set("foobar", testdata.Medium, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = pc.Get("bar", imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}
