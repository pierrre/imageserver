// Package groupcache provides a groupcache imageserver.Server implementation.
package groupcache

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/groupcache"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
)

// NewServer is a helper to create a new groupcache Server.
func NewServer(srv imageserver.Server, kg imageserver_cache.KeyGenerator, name string, cacheBytes int64) *Server {
	return &Server{
		Group:        groupcache.NewGroup(name, cacheBytes, &Getter{Server: srv}),
		KeyGenerator: kg,
	}
}

// Server is a groupcache imageserver.Server implementation.
//
// Group MUST use a Getter from this package.
type Server struct {
	Group        *groupcache.Group
	KeyGenerator imageserver_cache.KeyGenerator
}

// Get implements imageserver.Server.
func (srv *Server) Get(goctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
	myctx := &Context{
		Params: params,
	}
	key := srv.KeyGenerator.GetKey(params)
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	err := srv.Group.Get(myctx, key, dest)
	if err != nil {
		return nil, err
	}
	im := new(imageserver.Image)
	err = im.UnmarshalBinaryNoCopy(data)
	if err != nil {
		return nil, err
	}
	return im, nil
}

// Getter is a groupcache.Getter implementation for Server.
type Getter struct {
	Server imageserver.Server
}

// Get implements groupcache.Getter.
func (gt *Getter) Get(gcctx groupcache.Context, key string, dest groupcache.Sink) error {
	myctx, err := gt.getContext(gcctx)
	if err != nil {
		return err
	}
	goctx, cancel := gt.getGoContext(myctx)
	if cancel != nil {
		defer cancel()
	}
	return gt.get(goctx, myctx.Params, dest)
}

func (gt *Getter) getContext(gcctx groupcache.Context) (*Context, error) {
	myctx, ok := gcctx.(*Context)
	if !ok {
		return nil, fmt.Errorf("invalid context type: %T", gcctx)
	}
	if myctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	if myctx.Params == nil {
		return nil, fmt.Errorf("context has nil Params")
	}
	return myctx, nil
}

func (gt *Getter) getGoContext(myctx *Context) (goctx context.Context, cancel func()) {
	goctx = context.Background()
	if myctx.GoContext.Deadline != nil {
		goctx, cancel = context.WithDeadline(goctx, *myctx.GoContext.Deadline)
	}
	for k, v := range myctx.GoContext.Values {
		goctx = context.WithValue(goctx, k, v)
	}
	return goctx, cancel
}

func (gt *Getter) get(goctx context.Context, params imageserver.Params, dest groupcache.Sink) error {
	im, err := gt.Server.Get(goctx, params)
	if err != nil {
		return err
	}
	data, err := im.MarshalBinary()
	if err != nil {
		return err
	}
	return dest.SetBytes(data)
}

// Context is a groupcache.Context implementation used by Getter.
type Context struct {
	Params    imageserver.Params
	GoContext struct {
		Deadline *time.Time
		Values   map[interface{}]interface{}
	}
}
