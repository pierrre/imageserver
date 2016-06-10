package groupcache

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/golang/groupcache"
	"github.com/pierrre/compare"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_source "github.com/pierrre/imageserver/source"
	"github.com/pierrre/imageserver/testdata"
)

const (
	testSize = 100 * (1 << 20)
)

var _ imageserver.Server = &Server{}

func TestServer(t *testing.T) {
	srv := newTestServer(
		imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		imageserver_cache.KeyGeneratorFunc(func(params imageserver.Params) string {
			return "test"
		}),
	)
	im, err := srv.Get(context.Background(), imageserver.Params{
		imageserver_source.Param: testdata.MediumFileName,
	})
	if err != nil {
		t.Fatal(err)
	}
	diff := compare.Compare(im, testdata.Medium)
	if len(diff) != 0 {
		t.Fatalf("images not equal, diff:\n%+v", diff)
	}
}

func TestServerErrorGroup(t *testing.T) {
	srv := newTestServer(
		imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
		imageserver_cache.KeyGeneratorFunc(func(params imageserver.Params) string {
			return "test"
		}),
	)
	_, err := srv.Get(context.Background(), imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestServerErrorImageUnmarshal(t *testing.T) {
	srv := &Server{
		Group: groupcache.NewGroup(
			newTestServerName(),
			testSize,
			groupcache.GetterFunc(func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
				return dest.SetBytes(nil)
			}),
		),
		KeyGenerator: imageserver_cache.KeyGeneratorFunc(func(params imageserver.Params) string {
			return "test"
		}),
	}
	_, err := srv.Get(context.Background(), imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func newTestServer(srv imageserver.Server, kg imageserver_cache.KeyGenerator) *Server {
	return NewServer(srv, kg, newTestServerName(), testSize)
}

func newTestServerName() string {
	return fmt.Sprintf("test_%d_%d", time.Now().UnixNano(), rand.Int63())
}

var _ groupcache.Getter = &Getter{}

func TestGetter(t *testing.T) {
	ctx := &Context{
		Params: imageserver.Params{},
	}
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
	}
	err := gt.Get(ctx, "foo", dest)
	if err != nil {
		t.Fatal(err)
	}
	im := new(imageserver.Image)
	err = im.UnmarshalBinary(data)
	if err != nil {
		t.Fatal(err)
	}
	diff := compare.Compare(im, testdata.Medium)
	if len(diff) != 0 {
		t.Fatalf("images not equal, diff:\n%+v", diff)
	}
}

func TestGetterErrorContextType(t *testing.T) {
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
	}
	err := gt.Get("invalid", "foo", dest)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetterErrorContextNil(t *testing.T) {
	var ctx *Context
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
	}
	err := gt.Get(ctx, "foo", dest)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetterErrorParamsNil(t *testing.T) {
	ctx := &Context{
		Params: nil,
	}
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
	}
	err := gt.Get(ctx, "foo", dest)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetterErrorServer(t *testing.T) {
	ctx := &Context{
		Params: imageserver.Params{},
	}
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return nil, fmt.Errorf("error")
		}),
	}
	err := gt.Get(ctx, "foo", dest)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetterErrorImageMarshal(t *testing.T) {
	ctx := &Context{
		Params: imageserver.Params{},
	}
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return &imageserver.Image{
				Format: strings.Repeat("a", imageserver.ImageFormatMaxLen+1),
			}, nil
		}),
	}
	err := gt.Get(ctx, "foo", dest)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetterErrorSink(t *testing.T) {
	ctx := &Context{
		Params: imageserver.Params{},
	}
	dest := groupcache.AllocatingByteSliceSink(nil)
	gt := &Getter{
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
	}
	err := gt.Get(ctx, "foo", dest)
	if err == nil {
		t.Fatal("no error")
	}
}
