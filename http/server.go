// Http image server
package http

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var inmHeaderRegexp, _ = regexp.Compile("^\"(.+)\"$")

var expiresHeaderLocation, _ = time.LoadLocation("GMT")

var msgInternalError = "Internal error"

// Http image server
//
// Only GET method is supported.
//
// Supports ETag/If-None-Match (status code 304).
// It doesn't check if the image really exists.
//
// Status codes: 200 (everything is ok), 400 (user error), 500 (internal error).
//
// If Expire is defined, the "Expires" header is set.
//
// The HeaderFunc function allows to set custom headers.
type Server struct {
	Parser      Parser
	ImageServer *imageserver.Server

	Expire time.Duration // optional

	HeaderFunc func(http.Header, *http.Request, imageserver.Parameters) // optional
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		server.sendError(writer, fmt.Errorf("Invalid request method"))
		return
	}

	parameters := make(imageserver.Parameters)
	if err := server.Parser.Parse(request, parameters); err != nil {
		server.sendError(writer, err)
		return
	}

	if server.checkNotModified(writer, request, parameters) {
		return
	}

	image, err := server.ImageServer.Get(parameters)
	if err != nil {
		server.sendError(writer, err)
		return
	}

	server.sendImage(writer, request, parameters, image)
}

func (server *Server) checkNotModified(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) bool {
	inmHeader := request.Header.Get("If-None-Match")
	if len(inmHeader) > 0 {
		matches := inmHeaderRegexp.FindStringSubmatch(inmHeader)
		if matches != nil && len(matches) == 2 {
			inm := matches[1]
			if inm == parameters.Hash() {
				server.sendHeader(writer, request, parameters)
				writer.WriteHeader(http.StatusNotModified)
				return true
			}
		}
	}
	return false
}

func (server *Server) sendImage(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters, image *imageserver.Image) {
	server.sendHeader(writer, request, parameters)

	if len(image.Type) > 0 {
		writer.Header().Set("Content-Type", "image/"+image.Type)
	}

	writer.Header().Set("Content-Length", strconv.Itoa(len(image.Data)))

	writer.Write(image.Data)
}

func (server *Server) sendHeader(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) {
	header := writer.Header()
	if server.HeaderFunc != nil {
		server.HeaderFunc(header, request, parameters)
	}
	server.sendHeaderCache(header, parameters)
}

func (server *Server) sendHeaderCache(header http.Header, parameters imageserver.Parameters) {
	header.Set("Cache-Control", "public")

	header.Set("ETag", fmt.Sprintf("\"%s\"", parameters.Hash()))

	if server.Expire != 0 {
		t := time.Now()
		t = t.Add(server.Expire)
		t = t.In(expiresHeaderLocation)
		header.Set("Expires", t.Format(time.RFC1123))
	}
}

func (server *Server) sendError(writer http.ResponseWriter, err error) {
	if _, ok := err.(*imageserver.Error); ok {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(writer, msgInternalError, http.StatusInternalServerError)
	}
}
