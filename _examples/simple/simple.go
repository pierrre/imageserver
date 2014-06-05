package main

import (
	"crypto/sha256"
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_graphicsmagick "github.com/pierrre/imageserver/http/parser/graphicsmagick"
	imageserver_http_parser_list "github.com/pierrre/imageserver/http/parser/list"
	imageserver_http_parser_source "github.com/pierrre/imageserver/http/parser/source"
	imageserver_processor_graphicsmagick "github.com/pierrre/imageserver/processor/graphicsmagick"
	imageserver_provider_http "github.com/pierrre/imageserver/provider/http"
)

func main() {
	cache := imageserver_cache_memory.New(10 * 1024 * 1024)

	imageServer := &imageserver.ImageServer{
		Cache:        cache,
		CacheKeyFunc: imageserver.NewParametersHashCacheKeyFunc(sha256.New),
		Provider:     &imageserver_provider_http.HTTPProvider{},
		Processor: &imageserver_processor_graphicsmagick.GraphicsMagickProcessor{
			Executable: "gm",
		},
	}

	imageHTTPHandler := &imageserver_http.ImageHTTPHandler{
		Parser: &imageserver_http_parser_list.ListParser{
			&imageserver_http_parser_source.SourceParser{},
			&imageserver_http_parser_graphicsmagick.GraphicsMagickParser{},
		},
		ImageServer: imageServer,
	}

	http.Handle("/", imageHTTPHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
