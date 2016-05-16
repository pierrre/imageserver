// Package advanced provides a groupcache example.
//
// Run:
//  go run groupcache.go -http=:8081 -groupcache-peers=:8082,:8083
//  go run groupcache.go -http=:8082 -groupcache-peers=:8081,:8083
//  go run groupcache.go -http=:8083 -groupcache-peers=:8081,:8082
//
// Open http://localhost:8081/medium.jpg?width=100
//
// Stats are available on http://localhost:8081/stats
package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/disintegration/gift"
	"github.com/golang/groupcache"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_cache_groupcache "github.com/pierrre/imageserver/cache/groupcache"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_gift "github.com/pierrre/imageserver/http/gift"
	imageserver_http_image "github.com/pierrre/imageserver/http/image"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/gif"
	imageserver_image_gift "github.com/pierrre/imageserver/image/gift"
	_ "github.com/pierrre/imageserver/image/jpeg"
	_ "github.com/pierrre/imageserver/image/png"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

const (
	groupcacheName = "imageserver"
)

var (
	flagHTTP            = ":8080"
	flagGroupcache      = int64(128 * (1 << 20))
	flagGroupcachePeers string
)

func main() {
	parseFlags()
	startHTTPServer()
}

func parseFlags() {
	flag.StringVar(&flagHTTP, "http", flagHTTP, "HTTP")
	flag.Int64Var(&flagGroupcache, "groupcache", flagGroupcache, "Groupcache")
	flag.StringVar(&flagGroupcachePeers, "groupcache-peers", flagGroupcachePeers, "Groupcache peers")
	flag.Parse()
}

func startHTTPServer() {
	http.Handle("/", http.StripPrefix("/", newImageHTTPHandler()))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	initGroupcacheHTTPPool() // it automatically registers itself to "/_groupcache"
	http.HandleFunc("/stats", groupcacheStatsHTTPHandler)
	err := http.ListenAndServe(flagHTTP, nil)
	if err != nil {
		panic(err)
	}
}

func newImageHTTPHandler() http.Handler {
	return &imageserver_http.Handler{
		Parser: imageserver_http.ListParser([]imageserver_http.Parser{
			&imageserver_http.SourcePathParser{},
			&imageserver_http_gift.ResizeParser{},
			&imageserver_http_image.FormatParser{},
			&imageserver_http_image.QualityParser{},
		}),
		Server: newServer(),
	}
}

func newServer() imageserver.Server {
	srv := imageserver_testdata.Server
	srv = newServerImage(srv)
	srv = newServerGroupcache(srv)
	return srv
}

func newServerImage(srv imageserver.Server) imageserver.Server {
	return &imageserver.HandlerServer{
		Server: srv,
		Handler: &imageserver_image.Handler{
			Processor: &imageserver_image_gift.ResizeProcessor{
				DefaultResampling: gift.LanczosResampling,
			},
		},
	}
}

func newServerGroupcache(srv imageserver.Server) imageserver.Server {
	if flagGroupcache <= 0 {
		return srv
	}
	return imageserver_cache_groupcache.NewServer(
		srv,
		imageserver_cache.NewParamsHashKeyGenerator(sha256.New),
		groupcacheName,
		flagGroupcache,
	)
}

func initGroupcacheHTTPPool() {
	self := (&url.URL{Scheme: "http", Host: flagHTTP}).String()
	var peers []string
	peers = append(peers, self)
	for _, p := range strings.Split(flagGroupcachePeers, ",") {
		if p == "" {
			continue
		}
		peer := (&url.URL{Scheme: "http", Host: p}).String()
		peers = append(peers, peer)
	}
	pool := groupcache.NewHTTPPool(self)
	pool.Context = imageserver_cache_groupcache.HTTPPoolContext
	pool.Transport = imageserver_cache_groupcache.NewHTTPPoolTransport(nil)
	pool.Set(peers...)
}

func groupcacheStatsHTTPHandler(w http.ResponseWriter, req *http.Request) {
	gp := groupcache.GetGroup(groupcacheName)
	if gp == nil {
		http.Error(w, fmt.Sprintf("group %s not found", groupcacheName), http.StatusServiceUnavailable)
		return
	}
	type cachesStats struct {
		Main groupcache.CacheStats
		Hot  groupcache.CacheStats
	}
	type stats struct {
		Group  groupcache.Stats
		Caches cachesStats
	}
	data, err := json.MarshalIndent(
		stats{
			Group: gp.Stats,
			Caches: cachesStats{
				Main: gp.CacheStats(groupcache.MainCache),
				Hot:  gp.CacheStats(groupcache.HotCache),
			},
		},
		"",
		"	",
	)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
