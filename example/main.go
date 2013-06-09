package main

import (
	memcache_impl "github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageserver"
	imageserver_cache_chain "github.com/pierrre/imageserver/cache/chain"
	imageserver_cache_memcache "github.com/pierrre/imageserver/cache/memcache"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_cache_prefix "github.com/pierrre/imageserver/cache/prefix"
	imageserver_converter_graphicsmagick "github.com/pierrre/imageserver/converter/graphicsmagick"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_graphicsmagick "github.com/pierrre/imageserver/http/parser/graphicsmagick"
	imageserver_http_parser_merge "github.com/pierrre/imageserver/http/parser/merge"
	imageserver_http_parser_source "github.com/pierrre/imageserver/http/parser/source"
	"net/http"
)

func main() {
	cache := imageserver_cache_chain.ChainCache{
		imageserver_cache_memory.New(10 * 1024 * 1024),
		&imageserver_cache_memcache.MemcacheCache{
			Memcache: memcache_impl.New("localhost:11211"),
		},
	}

	server := imageserver_http.Server{
		HttpServer: http.Server{
			Addr: ":8080",
		},
		Parser: imageserver_http_parser_merge.MergeParser{
			&imageserver_http_parser_source.SourceParser{},
			&imageserver_http_parser_graphicsmagick.GraphicsMagickParser{},
		},
		ImageServer: imageserver.Server{
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
					"gif",
				},
				DefaultQualities: map[string]string{
					"jpeg": "85",
				},
			},
		},
	}
	server.Run()
}
