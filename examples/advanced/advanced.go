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
	"net"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"time"
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

	var cache imageserver.Cache
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
		CacheKeyFunc: imageserver_provider_cache.NewSourceHashCacheKeyFunc(sha256.New),
	}

	var processor imageserver.Processor
	processor = &imageserver_processor_graphicsmagick.GraphicsMagickProcessor{
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
	}
	processor = imageserver_processor_limit.New(processor, 16)

	imageServer := &imageserver.Server{
		Cache:        cache,
		CacheKeyFunc: imageserver.NewParametersHashCacheKeyFunc(sha256.New),
		Provider:     provider,
		Processor:    processor,
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
	http.Handle("/", httpImageServer)

	tcpAddr, err := net.ResolveTCPAddr("tcp", httpAddr)
	if err != nil {
		panic(err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	interrupted := false
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		<-interruptChan
		interrupted = true
		log.Println("Close TCP listener")
		err := tcpListener.Close()
		if err != nil {
			panic(err)
		}
	}()

	log.Println("Start HTTP server")
	err = http.Serve(tcpListener, nil)
	if err != nil {
		if interrupted {
			waitDuration := 10 * time.Second
			log.Printf("Wait clients %s (press CTRL+C again to stop the server immediatly)", waitDuration)
			select {
			case <-time.After(waitDuration):
			case <-interruptChan:
			}
		} else {
			panic(err)
		}
	}

	log.Println("Exit")
}
