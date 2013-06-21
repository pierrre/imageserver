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

type HttpProvider struct {
}

func (provider *HttpProvider) Get(source interface{}, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	sourceUrl, err := provider.getSourceUrl(source)
	if err != nil {
		return
	}
	response, err := provider.request(sourceUrl)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if err = provider.checkResponse(response); err != nil {
		return
	}
	image, err = provider.createImage(response)
	if err != nil {
		return
	}
	return
}

func (provider *HttpProvider) getSourceUrl(source interface{}) (sourceUrl *url.URL, err error) {
	sourceUrl, ok := source.(*url.URL)
	if !ok {
		sourceUrl, err = url.ParseRequestURI(fmt.Sprint(source))
		if err != nil {
			return
		}
	}
	if sourceUrl.Scheme != "http" && sourceUrl.Scheme != "https" {
		err = fmt.Errorf("Invalid scheme")
		return
	}
	return
}

func (provider *HttpProvider) request(sourceUrl *url.URL) (response *http.Response, err error) {
	//TODO optional http client
	return http.Get(sourceUrl.String())
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
