// Package http provides an http handler for the imageserver package
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

var inmHeaderRegexp, _ = regexp.Compile("^\"(.+)\"$")

var expiresHeaderLocation, _ = time.LoadLocation("GMT")

// Server represents an http handler for imageserver.Server
type Server struct {
	Parser                               // parse request to Parameters
	ImageServer *imageserver.ImageServer // handle image requests

	ETagFunc func(parameters imageserver.Parameters) string // optional
	Expire   time.Duration                                  // set the "Expires" header, optional

	RequestFunc  func(request *http.Request) error                                         // allows to handle incoming requests (and eventually return an error), optional
	HeaderFunc   func(header http.Header, request *http.Request, err error)                // allows to set custom headers, optional
	ErrorFunc    func(err error, request *http.Request)                                    // allows to handle internal errors, optional
	ResponseFunc func(request *http.Request, statusCode int, contentSize int64, err error) // allows to handle returned responses, optional
}

// ServeHTTP implements the http handler interface
//
// Only GET and HEAD methods are supported.
//
// Supports ETag/If-None-Match (status code 304).
// It doesn't check if the image really exists.
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" && request.Method != "HEAD" {
		server.sendError(writer, request, NewError(http.StatusMethodNotAllowed))
		return
	}

	if server.RequestFunc != nil {
		if err := server.RequestFunc(request); err != nil {
			server.sendError(writer, request, err)
			return
		}
	}

	parameters := make(imageserver.Parameters)
	if err := server.Parser.Parse(request, parameters); err != nil {
		server.sendError(writer, request, err)
		return
	}

	if server.checkNotModified(writer, request, parameters) {
		return
	}

	image, err := server.ImageServer.Get(parameters)
	if err != nil {
		server.sendError(writer, request, err)
		return
	}

	if err := server.sendImage(writer, request, parameters, image); err != nil {
		server.callErrFunc(err, request)
		return
	}
}

func (server *Server) checkNotModified(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) bool {
	if server.ETagFunc == nil {
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
	etag := server.ETagFunc(parameters)
	if inm != etag {
		return false
	}

	server.setImageHeaderCommon(writer, request, parameters)

	writer.WriteHeader(http.StatusNotModified)

	server.callResponseFunc(request, http.StatusNotModified, 0, nil)

	return true
}

func (server *Server) sendImage(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters, image *imageserver.Image) error {
	server.setImageHeaderCommon(writer, request, parameters)

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

	server.callResponseFunc(request, http.StatusOK, contentSize, nil)

	return nil
}

func (server *Server) setImageHeaderCommon(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) {
	header := writer.Header()

	header.Set("Cache-Control", "public")

	if server.ETagFunc != nil {
		header.Set("ETag", fmt.Sprintf("\"%s\"", server.ETagFunc(parameters)))
	}

	if server.Expire != 0 {
		t := time.Now()
		t = t.Add(server.Expire)
		t = t.In(expiresHeaderLocation)
		header.Set("Expires", t.Format(time.RFC1123))
	}

	server.callHeaderFunc(header, request, nil)
}

func (server *Server) sendError(writer http.ResponseWriter, request *http.Request, err error) {
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

		server.callErrFunc(err, request)
	}

	server.callHeaderFunc(writer.Header(), request, err)

	http.Error(writer, message, statusCode)

	server.callResponseFunc(request, statusCode, int64(len(message)), err)
}

func (server *Server) callErrFunc(err error, request *http.Request) {
	if server.ErrorFunc != nil {
		server.ErrorFunc(err, request)
	}
}

func (server *Server) callHeaderFunc(header http.Header, request *http.Request, err error) {
	if server.HeaderFunc != nil {
		server.HeaderFunc(header, request, err)
	}
}

func (server *Server) callResponseFunc(request *http.Request, statusCode int, contentSize int64, err error) {
	if server.ResponseFunc != nil {
		server.ResponseFunc(request, statusCode, contentSize, err)
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
