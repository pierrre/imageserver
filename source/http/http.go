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

type HttpSource struct {
}

func (source *HttpSource) Get(sourceId string, parameters imageserver.Parameters) (image *imageserver.Image, err error) {
	if sourceId, err = source.validate(sourceId); err != nil {
		return
	}
	response, err := source.request(sourceId)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if err = source.checkResponse(response); err != nil {
		return
	}
	if image, err = source.createImage(response); err != nil {
		return
	}
	return
}

func (source *HttpSource) validate(sourceIdIn string) (sourceIdOut string, err error) {
	sourceUrl, err := url.ParseRequestURI(sourceIdIn)
	if err != nil {
		return
	}
	if sourceUrl.Scheme != "http" && sourceUrl.Scheme != "https" {
		err = fmt.Errorf("Invalid scheme")
		return
	}
	sourceIdOut = sourceUrl.String()
	return
}

func (source *HttpSource) request(sourceId string) (response *http.Response, err error) {
	//TODO optional http client
	return http.Get(sourceId)
}

func (source *HttpSource) checkResponse(response *http.Response) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("Error while downloading")
	}
	return nil
}

func (source *HttpSource) createImage(response *http.Response) (image *imageserver.Image, err error) {
	image = &imageserver.Image{}
	source.parseType(response, image)
	if err = source.parseData(response, image); err != nil {
		image = nil
		return
	}
	return
}

func (source *HttpSource) parseType(response *http.Response, image *imageserver.Image) {
	contentType := response.Header.Get("Content-Type")
	if len(contentType) > 0 {
		matches := sourceContentTypeHeaderRegexp.FindStringSubmatch(contentType)
		if matches != nil && len(matches) == 2 {
			image.Type = matches[1]
		}
	}
}

func (source *HttpSource) parseData(response *http.Response, image *imageserver.Image) error {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	image.Data = data
	return nil
}
