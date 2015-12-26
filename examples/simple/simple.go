// Package simple provides a simple example.
package main

import (
	"net/http"

	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_gift "github.com/pierrre/imageserver/http/gift"
	imageserver_image "github.com/pierrre/imageserver/image"
	_ "github.com/pierrre/imageserver/image/gif"
	imageserver_image_gift "github.com/pierrre/imageserver/image/gift"
	_ "github.com/pierrre/imageserver/image/jpeg"
	_ "github.com/pierrre/imageserver/image/png"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

func main() {
	http.Handle("/", &imageserver_http.Handler{
		Parser: imageserver_http.ListParser([]imageserver_http.Parser{
			&imageserver_http.SourceParser{},
			&imageserver_http_gift.Parser{},
			&imageserver_http.FormatParser{},
			&imageserver_http.QualityParser{},
		}),
		Server: &imageserver.HandlerServer{
			Server: imageserver_testdata.Server,
			Handler: &imageserver_image.Handler{
				Processor: &imageserver_image_gift.Processor{},
			},
		},
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
