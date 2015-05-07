// Package httpsource provides a HTTP source Image Server.
package httpsource

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pierrre/imageserver"
)

var contentTypeRegexp = regexp.MustCompile("^image/(.+)$")

// Server is a HTTP source Image Server.
type Server struct{}

// Get returns an Image for a HTTP source.
//
// If the source is not an url, the string representation of the source will be used to create one.
//
// Returns an error if the HTTP status code is not 200 (OK).
//
// The image type is determined by the "Content-Type" header.
func (server *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	sourceURL, err := getSourceURL(params)
	if err != nil {
		return nil, err
	}
	response, err := doRequest(sourceURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	image, err := parseResponse(response)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func getSourceURL(params imageserver.Params) (*url.URL, error) {
	source, err := params.Get(imageserver.SourceParam)
	if err != nil {
		return nil, err
	}
	sourceURL, ok := source.(*url.URL)
	if !ok {
		sourceURL, err = url.ParseRequestURI(fmt.Sprint(source))
		if err != nil {
			return nil, &imageserver.ParamError{
				Param:   imageserver.SourceParam,
				Message: "parse url error",
			}
		}
	}
	if sourceURL.Scheme != "http" && sourceURL.Scheme != "https" {
		return nil, &imageserver.ParamError{
			Param:   imageserver.SourceParam,
			Message: "url scheme must be http(s)",
		}
	}
	return sourceURL, nil
}

func doRequest(sourceURL *url.URL) (*http.Response, error) {
	//TODO optional http client
	response, err := http.Get(sourceURL.String())
	if err != nil {
		return nil, &imageserver.ParamError{Param: imageserver.SourceParam, Message: err.Error()}
	}
	return response, nil
}

func parseResponse(response *http.Response) (*imageserver.Image, error) {
	if response.StatusCode != http.StatusOK {
		return nil, &imageserver.ParamError{
			Param:   imageserver.SourceParam,
			Message: fmt.Sprintf("http status code %d while downloading", response.StatusCode),
		}
	}
	im := new(imageserver.Image)
	contentType := response.Header.Get("Content-Type")
	if contentType != "" {
		matches := contentTypeRegexp.FindStringSubmatch(contentType)
		if matches != nil && len(matches) == 2 {
			im.Format = matches[1]
		}
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &imageserver.ParamError{
			Param:   imageserver.SourceParam,
			Message: fmt.Sprintf("error while downloading: %s", err),
		}
	}
	im.Data = data
	return im, nil
}
