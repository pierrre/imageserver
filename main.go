package main

import (
	"./src/imageproxy"
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
)

func main() {
	server := imageproxy.NewServer(
		imageproxy.NewMemcacheCache(
			"lol",
			memcache.New("localhost"),
		),
	)

	http.HandleFunc("/", server.HandleHttpRequest)
	http.ListenAndServe(":8080", nil)
}
