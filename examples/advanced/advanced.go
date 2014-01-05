// +build ignore

package main

import (
	"crypto/sha256"
	//_ "expvar"
	"flag"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/imageserver"
	imageserver_cache_async "github.com/pierrre/imageserver/cache/async"
	imageserver_cache_list "github.com/pierrre/imageserver/cache/list"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_cache_redis "github.com/pierrre/imageserver/cache/redis"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_graphicsmagick "github.com/pierrre/imageserver/http/parser/graphicsmagick"
	imageserver_http_parser_list "github.com/pierrre/imageserver/http/parser/list"
	imageserver_http_parser_source "github.com/pierrre/imageserver/http/parser/source"
	imageserver_processor_graphicsmagick "github.com/pierrre/imageserver/processor/graphicsmagick"
	imageserver_processor_limit "github.com/pierrre/imageserver/processor/limit"
	imageserver_provider_cache "github.com/pierrre/imageserver/provider/cache"
	imageserver_provider_http "github.com/pierrre/imageserver/provider/http"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"strconv"
	"time"
)

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Verbose")
	flag.Parse()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	cache := imageserver_cache_list.ListCache{
		imageserver_cache_memory.New(10 * 1024 * 1024),
		&imageserver_cache_async.AsyncCache{
			Cache: &imageserver_cache_redis.RedisCache{
				Pool: &redigo.Pool{
					Dial: func() (redigo.Conn, error) {
						return redigo.Dial("tcp", "localhost:6379")
					},
					MaxIdle: 50,
				},
				Expire: time.Duration(7 * 24 * time.Hour),
			},
			ErrFunc: func(err error, key string, image *imageserver.Image, parameters imageserver.Parameters) {
				if verbose {
					log.Println(err)
				}
			},
		},
	}

	imageServer := &imageserver.Server{
		Cache:        cache,
		CacheKeyFunc: imageserver.NewParametersHashCacheKeyFunc(sha256.New),
		Provider: &imageserver_provider_cache.CacheProvider{
			Provider:     &imageserver_provider_http.HTTPProvider{},
			Cache:        cache,
			CacheKeyFunc: imageserver_provider_cache.NewSourceHashCacheKeyFunc(sha256.New),
		},
		Processor: imageserver_processor_limit.New(&imageserver_processor_graphicsmagick.GraphicsMagickProcessor{
			Executable: "gm",
			Timeout:    time.Duration(10 * time.Second),
			AllowedFormats: []string{
				"jpeg",
				"png",
				"bmp",
				"gif",
			},
			DefaultQualities: map[string]string{
				"jpeg": "85",
			},
		}, 16),
	}

	httpImageServer := &imageserver_http.Server{
		Parser: &imageserver_http_parser_list.ListParser{
			&imageserver_http_parser_source.SourceParser{},
			&imageserver_http_parser_graphicsmagick.GraphicsMagickParser{},
		},
		ImageServer: imageServer,
		ETagFunc:    imageserver_http.NewParametersHashETagFunc(sha256.New),
		Expire:      time.Duration(7 * 24 * time.Hour),
		RequestFunc: func(request *http.Request) error {
			url := request.URL
			query := url.Query()
			errorCodeString := query.Get("error")
			if len(errorCodeString) == 0 {
				return nil
			}
			errorCode, err := strconv.Atoi(errorCodeString)
			if err != nil {
				return imageserver.NewError(err.Error())
			}
			return imageserver_http.NewError(errorCode)
		},
		HeaderFunc: func(header http.Header, request *http.Request, err error) {
			header.Set("X-Hostname", hostname)
		},
		ErrorFunc: func(err error, request *http.Request) {
			if verbose {
				log.Println(err)
			}
		},
		ResponseFunc: func(request *http.Request, statusCode int, contentSize int64, err error) {
			if verbose {
				var errString string
				if err != nil {
					errString = err.Error()
				}
				log.Println(request.RemoteAddr, request.Method, strconv.Quote(request.URL.String()), statusCode, contentSize, strconv.Quote(errString))
			}
		},
	}

	http.Handle("/", httpImageServer)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
