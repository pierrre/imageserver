// Package groupcache provides a groupcache imageserver.Server implementation.
package groupcache

import (
	"fmt"

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
func (srv *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	ctx := &Context{
		Params: params,
	}
	key := srv.KeyGenerator.GetKey(params)
	var data []byte
	dest := groupcache.AllocatingByteSliceSink(&data)
	err := srv.Group.Get(ctx, key, dest)
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
func (gt *Getter) Get(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	myctx, ok := ctx.(*Context)
	if !ok {
		return fmt.Errorf("invalid context type: %T", ctx)
	}
	if myctx == nil {
		return fmt.Errorf("context is nil")
	}
	if myctx.Params == nil {
		return fmt.Errorf("context has nil Params")
	}
	im, err := gt.Server.Get(myctx.Params)
	if err != nil {
		return err
	}
	data, err := im.MarshalBinary()
	if err != nil {
		return err
	}
	err = dest.SetBytes(data)
	if err != nil {
		return err
	}
	return nil
}

// Context is a groupcache.Context implementation used by Getter.
type Context struct {
	Params imageserver.Params
}
