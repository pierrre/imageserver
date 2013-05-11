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
		Cache: &imageproxy.MemcacheCache{
			Prefix:   "imageproxy",
			Memcache: memcache.New("localhost"),
		},
		Converter: &imageproxy.GraphicsMagickConverter{
			Executable: "/usr/local/bin/gm",
		},
	}

	server.Run()
}
