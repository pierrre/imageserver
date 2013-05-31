package main

import (
	"github.com/pierrre/imageproxy"
	imageproxy_cache_memory "github.com/pierrre/imageproxy/cache/memory"
	imageproxy_cache_prefix "github.com/pierrre/imageproxy/cache/prefix"
	imageproxy_converter_graphicsmagick "github.com/pierrre/imageproxy/converter/graphicsmagick"
	imageproxy_requestparser_simple "github.com/pierrre/imageproxy/requestparser/simple"
	"net/http"
)

func main() {
	cache := imageproxy_cache_memory.New(10 * 1024 * 1024)

	server := &imageproxy.Server{
		HttpServer: &http.Server{
			Addr: ":8080",
		},
		RequestParser: &imageproxy_requestparser_simple.SimpleRequestParser{},
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
		},
	}
	server.Run()
}
