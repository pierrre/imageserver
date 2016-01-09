// Package httpsource provides a imageserver.Server implementation that gets the Image from an HTTP URL.
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

// Server is a imageserver.Server implementation that gets the Image from an HTTP URL.
//
// It parses the "source" param as URL, then do a GET request.
// It returns an error if the HTTP status code is not 200 (OK).
//
// The Image type is determined by the "Content-Type" header.
type Server struct {
	// Client is an optional HTTP client.
	// http.DefaultClient is used by default.
	Client *http.Client
}

// Get implements Server.
func (srv *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	sourceURL, err := getSourceURL(params)
	if err != nil {
		return nil, err
	}
	response, err := srv.doRequest(sourceURL)
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
	source, err := params.GetString(imageserver.SourceParam)
	if err != nil {
		return nil, err
	}
	sourceURL, err := url.ParseRequestURI(source)
	if err != nil {
		return nil, &imageserver.ParamError{
			Param:   imageserver.SourceParam,
			Message: fmt.Sprintf("parse url error: %s", err),
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

func (srv *Server) doRequest(sourceURL *url.URL) (*http.Response, error) {
	c := srv.Client
	if c == nil {
		c = http.DefaultClient
	}
	response, err := c.Get(sourceURL.String())
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
