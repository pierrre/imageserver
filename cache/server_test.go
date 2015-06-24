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

var _ imageserver.Server = &Server{}

func TestServer(t *testing.T) {
	s := &Server{
		Server:       &imageserver.StaticServer{Image: testdata.Medium},
		Cache:        cachetest.NewMapCache(),
		KeyGenerator: StringKeyGenerator("test"),
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
				return nil, errors.New("error")
			},
		},
		KeyGenerator: StringKeyGenerator("test"),
	}
	_, err := s.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorServer(t *testing.T) {
	s := &Server{
		Server:       &imageserver.StaticServer{Error: errors.New("error")},
		Cache:        cachetest.NewMapCache(),
		KeyGenerator: StringKeyGenerator("test"),
	}
	_, err := s.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorCacheSet(t *testing.T) {
	s := &Server{
		Server: &imageserver.StaticServer{Image: testdata.Medium},
		Cache: &Func{
			GetFunc: func(key string, params imageserver.Params) (*imageserver.Image, error) {
				return nil, nil
			},
			SetFunc: func(key string, image *imageserver.Image, params imageserver.Params) error {
				return errors.New("error")
			},
		},
		KeyGenerator: StringKeyGenerator("test"),
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

var _ KeyGenerator = StringKeyGenerator("")

func TestStringKeyGenerator(t *testing.T) {
	g := StringKeyGenerator("foo")
	key := g.GetKey(imageserver.Params{})
	if key != "foo" {
		t.Fatal("not equal")
	}
}
