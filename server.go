package main

import (
	"time"
	"io/ioutil"
	"net/http"
	"github.com/bradfitz/gomemcache/memcache"
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

