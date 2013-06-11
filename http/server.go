package http

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"net/http"
	"regexp"
	"time"
)

var inmHeaderRegexp, _ = regexp.Compile("^\"(.+)\"$")

var expiresHeaderLocation, _ = time.LoadLocation("GMT")

type Server struct {
	HttpServer  *http.Server
	Parser      Parser
	ImageServer *imageserver.Server

	ExpiresDuration time.Duration
}

func (server *Server) Serve() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", server.handleHttpRequest)
	server.HttpServer.Handler = serveMux
	server.HttpServer.ListenAndServe()
}

func (server *Server) handleHttpRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		server.sendError(writer, fmt.Errorf("Invalid request method"))
		return
	}

	parameters, err := server.Parser.Parse(request)
	if err != nil {
		server.sendError(writer, err)
		return
	}

	if server.checkNotModified(writer, request, parameters) {
		return
	}

	image, err := server.ImageServer.GetImage(parameters)

	if err != nil {
		server.sendError(writer, err)
		return
	}

	server.sendImage(writer, image, parameters)
}

func (server *Server) checkNotModified(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) bool {
	inmHeader := request.Header.Get("If-None-Match")
	if len(inmHeader) > 0 {
		matches := inmHeaderRegexp.FindStringSubmatch(inmHeader)
		if matches != nil && len(matches) == 2 {
			inm := matches[1]
			if inm == parameters.Hash() {
				writer.WriteHeader(304)
				return true
			}
		}
	}
	return false
}

func (server *Server) sendImage(writer http.ResponseWriter, image *imageserver.Image, parameters imageserver.Parameters) {
	if len(image.Type) > 0 {
		writer.Header().Set("Content-Type", "image/"+image.Type)
	}

	writer.Header().Set("Cache-Control", "public")

	if server.ExpiresDuration != 0 {
		t := time.Now()
		t = t.Add(server.ExpiresDuration)
		t = t.In(expiresHeaderLocation)
		writer.Header().Set("Expires", t.Format(time.RFC1123))
	}

	writer.Header().Set("ETag", fmt.Sprintf("\"%s\"", parameters.Hash()))

	writer.Write(image.Data)
}

func (server *Server) sendError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusBadRequest)
}
