package main

import (
	imageproxy ".."
	"github.com/bradfitz/gomemcache/memcache"
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
	}

	server.Run()
}
