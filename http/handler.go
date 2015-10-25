// Package http provides a HTTP Handler for an Image Server.
package http

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/pierrre/imageserver"
)

var inmHeaderRegexp = regexp.MustCompile("^\"(.+)\"$")

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
	if req.Method != "GET" && req.Method != "HEAD" {
		handler.sendError(rw, req, NewErrorDefaultText(http.StatusMethodNotAllowed))
		return
	}
	params := imageserver.Params{}
	if err := handler.Parser.Parse(req, params); err != nil {
		handler.sendError(rw, req, err)
		return
	}
	if handler.checkNotModified(rw, req, params) {
		return
	}
	image, err := handler.Server.Get(params)
	if err != nil {
		handler.sendError(rw, req, err)
		return
	}
	if err := handler.sendImage(rw, req, params, image); err != nil {
		handler.callErrorFunc(err, req)
		return
	}
}

func (handler *Handler) checkNotModified(rw http.ResponseWriter, req *http.Request, params imageserver.Params) bool {
	if handler.ETagFunc == nil {
		return false
	}
	inmHeader := req.Header.Get("If-None-Match")
	if inmHeader == "" {
		return false
	}
	matches := inmHeaderRegexp.FindStringSubmatch(inmHeader)
	if matches == nil || len(matches) != 2 {
		return false
	}
	inm := matches[1]
	etag := handler.ETagFunc(params)
	if inm != etag {
		return false
	}
	handler.setImageHeaderCommon(rw, req, params)
	rw.WriteHeader(http.StatusNotModified)
	return true
}

func (handler *Handler) sendImage(rw http.ResponseWriter, req *http.Request, params imageserver.Params, image *imageserver.Image) error {
	handler.setImageHeaderCommon(rw, req, params)
	if image.Format != "" {
		rw.Header().Set("Content-Type", "image/"+image.Format)
	}
	rw.Header().Set("Content-Length", strconv.Itoa(len(image.Data)))
	if req.Method == "GET" {
		if _, err := rw.Write(image.Data); err != nil {
			return err
		}
	}
	return nil
}

func (handler *Handler) setImageHeaderCommon(rw http.ResponseWriter, req *http.Request, params imageserver.Params) {
	header := rw.Header()
	header.Set("Cache-Control", "public")
	if handler.ETagFunc != nil {
		header.Set("ETag", fmt.Sprintf("\"%s\"", handler.ETagFunc(params)))
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
		handler.callErrorFunc(err, req)
		return NewErrorDefaultText(http.StatusInternalServerError)
	}
}

func (handler *Handler) callErrorFunc(err error, req *http.Request) {
	if handler.ErrorFunc != nil {
		handler.ErrorFunc(err, req)
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
