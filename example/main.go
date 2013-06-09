package main

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
	imageserver_cache_chain "github.com/pierrre/imageserver/cache/chain"
	imageserver_cache_memcache "github.com/pierrre/imageserver/cache/memcache"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_cache_prefix "github.com/pierrre/imageserver/cache/prefix"
	imageserver_converter_graphicsmagick "github.com/pierrre/imageserver/converter/graphicsmagick"
	imageserver_requestparser_graphicsmagick "github.com/pierrre/imageserver/requestparser/graphicsmagick"
	imageserver_requestparser_merge "github.com/pierrre/imageserver/requestparser/merge"
	imageserver_requestparser_source "github.com/pierrre/imageserver/requestparser/source"
	"net/http"
)

func main() {
	cache := &imageserver_cache_chain.ChainCache{
		imageserver_cache_memory.New(10 * 1024 * 1024),
		&imageserver_cache_memcache.MemcacheCache{
			Memcache: memcache_impl.New("localhost:11211"),
		},
	}

	server := &imageserver.Server{
		HttpServer: &http.Server{
			Addr: ":8080",
		},
		RequestParser: &imageserver_requestparser_merge.MergeRequestParser{
			&imageserver_requestparser_source.SourceRequestParser{},
			&imageserver_requestparser_graphicsmagick.GraphicsMagickRequestParser{},
		},
		Cache: &imageserver_cache_prefix.PrefixCache{
			Prefix: "converted_",
			Cache:  cache,
		},
		SourceCache: &imageserver_cache_prefix.PrefixCache{
			Prefix: "source_",
			Cache:  cache,
		},
		Converter: &imageserver_converter_graphicsmagick.GraphicsMagickConverter{
			Executable: "/usr/local/bin/gm",
			AllowedFormats: []string{
				"jpeg",
				"png",
				"bmp",
			},
			DefaultQualities: map[string]string{
				"jpeg": "85",
			},
		},
	}
	server.Run()
}
