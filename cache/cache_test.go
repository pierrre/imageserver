package cache_test

import (
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

var _ Cache = &IgnoreError{}

func TestIgnoreErrorGetSet(t *testing.T) {
	c := &IgnoreError{
		Cache: cachetest.NewMapCache(),
	}
	cachetest.TestGetSet(t, c)
}

func TestIgnoreErrorGetSetError(t *testing.T) {
	c := &IgnoreError{
		Cache: &Func{
			GetFunc: func(key string, params imageserver.Params) (*imageserver.Image, error) {
				return nil, fmt.Errorf("error")
			},
			SetFunc: func(key string, image *imageserver.Image, params imageserver.Params) error {
				return fmt.Errorf("error")
			},
		},
	}
	_, err := c.Get("test", imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set("test", testdata.Medium, imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
}

var _ Cache = &Async{}

func TestAsyncGetSet(t *testing.T) {
	mapCache := cachetest.NewMapCache()
	setCallCh := make(chan struct{})
	asyncCache := &Async{
		Cache: &Func{
			GetFunc: func(key string, params imageserver.Params) (*imageserver.Image, error) {
				return mapCache.Get(key, params)
			},
			SetFunc: func(key string, image *imageserver.Image, params imageserver.Params) error {
				err := mapCache.Set(key, image, params)
				setCallCh <- struct{}{}
				return err
			},
		},
	}

	err := asyncCache.Set("foo", testdata.Small, imageserver.Params{})
	if err != nil {
		panic(err)
	}
	<-setCallCh
	im, err := asyncCache.Get("foo", imageserver.Params{})
	if err != nil {
		panic(err)
	}
	if im == nil {
		t.Fatal("no image")
	}
}

var _ Cache = &Func{}
