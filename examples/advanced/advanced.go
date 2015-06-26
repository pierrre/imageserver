package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/githubhook"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_cache_redis "github.com/pierrre/imageserver/cache/redis"
	imageserver_graphicsmagick "github.com/pierrre/imageserver/graphicsmagick"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_graphicsmagick "github.com/pierrre/imageserver/http/graphicsmagick"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

var (
	flagHTTPAddr            = ":8080"
	flagGitHubWebhookSecret string
	flagHTTPExpires         = time.Duration(7 * 24 * time.Hour)
	flagRedis               = "localhost:6379"
	flagRedisExpire         = time.Duration(7 * 24 * time.Hour)
	flagCacheMemory         = int64(64 * (1 << 20))
)

func main() {
	parseFlags()
	log.Println("Start")
	logEnv()
	startHTTPServer()
}

func parseFlags() {
	flag.StringVar(&flagHTTPAddr, "http", flagHTTPAddr, "HTTP addr")
	flag.DurationVar(&flagHTTPExpires, "http-expires", flagHTTPExpires, "HTTP expires")
	flag.StringVar(&flagGitHubWebhookSecret, "github-webhook-secret", "", "GitHub webhook secret")
	flag.StringVar(&flagRedis, "redis", flagRedis, "Redis addr")
	flag.DurationVar(&flagRedisExpire, "redis-expire", flagRedisExpire, "Redis expire")
	flag.Int64Var(&flagCacheMemory, "cache-memory", flagCacheMemory, "Cache memory")
	flag.Parse()
}

func logEnv() {
	log.Printf("Go version: %s", runtime.Version())
	log.Printf("Go max procs: %d", runtime.GOMAXPROCS(0))
}

func startHTTPServer() {
	http.Handle("/", newImageHTTPHandler())
	if h := newGitHubWebhookHTTPHandler(); h != nil {
		http.Handle("/github_webhook", h)
	}
	log.Printf("Start HTTP server on %s", flagHTTPAddr)
	err := http.ListenAndServe(flagHTTPAddr, nil)
	if err != nil {
		log.Panic(err)
	}
}

func newGitHubWebhookHTTPHandler() http.Handler {
	if flagGitHubWebhookSecret == "" {
		return nil
	}
	return &githubhook.Handler{
		Secret: flagGitHubWebhookSecret,
		Delivery: func(event string, deliveryID string, payload interface{}) {
			log.Printf("Received GitHub webhook: %s", event)
			if event == "push" {
				delay := time.Duration(5 * time.Second)
				log.Printf("Killing process in %s", delay)
				time.AfterFunc(delay, func() {
					log.Println("Killing process now")
					os.Exit(0)
				})
			}
		},
		Error: func(err error, req *http.Request) {
			log.Printf("GitHub webhook error: %s", err)
		},
	}
}

func newImageHTTPHandler() http.Handler {
	handler := http.Handler(&imageserver_http.Handler{
		Parser:   newParser(),
		Server:   newServer(),
		ETagFunc: imageserver_http.NewParamsHashETagFunc(sha256.New),
		ErrorFunc: func(err error, request *http.Request) {
			log.Println("Error:", err)
		},
	})
	if flagHTTPExpires != 0 {
		handler = &imageserver_http.ExpiresHandler{
			Handler: handler,
			Expires: flagHTTPExpires,
		}
	}
	return handler
}

func newParser() imageserver_http.Parser {
	return &imageserver_http.ListParser{
		&imageserver_http.SourceParser{},
		&imageserver_http_graphicsmagick.Parser{},
	}
}

func newServer() imageserver.Server {
	server := newServerTestData()
	server = newServerGraphicsMagick(server)
	server = newServerLimit(server)
	server = newServerCache(server)
	return server
}

func newServerTestData() imageserver.Server {
	return imageserver_testdata.Server
}

func newServerGraphicsMagick(server imageserver.Server) imageserver.Server {
	return &imageserver_graphicsmagick.Server{
		Server:     server,
		Executable: "gm",
		Timeout:    time.Duration(10 * time.Second),
		AllowedFormats: []string{
			"jpeg",
			"png",
			"bmp",
			"gif",
		},
	}
}

func newServerLimit(server imageserver.Server) imageserver.Server {
	return imageserver.NewLimitServer(server, runtime.GOMAXPROCS(0)*2)
}

func newServerCache(server imageserver.Server) imageserver.Server {
	keyGenerator := imageserver_cache.NewParamsHashKeyGenerator(sha256.New)
	if cache := newCacheRedis(); cache != nil {
		server = &imageserver_cache.Server{
			Server:       server,
			Cache:        cache,
			KeyGenerator: keyGenerator,
		}
	}
	if cache := newCacheMemory(); cache != nil {
		server = &imageserver_cache.Server{
			Server:       server,
			Cache:        cache,
			KeyGenerator: keyGenerator,
		}
	}
	return server
}

func newCacheRedis() imageserver_cache.Cache {
	if flagRedis == "" {
		return nil
	}
	cache := imageserver_cache.Cache(&imageserver_cache_redis.Cache{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", flagRedis)
			},
			MaxIdle: 50,
		},
		Expire: flagRedisExpire,
	})
	cache = &imageserver_cache.IgnoreError{
		Cache: cache,
	}
	cache = &imageserver_cache.Async{
		Cache: cache,
	}
	return cache
}

func newCacheMemory() imageserver_cache.Cache {
	if flagCacheMemory <= 0 {
		return nil
	}
	return imageserver_cache_memory.New(flagCacheMemory)
}
