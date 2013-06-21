package http

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var sourceContentTypeHeaderRegexp, _ = regexp.Compile("^image/(.+)$")

type HttpProvider struct {
}

func (provider *HttpProvider) Get(source string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	if source, err = provider.validate(source); err != nil {
		return
	}
	response, err := provider.request(source)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if err = provider.checkResponse(response); err != nil {
		return
	}
	if image, err = provider.createImage(response); err != nil {
		return
	}
	return
}

func (provider *HttpProvider) validate(sourceIn string) (sourceOut string, err error) {
	sourceUrl, err := url.ParseRequestURI(sourceIn)
	if err != nil {
		return
	}
	if sourceUrl.Scheme != "http" && sourceUrl.Scheme != "https" {
		err = fmt.Errorf("Invalid scheme")
		return
	}
	sourceOut = sourceUrl.String()
	return
}

func (provider *HttpProvider) request(source string) (response *http.Response, err error) {
	//TODO optional http client
	return http.Get(source)
}

func (provider *HttpProvider) checkResponse(response *http.Response) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("Error while downloading")
	}
	return nil
}

func (provider *HttpProvider) createImage(response *http.Response) (image *imageserver.Image, err error) {
	image = &imageserver.Image{}
	provider.parseType(response, image)
	if err = provider.parseData(response, image); err != nil {
		image = nil
		return
	}
	return
}

func (provider *HttpProvider) parseType(response *http.Response, image *imageserver.Image) {
	contentType := response.Header.Get("Content-Type")
	if len(contentType) > 0 {
		matches := sourceContentTypeHeaderRegexp.FindStringSubmatch(contentType)
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
