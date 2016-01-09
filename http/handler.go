// Package http provides a HTTP Handler for an Image Server.
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

// Handler is a HTTP Handler for imageserver.Server.
//
// Only GET and HEAD methods are supported.
//
// Supports ETag/If-None-Match (status code 304).
// It doesn't check if the image really exists.
type Handler struct {
	Parser    Parser                                 // parse request to Params
	Server    imageserver.Server                     // handle image requests
	ETagFunc  func(params imageserver.Params) string // optional
	ErrorFunc func(err error, req *http.Request)     // allows to handle internal errors, optional
}

// ServeHTTP implements http.Handler.
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
