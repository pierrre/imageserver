// Package httpsource provides a imageserver.Server implementation that gets the Image from an HTTP URL.
package httpsource

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pierrre/imageserver"
)

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

// Get implements imageserver.Server.
func (srv *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	response, err := srv.doRequest(params)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	return parseResponse(response)
}

func (srv *Server) doRequest(params imageserver.Params) (*http.Response, error) {
	src, err := params.GetString(imageserver.SourceParam)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return nil, newSourceError(err.Error())
	}
	c := srv.Client
	if c == nil {
		c = http.DefaultClient
	}
	response, err := c.Do(req)
	if err != nil {
		return nil, newSourceError(err.Error())
	}
	return response, nil
}

func parseResponse(response *http.Response) (*imageserver.Image, error) {
	if response.StatusCode != http.StatusOK {
		return nil, newSourceError(fmt.Sprintf("http status code %d while downloading", response.StatusCode))
	}
	im := new(imageserver.Image)
	contentType := response.Header.Get("Content-Type")
	if contentType != "" {
		const pref = "image/"
		if strings.HasPrefix(contentType, pref) {
			im.Format = contentType[len(pref):]
		}
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, newSourceError(fmt.Sprintf("error while downloading: %s", err))
	}
	im.Data = data
	return im, nil
}

func newSourceError(msg string) error {
	return &imageserver.ParamError{
		Param:   imageserver.SourceParam,
		Message: msg,
	}
}
