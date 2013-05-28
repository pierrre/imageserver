package main

import (
	//"github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageproxy"
	//imageproxy_cache_memcache "github.com/pierrre/imageproxy/cache/memcache"
	imageproxy_cache_memory "github.com/pierrre/imageproxy/cache/memory"
	imageproxy_converter_graphicsmagick "github.com/pierrre/imageproxy/converter/graphicsmagick"
	"net/http"
)

func main() {
	server := &imageproxy.Server{
		HttpServer: &http.Server{
			Addr: ":8080",
		},
		RequestParser: &imageproxy.SimpleRequestParser{},
		/*
			Cache: &imageproxy_cache_memcache.MemcacheCache{
				Prefix:   "imageproxy",
				Memcache: memcache.New("localhost:11211"),
			},
		*/
		Cache: imageproxy_cache_memory.New(10 * 1024 * 1024),
		Converter: &imageproxy_converter_graphicsmagick.GraphicsMagickConverter{
			Executable: "/usr/local/bin/gm",
		},
	}

	server.Run()
}
