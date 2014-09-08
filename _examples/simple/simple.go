package main

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_graphicsmagick "github.com/pierrre/imageserver/http/parser/graphicsmagick"
	imageserver_http_parser_list "github.com/pierrre/imageserver/http/parser/list"
	imageserver_http_parser_source "github.com/pierrre/imageserver/http/parser/source"
	imageserver_processor "github.com/pierrre/imageserver/processor"
	imageserver_processor_graphicsmagick "github.com/pierrre/imageserver/processor/graphicsmagick"
	imageserver_provider "github.com/pierrre/imageserver/provider"
	imageserver_provider_http "github.com/pierrre/imageserver/provider/http"
)

func main() {
	var server imageserver.Server
	server = &imageserver_provider.Server{
		Provider: &imageserver_provider_http.Provider{},
	}
	server = &imageserver_processor.Server{
		Server: server,
		Processor: &imageserver_processor_graphicsmagick.Processor{
			Executable: "gm",
		},
	}

	handler := &imageserver_http.Handler{
		Parser: &imageserver_http_parser_list.Parser{
			&imageserver_http_parser_source.Parser{},
			&imageserver_http_parser_graphicsmagick.Parser{},
		},
		Server: server,
	}

	http.Handle("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
