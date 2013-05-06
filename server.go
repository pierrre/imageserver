package imageproxy

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
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

	sourceString := query.Get("source")
	if len(sourceString) == 0 {
		err = errors.New("Missing source parameter")
		return
	}
	source, err = url.ParseRequestURI(sourceString)
	if err != nil {
		err = fmt.Errorf("Invalid source parameter (%s)", err)
		return
	}

	parameters = &Parameters{}

	widthString := query.Get("width")
	if len(widthString) > 0 {
		parameters.Width, err = strconv.Atoi(widthString)
		if err != nil {
			return
		}
		if parameters.Width < 0 {
			err = errors.New("Invalid width parameter")
			return
		}
	}

	heightString := query.Get("height")
	if len(heightString) > 0 {
		parameters.Height, err = strconv.Atoi(heightString)
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

	image = &Image{}

	contentType := response.Header.Get("Content-Type")
	if len(contentType) > 0 {
		r, _ := regexp.Compile("image/(.+)")
		matches := r.FindStringSubmatch(contentType)
		if matches != nil && len(matches) == 2 {
			image.Type = matches[1]
		}
	}

	defer response.Body.Close()
	image.Data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	return
}
