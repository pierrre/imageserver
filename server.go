package imageproxy

import (
	"errors"
	"io/ioutil"
	"net/http"
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
	image, err := server.getImage(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	server.sendImage(writer, image)
}

func (server *Server) getImage(request *http.Request) (image *Image, err error) {
	parameters, err := server.RequestParser.ParseRequest(request)
	if err != nil {
		return
	}
	err = parameters.Validate()
	if err != nil {
		return
	}

	image, _ = server.Cache.Get("lol")
	if image != nil {
		return
	}

	sourceImage, err := server.getSourceImage(parameters)
	if err != nil {
		return
	}

	image, err = server.convertImage(sourceImage, parameters)
	if err != nil {
		return
	}

	_ = server.Cache.Set("lol", image)

	return
}

func (server *Server) getSourceImage(parameters *Parameters) (image *Image, err error) {
	response, err := http.Get(parameters.Source.String())
	if err != nil {
		return
	}
	defer response.Body.Close()

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
