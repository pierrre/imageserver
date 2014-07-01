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
	"time"

	"github.com/pierrre/imageserver"
)

var inmHeaderRegexp = regexp.MustCompile("^\"(.+)\"$")

var expiresHeaderLocation, _ = time.LoadLocation("GMT")

// ImageHTTPHandler represents an HTTP Handler for imageserver.Server
type ImageHTTPHandler struct {
	Parser      Parser                  // parse request to Parameters
	ImageServer imageserver.ImageServer // handle image requests

	ETagFunc func(parameters imageserver.Parameters) string // optional
	Expire   time.Duration                                  // set the "Expires" header, optional

	RequestFunc  func(request *http.Request) error                                         // allows to handle incoming requests (and eventually return an error), optional
	HeaderFunc   func(header http.Header, request *http.Request, err error)                // allows to set custom headers, optional
	ErrorFunc    func(err error, request *http.Request)                                    // allows to handle internal errors, optional
	ResponseFunc func(request *http.Request, statusCode int, contentSize int64, err error) // allows to handle returned responses, optional
}

// ServeHTTP implements the HTTP Handler interface
//
// Only GET and HEAD methods are supported.
//
// Supports ETag/If-None-Match (status code 304).
// It doesn't check if the image really exists.
func (handler *ImageHTTPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" && request.Method != "HEAD" {
		handler.sendError(writer, request, NewError(http.StatusMethodNotAllowed))
		return
	}

	if handler.RequestFunc != nil {
		if err := handler.RequestFunc(request); err != nil {
			handler.sendError(writer, request, err)
			return
		}
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
	if len(inmHeader) == 0 {
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

	handler.callResponseFunc(request, http.StatusNotModified, 0, nil)

	return true
}

func (handler *ImageHTTPHandler) sendImage(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters, image *imageserver.Image) error {
	handler.setImageHeaderCommon(writer, request, parameters)

	if len(image.Format) > 0 {
		writer.Header().Set("Content-Type", "image/"+image.Format)
	}

	contentLength := len(image.Data)
	writer.Header().Set("Content-Length", strconv.Itoa(contentLength))

	var contentSize int64
	if request.Method == "GET" {
		contentSize = int64(contentLength)
		if _, err := writer.Write(image.Data); err != nil {
			return err
		}
	}

	handler.callResponseFunc(request, http.StatusOK, contentSize, nil)

	return nil
}

func (handler *ImageHTTPHandler) setImageHeaderCommon(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) {
	header := writer.Header()

	header.Set("Cache-Control", "public")

	if handler.ETagFunc != nil {
		header.Set("ETag", fmt.Sprintf("\"%s\"", handler.ETagFunc(parameters)))
	}

	if handler.Expire != 0 {
		t := time.Now()
		t = t.Add(handler.Expire)
		t = t.In(expiresHeaderLocation)
		header.Set("Expires", t.Format(time.RFC1123))
	}

	handler.callHeaderFunc(header, request, nil)
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
		message = err.Error()
	default:
		statusCode = http.StatusInternalServerError
		message = http.StatusText(statusCode)

		handler.callErrFunc(err, request)
	}

	handler.callHeaderFunc(writer.Header(), request, err)

	http.Error(writer, message, statusCode)

	handler.callResponseFunc(request, statusCode, int64(len(message)), err)
}

func (handler *ImageHTTPHandler) callErrFunc(err error, request *http.Request) {
	if handler.ErrorFunc != nil {
		handler.ErrorFunc(err, request)
	}
}

func (handler *ImageHTTPHandler) callHeaderFunc(header http.Header, request *http.Request, err error) {
	if handler.HeaderFunc != nil {
		handler.HeaderFunc(header, request, err)
	}
}

func (handler *ImageHTTPHandler) callResponseFunc(request *http.Request, statusCode int, contentSize int64, err error) {
	if handler.ResponseFunc != nil {
		handler.ResponseFunc(request, statusCode, contentSize, err)
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
