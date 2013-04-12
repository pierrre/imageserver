package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	cache = memcache.New("localhost:11211")
)

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	item, err := cache.Get("foo")
	var bytes []byte
	if err == nil {
		bytes = item.Value
	} else {
		time.Sleep(200 * time.Millisecond)
		bytes, _ = ioutil.ReadFile("arbre.webp")
		item = &memcache.Item{Key: "foo", Value: bytes}
		cache.Set(item)
	}
	w.Write(bytes)
}
