package cache_test

import (
	"crypto/sha256"
	"errors"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestServerInterface(t *testing.T) {
	var _ imageserver.Server = &Server{}
}

func TestServer(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		Cache: cachetest.NewMapCache(),
		KeyGenerator: KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
			return "test"
		}),
	}
	image1, err := s.Get(imageserver.Parameters{})
	if err != nil {
		t.Fatal(err)
	}
	image2, err := s.Get(imageserver.Parameters{})
	if err != nil {
		t.Fatal(err)
	}
	if !imageserver.ImageEqual(image1, image2) {
		t.Fatal("not equal")
	}
}

func TestServerErrorServer(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, errors.New("error")
		}),
		Cache: cachetest.NewMapCache(),
		KeyGenerator: KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
			return "test"
		}),
	}
	_, err := s.Get(imageserver.Parameters{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorCacheSet(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		Cache: &cachetest.FuncCache{
			GetFunc: func(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
				return nil, &MissError{Key: key}
			},
			SetFunc: func(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
				return errors.New("error")
			},
		},
		KeyGenerator: KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
			return "test"
		}),
	}
	_, err := s.Get(imageserver.Parameters{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestKeyGeneratorFuncInterface(t *testing.T) {
	var _ KeyGenerator = KeyGeneratorFunc(nil)
}

func TestNewParametersHashKeyGenerator(t *testing.T) {
	g := NewParametersHashKeyGenerator(sha256.New)
	parameters := imageserver.Parameters{
		"foo": "bar",
	}
	g.GetKey(parameters)
}
