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
	ErrorFunc func(err error, request *http.Request) // allows to handle internal errors, optional
}

// ServeHTTP implements http.Handler.
func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" && request.Method != "HEAD" {
		handler.sendError(writer, request, NewErrorDefaultText(http.StatusMethodNotAllowed))
		return
	}
	params := make(imageserver.Params)
	if err := handler.Parser.Parse(request, params); err != nil {
		handler.sendError(writer, request, err)
		return
	}
	if handler.checkNotModified(writer, request, params) {
		return
	}
	image, err := handler.Server.Get(params)
	if err != nil {
		handler.sendError(writer, request, err)
		return
	}
	if err := handler.sendImage(writer, request, params, image); err != nil {
		handler.callErrorFunc(err, request)
		return
	}
}

func (handler *Handler) checkNotModified(writer http.ResponseWriter, request *http.Request, params imageserver.Params) bool {
	if handler.ETagFunc == nil {
		return false
	}
	inmHeader := request.Header.Get("If-None-Match")
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
	handler.setImageHeaderCommon(writer, request, params)
	writer.WriteHeader(http.StatusNotModified)
	return true
}

func (handler *Handler) sendImage(writer http.ResponseWriter, request *http.Request, params imageserver.Params, image *imageserver.Image) error {
	handler.setImageHeaderCommon(writer, request, params)
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

func (handler *Handler) setImageHeaderCommon(writer http.ResponseWriter, request *http.Request, params imageserver.Params) {
	header := writer.Header()
	header.Set("Cache-Control", "public")
	if handler.ETagFunc != nil {
		header.Set("ETag", fmt.Sprintf("\"%s\"", handler.ETagFunc(params)))
	}
}

func (handler *Handler) sendError(writer http.ResponseWriter, request *http.Request, err error) {
	httpErr := handler.convertGenericErrorToHTTP(err, request)
	http.Error(writer, httpErr.Text, httpErr.Code)
}

func (handler *Handler) convertGenericErrorToHTTP(err error, request *http.Request) *Error {
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
		handler.callErrorFunc(err, request)
		return NewErrorDefaultText(http.StatusInternalServerError)
	}
}

func (handler *Handler) callErrorFunc(err error, request *http.Request) {
	if handler.ErrorFunc != nil {
		handler.ErrorFunc(err, request)
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
