// Package http provides a net/http.Handler implementation that wraps a imageserver.Server.
package http

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/pierrre/imageserver"
)

// Handler is a net/http.Handler implementation that wraps a imageserver.Server.
//
// Supported methods are: GET and HEAD.
// Other method will return a StatusMethodNotAllowed/405 response.
//
// It supports ETag/If-None-Match headers and returns a StatusNotModified/304 response accordingly.
// But it doesn't check if the Image really exists (the Server is not called).
//
// Steps:
//  - Parse the HTTP request, and fill the Params.
//  - If the given If-None-Match header matches the ETag, return a StatusNotModified/304 response.
//  - Call the Server and get the Image.
//  - Return a StatusOK/200 response containing the Image.
//
// Errors (returned by Parser or Server):
//  - *imageserver/http.Error will return a response with the given status code and message.
//  - *imageserver.ParamError will return a StatusBadRequest/400 response, with a message including the resolved HTTP param.
//  - *imageserver.ImageError will return a StatusBadRequest/400 response, with the given message.
//  - Other error will return a StatusInternalServerError/500 response, and ErrorFunc will be called.
//
// Returned headers:
//  - Content-Type is set for StatusOK/200 response, and contains "image/{Image.Format}".
//  - Content-Length is set for StatusOK/200 response, and contains the Image size.
//  - ETag is set for StatusOK/200 and StatusNotModified/304 response, and contains the ETag value.
type Handler struct {
	// Parser parses the HTTP request and fills the Params.
	Parser Parser

	// Server handles the image request.
	Server imageserver.Server

	// ETagFunc is an optional function that returns the ETag value for the given Params.
	// See https://en.wikipedia.org/wiki/HTTP_ETag .
	// The returned value must not be enclosed in quotes (they are added automatically).
	ETagFunc func(params imageserver.Params) string

	// ErrorFunc is an optional function that is called if there is an internal error.
	ErrorFunc func(err error, req *http.Request)
}

// ServeHTTP implements net/http.Handler.
func (handler *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	err := handler.serveHTTP(rw, req)
	if err != nil {
		handler.sendError(rw, req, err)
	}
}

func (handler *Handler) serveHTTP(rw http.ResponseWriter, req *http.Request) error {
	if req.Method != "GET" && req.Method != "HEAD" {
		return NewErrorDefaultText(http.StatusMethodNotAllowed)
	}
	params := imageserver.Params{}
	err := handler.Parser.Parse(req, params)
	if err != nil {
		return err
	}
	etag := handler.getETag(params)
	if handler.checkNotModified(rw, req, etag) {
		return nil
	}
	image, err := handler.Server.Get(params)
	if err != nil {
		return err
	}
	handler.sendImage(rw, req, image, etag)
	return nil
}

func (handler *Handler) getETag(params imageserver.Params) string {
	if handler.ETagFunc != nil {
		return "\"" + handler.ETagFunc(params) + "\""
	}
	return ""
}

func (handler *Handler) checkNotModified(rw http.ResponseWriter, req *http.Request, etag string) bool {
	if etag == "" {
		return false
	}
	inm := req.Header.Get("If-None-Match")
	if inm != etag {
		return false
	}
	handler.setImageHeaderCommon(rw, req, etag)
	rw.WriteHeader(http.StatusNotModified)
	return true
}

func (handler *Handler) sendImage(rw http.ResponseWriter, req *http.Request, image *imageserver.Image, etag string) {
	handler.setImageHeaderCommon(rw, req, etag)
	if image.Format != "" {
		rw.Header().Set("Content-Type", "image/"+image.Format)
	}
	rw.Header().Set("Content-Length", strconv.Itoa(len(image.Data)))
	if req.Method == "GET" {
		rw.Write(image.Data)
	}
}

func (handler *Handler) setImageHeaderCommon(rw http.ResponseWriter, req *http.Request, etag string) {
	if etag != "" {
		rw.Header().Set("ETag", etag)
	}
}

func (handler *Handler) sendError(rw http.ResponseWriter, req *http.Request, err error) {
	httpErr := handler.convertGenericErrorToHTTP(err, req)
	http.Error(rw, httpErr.Text, httpErr.Code)
}

func (handler *Handler) convertGenericErrorToHTTP(err error, req *http.Request) *Error {
	switch err := err.(type) {
	case *Error:
		return err
	case *imageserver.ParamError:
		httpParam := handler.Parser.Resolve(err.Param)
		if httpParam == "" {
			httpParam = err.Param
		}
		text := fmt.Sprintf("invalid param \"%s\": %s", httpParam, err.Message)
		return &Error{Code: http.StatusBadRequest, Text: text}
	case *imageserver.ImageError:
		text := fmt.Sprintf("image error: %s", err.Message)
		return &Error{Code: http.StatusBadRequest, Text: text}
	default:
		if handler.ErrorFunc != nil {
			handler.ErrorFunc(err, req)
		}
		return NewErrorDefaultText(http.StatusInternalServerError)
	}
}

// NewParamsHashETagFunc returns a function that hashes the params and returns an ETag value.
//
// It is intended to be used in Handler.ETagFunc.
func NewParamsHashETagFunc(newHashFunc func() hash.Hash) func(params imageserver.Params) string {
	pool := &sync.Pool{
		New: func() interface{} {
			return newHashFunc()
		},
	}
	return func(params imageserver.Params) string {
		h := pool.Get().(hash.Hash)
		io.WriteString(h, params.String())
		data := h.Sum(nil)
		h.Reset()
		pool.Put(h)
		return hex.EncodeToString(data)
	}
}
