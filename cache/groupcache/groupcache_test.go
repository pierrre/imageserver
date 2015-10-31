package groupcache

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/groupcache"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	"github.com/pierrre/imageserver/testdata"
)

const (
	testSize = 100 * (1 << 20)
)

var _ imageserver.Server = &Server{}

func TestServer(t *testing.T) {
	srv := newTestServer(
		&imageserver.StaticServer{
			Image: testdata.Medium,
		},
		imageserver_cache.StringKeyGenerator("test"),
	)
	im, err := srv.Get(imageserver.Params{
		imageserver.SourceParam: testdata.MediumFileName,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !imageserver.ImageEqual(im, testdata.Medium) {
		t.Fatal("not equal")
	}
}

func TestServerErrorGroup(t *testing.T) {
	srv := newTestServer(
		&imageserver.StaticServer{
			Error: fmt.Errorf("error"),
		},
		imageserver_cache.StringKeyGenerator("test"),
	)
	_, err := srv.Get(imageserver.Params{})
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
				dest.SetBytes(nil)
				return nil
			}),
		),
		KeyGenerator: imageserver_cache.StringKeyGenerator("test"),
	}
	_, err := srv.Get(imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func newTestServer(srv imageserver.Server, kg imageserver_cache.KeyGenerator) *Server {
	return NewServer(srv, kg, newTestServerName(), testSize)
}

func newTestServerName() string {
	return fmt.Sprintf("test_%d", time.Now().UnixNano())
}

var _ groupcache.Getter = &Getter{}

func TestGetter(t *testing.T) {
	ctx := &Context{
		Params: imageserver.Params{},
	}
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
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
	if !imageserver.ImageEqual(im, testdata.Medium) {
		t.Fatal("not equal")
	}
}

func TestGetterErrorContextType(t *testing.T) {
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	gt := &Getter{
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
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
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
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
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
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
		Server: &imageserver.StaticServer{
			Error: fmt.Errorf("error"),
		},
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
		Server: &imageserver.StaticServer{
			Image: &imageserver.Image{
				Format: strings.Repeat("a", imageserver.ImageFormatMaxLen+1),
			},
		},
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
		Server: &imageserver.StaticServer{
			Image: testdata.Medium,
		},
	}
	err := gt.Get(ctx, "foo", dest)
	if err == nil {
		t.Fatal("no error")
	}
}
