package main

import (
	imageproxy ".."
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
)

func main() {
	server := imageproxy.NewServer(
		imageproxy.NewMemcacheCache(
			"imageproxy",
			memcache.New("localhost"),
		),
	)

	http.HandleFunc("/", server.HandleHttpRequest)
	http.ListenAndServe(":8080", nil)
}
