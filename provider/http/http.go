// Package http provides a http Image Provider
package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pierrre/imageserver"
)

var contentTypeRegexp = regexp.MustCompile("^image/(.+)$")

// HTTPProvider represents a http Image Provider
type HTTPProvider struct{}

// Get returns an Image for an http source
//
// If the source is not an url, the string representation of the source will be used to create one.
//
// Returns an error if the http status code is not 200 (OK).
//
// The image type is determined by the "Content-Type" header.
func (provider *HTTPProvider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	response, err := provider.getResponse(source)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	image, err := provider.createImage(response)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (provider *HTTPProvider) getResponse(source interface{}) (*http.Response, error) {
	sourceURL, err := provider.getSourceURL(source)
	if err != nil {
		return nil, err
	}

	response, err := provider.request(sourceURL)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (provider *HTTPProvider) getSourceURL(source interface{}) (*url.URL, error) {
	sourceURL, ok := source.(*url.URL)
	if !ok {
		var err error
		sourceURL, err = url.ParseRequestURI(fmt.Sprint(source))
		if err != nil {
			return nil, imageserver.NewError("Invalid source url")
		}
	}

	if sourceURL.Scheme != "http" && sourceURL.Scheme != "https" {
		return nil, imageserver.NewError("Invalid source scheme")
	}

	return sourceURL, nil
}

func (provider *HTTPProvider) request(sourceURL *url.URL) (*http.Response, error) {
	//TODO optional http client
	return http.Get(sourceURL.String())
}

func (provider *HTTPProvider) createImage(response *http.Response) (*imageserver.Image, error) {
	if err := provider.checkResponse(response); err != nil {
		return nil, err
	}

	image := new(imageserver.Image)

	provider.parseFormat(response, image)

	if err := provider.parseData(response, image); err != nil {
		return nil, err
	}

	return image, nil
}

func (provider *HTTPProvider) checkResponse(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		return imageserver.NewError(fmt.Sprintf("http status code %d while downloading source", response.StatusCode))
	}
	return nil
}

func (provider *HTTPProvider) parseFormat(response *http.Response, image *imageserver.Image) {
	contentType := response.Header.Get("Content-Type")
	if len(contentType) == 0 {
		return
	}

	matches := contentTypeRegexp.FindStringSubmatch(contentType)
	if matches == nil || len(matches) != 2 {
		return
	}

	image.Format = matches[1]
}

func (provider *HTTPProvider) parseData(response *http.Response, image *imageserver.Image) error {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return imageserver.NewError("error while downloading source")
	}

	image.Data = data

	return nil
}
