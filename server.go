package imageproxy

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Server struct {
	HttpServer    *http.Server
	RequestParser RequestParser
	Cache         Cache
	SourceCache   Cache
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

	cacheKey := hashCacheKey(fmt.Sprint(parameters))

	if server.Cache != nil {
		image, _ = server.Cache.Get(cacheKey)
		//FIX error
		if image != nil {
			return
		}
	}

	sourceImage, err := server.getSourceImage(parameters)
	if err != nil {
		return
	}

	image, err = server.convertImage(sourceImage, parameters)
	if err != nil {
		return
	}

	if server.Cache != nil {
		go func() {
			_ = server.Cache.Set(cacheKey, image)
		}()
	}

	return
}

func (server *Server) getSourceImage(parameters *Parameters) (image *Image, err error) {
	source := parameters.Source.String()
	cacheKey := hashCacheKey(source)

	if server.SourceCache != nil {
		image, _ = server.SourceCache.Get(cacheKey)
		if image != nil {
			return
		}
	}

	response, err := http.Get(source)
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

	if server.SourceCache != nil {
		go func() {
			_ = server.SourceCache.Set(cacheKey, image)
		}()
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

func hashCacheKey(key string) string {
	hash := md5.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	hashedKey := hex.EncodeToString(data)
	return hashedKey
}
