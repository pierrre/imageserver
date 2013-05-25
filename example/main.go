package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/pierrre/imageproxy"
	"net/http"
)

func main() {
	server := &imageproxy.Server{
		HttpServer: &http.Server{
			Addr: ":8080",
		},
		RequestParser: &imageproxy.SimpleRequestParser{},
		Cache: &imageproxy.MemcacheCache{
			Prefix:   "imageproxy",
			Memcache: memcache.New("localhost:11211"),
		},
		Converter: &imageproxy.GraphicsMagickConverter{
			Executable: "/usr/local/bin/gm",
		},
	}

	server.Run()
}
