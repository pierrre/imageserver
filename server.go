package imageproxy

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	source, parameters, err := server.parseRequest(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(source, parameters, err)

	sourceImage, err := server.getSourceImage(source)
	fmt.Println(sourceImage, err)
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

func (server *Server) getSourceImage(source *url.URL) (image *Image, err error) {
	response, err := http.Get(source.String())
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		err = errors.New("Error while downloading source")
		return
	}

	fmt.Println(response)

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	image = &Image{
		Data: data,
	}

	return
}
