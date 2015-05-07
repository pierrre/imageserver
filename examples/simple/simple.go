package main

import (
	"net/http"

	imageserver_graphicsmagick "github.com/pierrre/imageserver/graphicsmagick"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_graphicsmagick "github.com/pierrre/imageserver/http/graphicsmagick"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

func main() {
	server := imageserver_testdata.Server
	server = &imageserver_graphicsmagick.Server{
		Server:     server,
		Executable: "gm",
	}
	handler := &imageserver_http.Handler{
		Parser: &imageserver_http.ListParser{
			&imageserver_http.SourceParser{},
			&imageserver_http_graphicsmagick.Parser{},
		},
		Server: server,
	}
	http.Handle("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
