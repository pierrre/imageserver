package main

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageproxy"
	imageproxy_cache_chain "github.com/pierrre/imageproxy/cache/chain"
	imageproxy_cache_memcache "github.com/pierrre/imageproxy/cache/memcache"
	imageproxy_cache_memory "github.com/pierrre/imageproxy/cache/memory"
	imageproxy_cache_prefix "github.com/pierrre/imageproxy/cache/prefix"
	imageproxy_converter_graphicsmagick "github.com/pierrre/imageproxy/converter/graphicsmagick"
	imageproxy_requestparser_graphicsmagick "github.com/pierrre/imageproxy/requestparser/graphicsmagick"
	imageproxy_requestparser_merge "github.com/pierrre/imageproxy/requestparser/merge"
	imageproxy_requestparser_source "github.com/pierrre/imageproxy/requestparser/source"
	"net/http"
)

func main() {
	cache := &imageproxy_cache_chain.ChainCache{
		imageproxy_cache_memory.New(10 * 1024 * 1024),
		&imageproxy_cache_memcache.MemcacheCache{
			Memcache: memcache_impl.New("localhost:11211"),
		},
	}

	server := &imageproxy.Server{
		HttpServer: &http.Server{
			Addr: ":8080",
		},
		RequestParser: &imageproxy_requestparser_merge.MergeRequestParser{
			&imageproxy_requestparser_source.SourceRequestParser{},
			&imageproxy_requestparser_graphicsmagick.GraphicsMagickRequestParser{},
		},
		Cache: &imageproxy_cache_prefix.PrefixCache{
			Prefix: "converted_",
			Cache:  cache,
		},
		SourceCache: &imageproxy_cache_prefix.PrefixCache{
			Prefix: "source_",
			Cache:  cache,
		},
		Converter: &imageproxy_converter_graphicsmagick.GraphicsMagickConverter{
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
