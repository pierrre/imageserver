package main

import (
	"net/http"

	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_parser_graphicsmagick "github.com/pierrre/imageserver/http/parser/graphicsmagick"
	imageserver_processor "github.com/pierrre/imageserver/processor"
	imageserver_processor_graphicsmagick "github.com/pierrre/imageserver/processor/graphicsmagick"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

func main() {
	server := imageserver_testdata.Server
	server = &imageserver_processor.Server{
		Server: server,
		Processor: &imageserver_processor_graphicsmagick.Processor{
			Executable: "gm",
		},
	}
	handler := &imageserver_http.Handler{
		Parser: &imageserver_http.ListParser{
			&imageserver_http.SourceParser{},
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
