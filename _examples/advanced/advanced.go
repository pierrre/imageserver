package main

import (
	"crypto/sha256"
	//_ "expvar"
	"flag"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_cache_async "github.com/pierrre/imageserver/cache/async"
	imageserver_cache_list "github.com/pierrre/imageserver/cache/list"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_cache_redis "github.com/pierrre/imageserver/cache/redis"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_graphicsmagick "github.com/pierrre/imageserver/http/parser/graphicsmagick"
	imageserver_http_parser_list "github.com/pierrre/imageserver/http/parser/list"
	imageserver_http_parser_source "github.com/pierrre/imageserver/http/parser/source"
	imageserver_processor "github.com/pierrre/imageserver/processor"
	imageserver_processor_graphicsmagick "github.com/pierrre/imageserver/processor/graphicsmagick"
	imageserver_processor_limit "github.com/pierrre/imageserver/processor/limit"
	imageserver_provider "github.com/pierrre/imageserver/provider"
	imageserver_provider_cache "github.com/pierrre/imageserver/provider/cache"
	imageserver_provider_http "github.com/pierrre/imageserver/provider/http"
)

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", ":8080", "Http")
	flag.Parse()

	log.Println("Start")

	var cache imageserver_cache.Cache
	cache = &imageserver_cache_redis.Cache{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", "localhost:6379")
			},
			MaxIdle: 50,
		},
		Expire: time.Duration(7 * 24 * time.Hour),
	}
	cache = &imageserver_cache_async.Cache{
		Cache: cache,
		ErrFunc: func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters) {
			log.Println("Cache error:", err)
		},
	}
	cache = imageserver_cache_list.Cache{
		imageserver_cache_memory.New(10 * 1024 * 1024),
		cache,
	}

	provider := &imageserver_provider_cache.Provider{
		Provider:     &imageserver_provider_http.Provider{},
		Cache:        cache,
		KeyGenerator: imageserver_provider_cache.NewSourceHashKeyGenerator(sha256.New),
	}

	var processor imageserver_processor.Processor
	processor = &imageserver_processor_graphicsmagick.Processor{
		Executable: "gm",
		Timeout:    time.Duration(10 * time.Second),
		AllowedFormats: []string{
			"jpeg",
			"png",
			"bmp",
			"gif",
		},
	}
	processor = imageserver_processor_limit.New(processor, 16)

	var server imageserver.Server
	server = &imageserver_provider.Server{
		Provider: provider,
	}
	server = &imageserver_processor.Server{
		Server:    server,
		Processor: processor,
	}
	server = &imageserver_cache.Server{
		Server:       server,
		Cache:        cache,
		KeyGenerator: imageserver_cache.NewParametersHashKeyGenerator(sha256.New),
	}

	var handler http.Handler
	handler = &imageserver_http.Handler{
		Parser: &imageserver_http_parser_list.Parser{
			&imageserver_http_parser_source.Parser{},
			&imageserver_http_parser_graphicsmagick.Parser{},
		},
		Server:   server,
		ETagFunc: imageserver_http.NewParametersHashETagFunc(sha256.New),
		ErrorFunc: func(err error, request *http.Request) {
			log.Println("Error:", err)
		},
	}
	handler = &imageserver_http.ExpiresHandler{
		Handler: handler,
		Expires: time.Duration(7 * 24 * time.Hour),
	}
	http.Handle("/", handler)

	err := http.ListenAndServe(httpAddr, nil)
	if err != nil {
		log.Panic(err)
	}
}
