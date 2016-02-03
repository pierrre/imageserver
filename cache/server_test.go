package cache_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestServer(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		Cache: cachetest.NewMapCache(),
		KeyGenerator: KeyGeneratorFunc(func(params imageserver.Params) string {
			return "test"
		}),
	}
	image1, err := s.Get(imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	image2, err := s.Get(imageserver.Params{})
	if err != nil {
		t.Fatal(err)
	}
	if !imageserver.ImageEqual(image1, image2) {
		t.Fatal("not equal")
	}
}

func TestServerErrorCacheGet(t *testing.T) {
	s := &Server{
		Cache: &Func{
			GetFunc: func(key string, params imageserver.Params) (*imageserver.Image, error) {
				return nil, fmt.Errorf("error")
			},
		},
		KeyGenerator: KeyGeneratorFunc(func(params imageserver.Params) string {
			return "test"
		}),
	}
	_, err := s.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorServer(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
		Cache: cachetest.NewMapCache(),
		KeyGenerator: KeyGeneratorFunc(func(params imageserver.Params) string {
			return "test"
		}),
	}
	_, err := s.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorCacheSet(t *testing.T) {
	s := &Server{
		Server: imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		Cache: &Func{
			GetFunc: func(key string, params imageserver.Params) (*imageserver.Image, error) {
				return nil, nil
			},
			SetFunc: func(key string, image *imageserver.Image, params imageserver.Params) error {
				return fmt.Errorf("error")
			},
		},
		KeyGenerator: KeyGeneratorFunc(func(params imageserver.Params) string {
			return "test"
		}),
	}
	_, err := s.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

var _ KeyGenerator = KeyGeneratorFunc(nil)

func TestNewParamsHashKeyGenerator(t *testing.T) {
	NewParamsHashKeyGenerator(sha256.New).GetKey(imageserver.Params{
		"foo": "bar",
	})
}

var _ KeyGenerator = &PrefixKeyGenerator{}

func TestPrefixKeyGenerator(t *testing.T) {
	g := &PrefixKeyGenerator{
		KeyGenerator: KeyGeneratorFunc(func(params imageserver.Params) string {
			return "bar"
		}),
		Prefix: "foo",
	}
	key := g.GetKey(imageserver.Params{})
	if key != "foobar" {
		t.Fatal("not equal")
	}
}
