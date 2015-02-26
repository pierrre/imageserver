package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"net/http"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_cache_redis "github.com/pierrre/imageserver/cache/redis"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_graphicsmagick "github.com/pierrre/imageserver/http/parser/graphicsmagick"
	imageserver_processor "github.com/pierrre/imageserver/processor"
	imageserver_processor_graphicsmagick "github.com/pierrre/imageserver/processor/graphicsmagick"
	imageserver_provider "github.com/pierrre/imageserver/provider"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", ":8080", "Http")
	flag.Parse()

	log.Println("Start")

	cache := imageserver_cache.Cache(&imageserver_cache_redis.Cache{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", "localhost:6379")
			},
			MaxIdle: 50,
		},
		Expire: time.Duration(7 * 24 * time.Hour),
	})
	cache = &imageserver_cache.Async{
		Cache: cache,
		ErrFunc: func(err error, key string, image *imageserver.Image, params imageserver.Params) {
			log.Println("Cache error:", err)
		},
	}
	cache = imageserver_cache.List{
		imageserver_cache_memory.New(10 * 1024 * 1024),
		cache,
	}

	processor := imageserver_processor.Processor(&imageserver_processor_graphicsmagick.Processor{
		Executable: "gm",
		Timeout:    time.Duration(10 * time.Second),
		AllowedFormats: []string{
			"jpeg",
			"png",
			"bmp",
			"gif",
		},
	})
	processor = imageserver_processor.NewLimit(processor, 16)

	server := imageserver.Server(&imageserver_provider.Server{
		Provider: imageserver_testdata.Provider,
	})
	server = &imageserver_processor.Server{
		Server:    server,
		Processor: processor,
	}
	server = &imageserver_cache.Server{
		Server:       server,
		Cache:        cache,
		KeyGenerator: imageserver_cache.NewParamsHashKeyGenerator(sha256.New),
	}

	handler := http.Handler(&imageserver_http.Handler{
		Parser: &imageserver_http.ListParser{
			&imageserver_http.SourceParser{},
			&imageserver_http_parser_graphicsmagick.Parser{},
		},
		Server:   server,
		ETagFunc: imageserver_http.NewParamsHashETagFunc(sha256.New),
		ErrorFunc: func(err error, request *http.Request) {
			log.Println("Error:", err)
		},
	})
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
