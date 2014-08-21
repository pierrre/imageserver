package main

import (
	"crypto/sha256"
	//_ "expvar"
	"flag"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"strconv"
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
	var verbose bool
	var httpAddr string
	flag.BoolVar(&verbose, "verbose", false, "Verbose")
	flag.StringVar(&httpAddr, "http", ":8080", "Http")
	flag.Parse()

	log.Println("Start")

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	var cache imageserver_cache.Cache
	cache = &imageserver_cache_redis.RedisCache{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", "localhost:6379")
			},
			MaxIdle: 50,
		},
		Expire: time.Duration(7 * 24 * time.Hour),
	}
	cache = &imageserver_cache_async.AsyncCache{
		Cache: cache,
		ErrFunc: func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters) {
			if verbose {
				log.Println("Cache error:", err)
			}
		},
	}
	cache = imageserver_cache_list.ListCache{
		imageserver_cache_memory.New(10 * 1024 * 1024),
		cache,
	}

	provider := &imageserver_provider_cache.CacheProvider{
		Provider:     &imageserver_provider_http.HTTPProvider{},
		Cache:        cache,
		KeyGenerator: imageserver_provider_cache.NewSourceHashKeyGenerator(sha256.New),
	}

	var processor imageserver_processor.Processor
	processor = &imageserver_processor_graphicsmagick.GraphicsMagickProcessor{
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

	var imageServer imageserver.ImageServer
	imageServer = &imageserver_provider.ProviderImageServer{
		Provider: provider,
	}
	imageServer = &imageserver_processor.ProcessorImageServer{
		ImageServer: imageServer,
		Processor:   processor,
	}
	imageServer = &imageserver_cache.CacheImageServer{
		ImageServer:  imageServer,
		Cache:        cache,
		KeyGenerator: imageserver_cache.NewParametersHashKeyGenerator(sha256.New),
	}

	imageHTTPHandler := &imageserver_http.ImageHTTPHandler{
		Parser: &imageserver_http_parser_list.ListParser{
			&imageserver_http_parser_source.SourceParser{},
			&imageserver_http_parser_graphicsmagick.GraphicsMagickParser{},
		},
		ImageServer: imageServer,
		ETagFunc:    imageserver_http.NewParametersHashETagFunc(sha256.New),
		Expire:      time.Duration(7 * 24 * time.Hour),
		RequestFunc: func(request *http.Request) error {
			if verbose {
				log.Println("Request:", strconv.Quote(request.URL.String()))
			}
			return nil
		},
		HeaderFunc: func(header http.Header, request *http.Request, err error) {
			header.Set("X-Hostname", hostname)
		},
		ErrorFunc: func(err error, request *http.Request) {
			if verbose {
				log.Println("Error:", err)
			}
		},
		ResponseFunc: func(request *http.Request, statusCode int, contentSize int64, err error) {
			if verbose {
				var errString string
				if err != nil {
					errString = err.Error()
				}
				log.Println("Response:", request.RemoteAddr, request.Method, strconv.Quote(request.URL.String()), statusCode, contentSize, strconv.Quote(errString))
			}
		},
	}
	http.Handle("/", imageHTTPHandler)

	err = http.ListenAndServe(httpAddr, nil)
	if err != nil {
		log.Panic(err)
	}
}
