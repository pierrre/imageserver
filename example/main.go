package main

import (
	imageproxy ".."
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
)

func main() {
	server := imageproxy.NewServer(
		&http.Server{
			Addr: ":8080",
		},
		imageproxy.NewMemcacheCache(
			"imageproxy",
			memcache.New("localhost"),
		),
	)
	server.Run();
}
