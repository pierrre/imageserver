package imageproxy

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

type Server struct {
	HttpServer    *http.Server
	RequestParser RequestParser
	Cache         Cache
	Converter     Converter
}

func (server *Server) Run() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", server.handleHttpRequest)
	server.HttpServer.Handler = serveMux
	server.HttpServer.ListenAndServe()
}

func (server *Server) handleHttpRequest(writer http.ResponseWriter, request *http.Request) {
	source, parameters, err := server.RequestParser.ParseRequest(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = parameters.Validate()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	sourceImage, err := server.getSourceImage(source)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	image, err := server.convertImage(sourceImage, parameters)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	server.sendImage(writer, image)
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

func (server *Server) convertImage(sourceImage *Image, parameters *Parameters) (image *Image, err error) {
	if server.Converter != nil {
		image, err = server.Converter.Convert(sourceImage, parameters)
	} else {
		image = sourceImage
	}

	return
}

func (server *Server) sendImage(writer http.ResponseWriter, image *Image) {
	if len(image.Type) > 0 {
		writer.Header().Set("Content-Type", "image/"+image.Type)
	}

	writer.Write(image.Data)
}
