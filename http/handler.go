// Package http provides an HTTP Handler for the imageserver package
package http

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pierrre/imageserver"
)

var inmHeaderRegexp = regexp.MustCompile("^\"(.+)\"$")

// ImageHTTPHandler represents an HTTP Handler for imageserver.Server
type ImageHTTPHandler struct {
	Parser      Parser                                         // parse request to Parameters
	ImageServer imageserver.ImageServer                        // handle image requests
	ETagFunc    func(parameters imageserver.Parameters) string // optional
	ErrorFunc   func(err error, request *http.Request)         // allows to handle internal errors, optional
}

// ServeHTTP implements the HTTP Handler interface
//
// Only GET and HEAD methods are supported.
//
// Supports ETag/If-None-Match (status code 304).
// It doesn't check if the image really exists.
func (handler *ImageHTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" && request.Method != "HEAD" {
		handler.sendError(writer, request, NewErrorDefaultText(http.StatusMethodNotAllowed))
		return
	}

	parameters := make(imageserver.Parameters)
	if err := handler.Parser.Parse(request, parameters); err != nil {
		handler.sendError(writer, request, err)
		return
	}

	if handler.checkNotModified(writer, request, parameters) {
		return
	}

	image, err := handler.ImageServer.Get(parameters)
	if err != nil {
		handler.sendError(writer, request, err)
		return
	}

	if err := handler.sendImage(writer, request, parameters, image); err != nil {
		handler.callErrFunc(err, request)
		return
	}
}

func (handler *ImageHTTPHandler) checkNotModified(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) bool {
	if handler.ETagFunc == nil {
		return false
	}

	inmHeader := request.Header.Get("If-None-Match")
	if inmHeader == "" {
		return false
	}

	matches := inmHeaderRegexp.FindStringSubmatch(inmHeader)
	if matches == nil {
		return false
	}

	inm := matches[1]
	etag := handler.ETagFunc(parameters)
	if inm != etag {
		return false
	}

	handler.setImageHeaderCommon(writer, request, parameters)

	writer.WriteHeader(http.StatusNotModified)

	return true
}

func (handler *ImageHTTPHandler) sendImage(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters, image *imageserver.Image) error {
	handler.setImageHeaderCommon(writer, request, parameters)

	if image.Format != "" {
		writer.Header().Set("Content-Type", "image/"+image.Format)
	}

	writer.Header().Set("Content-Length", strconv.Itoa(len(image.Data)))

	if request.Method == "GET" {
		if _, err := writer.Write(image.Data); err != nil {
			return err
		}
	}

	return nil
}

func (handler *ImageHTTPHandler) setImageHeaderCommon(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) {
	header := writer.Header()

	header.Set("Cache-Control", "public")

	if handler.ETagFunc != nil {
		header.Set("ETag", fmt.Sprintf("\"%s\"", handler.ETagFunc(parameters)))
	}
}

func (handler *ImageHTTPHandler) sendError(writer http.ResponseWriter, request *http.Request, err error) {
	var statusCode int
	var message string

	switch err := err.(type) {
	case *imageserver.Error:
		statusCode = http.StatusBadRequest
		message = err.Error()
	case *Error:
		statusCode = err.Code
		message = err.Text
	default:
		statusCode = http.StatusInternalServerError
		message = http.StatusText(statusCode)

		handler.callErrFunc(err, request)
	}

	http.Error(writer, message, statusCode)
}

func (handler *ImageHTTPHandler) callErrFunc(err error, request *http.Request) {
	if handler.ErrorFunc != nil {
		handler.ErrorFunc(err, request)
	}
}

// NewParametersHashETagFunc returns a function that hashes the parameters and returns an ETag value
func NewParametersHashETagFunc(newHashFunc func() hash.Hash) func(parameters imageserver.Parameters) string {
	return func(parameters imageserver.Parameters) string {
		hash := newHashFunc()
		io.WriteString(hash, parameters.String())
		data := hash.Sum(nil)
		return hex.EncodeToString(data)
	}
}
