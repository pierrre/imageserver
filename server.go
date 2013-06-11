package imageserver

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var sourceContentTypeHeaderRegexp, _ = regexp.Compile("^image/(.+)$")

type Server struct {
	Cache       Cache
	SourceCache Cache
	Converter   Converter
}

func (server *Server) GetImage(parameters Parameters) (image *Image, err error) {
	cacheKey := parameters.Hash()

	if server.Cache != nil {
		image, err = server.Cache.Get(cacheKey)
		if err == nil {
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

func (server *Server) getSourceImage(parameters Parameters) (image *Image, err error) {
	//TODO source provider
	source, err := parameters.GetString("source")
	if err != nil {
		return
	}

	cacheKey := server.hashCacheKey(source) //TODO cache source provider

	if server.SourceCache != nil {
		image, _ = server.SourceCache.Get(cacheKey)
		if image != nil {
			return
		}
	}

	sourceUrl, err := url.ParseRequestURI(source)
	if err != nil {
		return
	}
	if sourceUrl.Scheme != "http" && sourceUrl.Scheme != "https" {
		err = fmt.Errorf("Invalid source scheme")
		return
	}
	source = sourceUrl.String()

	response, err := http.Get(source)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err = fmt.Errorf("Error while downloading source")
		return
	}

	image = &Image{}

	contentType := response.Header.Get("Content-Type")
	if len(contentType) > 0 {
		matches := sourceContentTypeHeaderRegexp.FindStringSubmatch(contentType)
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

func (server *Server) convertImage(sourceImage *Image, parameters Parameters) (image *Image, err error) {
	if server.Converter != nil {
		image, err = server.Converter.Convert(sourceImage, parameters)
	} else {
		image = sourceImage
	}

	return
}

func (server *Server) hashCacheKey(key string) string {
	hash := md5.New()
	io.WriteString(hash, key)
	data := hash.Sum(nil)
	hashedKey := hex.EncodeToString(data)
	return hashedKey
}
