// Package http provides a HTTP Image Provider
package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pierrre/imageserver"
	imageserver_provider "github.com/pierrre/imageserver/provider"
)

var contentTypeRegexp = regexp.MustCompile("^image/(.+)$")

// Provider represents a HTTP Image Provider
type Provider struct{}

// Get returns an Image for an HTTP source
//
// If the source is not an url, the string representation of the source will be used to create one.
//
// Returns an error if the HTTP status code is not 200 (OK).
//
// The image type is determined by the "Content-Type" header.
func (provider *Provider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	sourceURL, err := provider.getSourceURL(source)
	if err != nil {
		return nil, err
	}

	response, err := provider.doRequest(sourceURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	image, err := provider.parseResponse(response)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (provider *Provider) getSourceURL(source interface{}) (*url.URL, error) {
	sourceURL, ok := source.(*url.URL)
	if !ok {
		var err error
		sourceURL, err = url.ParseRequestURI(fmt.Sprint(source))
		if err != nil {
			return nil, &imageserver_provider.SourceError{Message: "parse url error"}
		}
	}

	if sourceURL.Scheme != "http" && sourceURL.Scheme != "https" {
		return nil, &imageserver_provider.SourceError{Message: "url scheme must be http(s)"}
	}

	return sourceURL, nil
}

func (provider *Provider) doRequest(sourceURL *url.URL) (*http.Response, error) {
	//TODO optional http client
	response, err := http.Get(sourceURL.String())
	if err != nil {
		return nil, &imageserver_provider.SourceError{Message: err.Error()}
	}
	return response, nil
}

func (provider *Provider) parseResponse(response *http.Response) (*imageserver.Image, error) {
	if response.StatusCode != http.StatusOK {
		return nil, &imageserver_provider.SourceError{Message: fmt.Sprintf("http status code %d while downloading", response.StatusCode)}
	}

	image := new(imageserver.Image)

	contentType := response.Header.Get("Content-Type")
	if contentType != "" {
		matches := contentTypeRegexp.FindStringSubmatch(contentType)
		if matches != nil && len(matches) == 2 {
			image.Format = matches[1]
		}
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &imageserver_provider.SourceError{Message: fmt.Sprintf("error while downloading: %s", err)}
	}
	image.Data = data

	return image, nil
}
