package http

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"net/http"
)

type Server struct {
	HttpServer  http.Server
	Parser      Parser
	ImageServer imageserver.Server
}

func (server *Server) Serve() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", server.handleHttpRequest)
	server.HttpServer.Handler = serveMux
	server.HttpServer.ListenAndServe()
}

func (server *Server) handleHttpRequest(writer http.ResponseWriter, request *http.Request) {
	image, err := server.getImage(request)

	if err != nil {
		server.sendError(writer, err)
		return
	}

	server.sendImage(writer, image)
}

func (server *Server) getImage(request *http.Request) (image *imageserver.Image, err error) {
	if request.Method != "GET" {
		err = fmt.Errorf("Invalid request method")
		return
	}

	parameters, err := server.Parser.Parse(request)
	if err != nil {
		return
	}

	image, err = server.ImageServer.GetImage(parameters)

	return
}

func (server *Server) sendImage(writer http.ResponseWriter, image *imageserver.Image) {
	if len(image.Type) > 0 {
		writer.Header().Set("Content-Type", "image/"+image.Type)
	}

	writer.Write(image.Data)
}

func (server *Server) sendError(writer http.ResponseWriter, err error) {
	http.Error(writer, err.Error(), http.StatusBadRequest)
}
