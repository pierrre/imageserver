package imageproxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Server struct {
	httpServer *http.Server
	cache      Cache
}

func NewServer(httpServer *http.Server, cache Cache) *Server {
	return &Server{
		httpServer: httpServer,
		cache:      cache,
	}
}

func (server *Server) Run() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", server.handleHttpRequest)
	server.httpServer.Handler = serveMux
	server.httpServer.ListenAndServe()
}

func (server *Server) handleHttpRequest(writer http.ResponseWriter, request *http.Request) {
	source, parameters, error := server.parseRequest(request)
	fmt.Println(source, parameters, error)
}

func (server *Server) parseRequest(request *http.Request) (source *url.URL, parameters *Parameters, err error) {
	if request.Method != "GET" {
		err = errors.New("Invalid request method")
		return
	}

	query := request.URL.Query()

	if len(query["source"]) == 0 {
		err = errors.New("Missing source parameter")
		return
	}
	source, err = url.ParseRequestURI(query["source"][0])
	if err != nil {
		err = fmt.Errorf("Invalid source parameter (%s)", err)
		return
	}

	parameters = &Parameters{}

	if len(query["width"]) > 0 {
		parameters.Width, err = strconv.Atoi(query["width"][0])
		if err != nil {
			return
		}
		if parameters.Width < 0 {
			err = errors.New("Invalid width parameter")
			return
		}
	}

	if len(query["height"]) > 0 {
		parameters.Height, err = strconv.Atoi(query["height"][0])
		if err != nil {
			return
		}
		if parameters.Height < 0 {
			err = errors.New("Invalid height parameter")
			return
		}
	}

	return
}
