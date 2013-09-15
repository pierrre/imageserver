// Http provider
package http

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var contentTypeRegexp, _ = regexp.Compile("^image/(.+)$")

// Returns image from an http source
//
// If the source is not an url, the string representation of the source will be used to create one.
//
// Returns an error if the http status code is not 200 (OK).
//
// The image type is determined by the "Content-Type" header.
type HttpProvider struct {
}

func (provider *HttpProvider) Get(source interface{}, parameters imageserver.Parameters) (*imageserver.Image, error) {
	sourceUrl, err := provider.getSourceUrl(source)
	if err != nil {
		return nil, err
	}
	response, err := provider.request(sourceUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if err = provider.checkResponse(response); err != nil {
		return nil, err
	}
	image, err := provider.createImage(response)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (provider *HttpProvider) getSourceUrl(source interface{}) (*url.URL, error) {
	sourceUrl, ok := source.(*url.URL)
	if !ok {
		var err error
		sourceUrl, err = url.ParseRequestURI(fmt.Sprint(source))
		if err != nil {
			return nil, imageserver.NewError("Invalid source url")
		}
	}
	if sourceUrl.Scheme != "http" && sourceUrl.Scheme != "https" {
		return nil, imageserver.NewError("Invalid source scheme")
	}
	return sourceUrl, nil
}

func (provider *HttpProvider) request(sourceUrl *url.URL) (*http.Response, error) {
	//TODO optional http client
	return http.Get(sourceUrl.String())
}

func (provider *HttpProvider) checkResponse(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		return imageserver.NewError(fmt.Sprintf("Error %d while downloading source", response.StatusCode))
	}
	return nil
}

func (provider *HttpProvider) createImage(response *http.Response) (*imageserver.Image, error) {
	image := &imageserver.Image{}
	provider.parseType(response, image)
	if err := provider.parseData(response, image); err != nil {
		return nil, err
	}
	return image, nil
}

func (provider *HttpProvider) parseType(response *http.Response, image *imageserver.Image) {
	contentType := response.Header.Get("Content-Type")
	if len(contentType) > 0 {
		matches := contentTypeRegexp.FindStringSubmatch(contentType)
		if matches != nil && len(matches) == 2 {
			image.Type = matches[1]
		}
	}
}

func (provider *HttpProvider) parseData(response *http.Response, image *imageserver.Image) error {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	image.Data = data
	return nil
}
