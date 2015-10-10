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
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_nfntresize "github.com/pierrre/imageserver/http/nfntresize"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/bmp"
	imageserver_image_gamma "github.com/pierrre/imageserver/image/gamma"
	_ "github.com/pierrre/imageserver/image/gif"
	_ "github.com/pierrre/imageserver/image/jpeg"
	imageserver_image_nfntresize "github.com/pierrre/imageserver/image/nfntresize"
	_ "github.com/pierrre/imageserver/image/png"
	_ "github.com/pierrre/imageserver/image/tiff"
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
	http.Handle("/favicon.ico", http.NotFoundHandler())
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
	var handler http.Handler
	handler = &imageserver_http.Handler{
		Parser: &imageserver_http.ListParser{
			&imageserver_http.SourceParser{},
			&imageserver_http_nfntresize.Parser{},
			&imageserver_http.FormatParser{},
			&imageserver_http.QualityParser{},
			&imageserver_http.GammaCorrectionParser{},
		},
		Server:   newServer(),
		ETagFunc: imageserver_http.NewParamsHashETagFunc(sha256.New),
		ErrorFunc: func(err error, request *http.Request) {
			log.Printf("Internal error: %s", err)
		},
	}
	if flagHTTPExpires != 0 {
		handler = &imageserver_http.ExpiresHandler{
			Handler: handler,
			Expires: flagHTTPExpires,
		}
	}
	return handler
}

func newServer() imageserver.Server {
	server := imageserver_testdata.Server
	server = newServerImage(server)
	server = newServerLimit(server)
	server = newServerValidate(server)
	server = newServerCacheRedis(server)
	server = newServerCacheMemory(server)
	return server
}

func newServerImage(server imageserver.Server) imageserver.Server {
	return &imageserver_image.Server{
		Server: server,
		Processor: imageserver_image_gamma.NewCorrectionProcessor(
			&imageserver_image_nfntresize.Processor{},
			true,
		),
	}
}

func newServerLimit(server imageserver.Server) imageserver.Server {
	return imageserver.NewLimitServer(server, runtime.GOMAXPROCS(0)*2)
}

func newServerValidate(server imageserver.Server) imageserver.Server {
	return &imageserver_image_nfntresize.ValidateParamsServer{
		Server:    server,
		WidthMax:  2048,
		HeightMax: 2048,
	}
}

func newServerCacheRedis(server imageserver.Server) imageserver.Server {
	if flagRedis == "" {
		return server
	}
	var cache imageserver_cache.Cache
	cache = &imageserver_cache_redis.Cache{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", flagRedis)
			},
			MaxIdle: 50,
		},
		Expire: flagRedisExpire,
	}
	cache = &imageserver_cache.IgnoreError{
		Cache: cache,
	}
	cache = &imageserver_cache.Async{
		Cache: cache,
	}
	return &imageserver_cache.Server{
		Server:       server,
		Cache:        cache,
		KeyGenerator: imageserver_cache.NewParamsHashKeyGenerator(sha256.New),
	}
}

func newServerCacheMemory(server imageserver.Server) imageserver.Server {
	if flagCacheMemory <= 0 {
		return server
	}
	return &imageserver_cache.Server{
		Server:       server,
		Cache:        imageserver_cache_memory.New(flagCacheMemory),
		KeyGenerator: imageserver_cache.NewParamsHashKeyGenerator(sha256.New),
	}
}
