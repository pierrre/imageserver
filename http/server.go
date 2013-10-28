// Http image server
package http

import (
	"fmt"
	"github.com/pierrre/imageserver"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var inmHeaderRegexp, _ = regexp.Compile("^\"(.+)\"$")

var expiresHeaderLocation, _ = time.LoadLocation("GMT")

// Http image server
//
// Only GET and HEAD methods are supported.
//
// Supports ETag/If-None-Match (status code 304).
// It doesn't check if the image really exists.
//
// Status codes: 200 (everything is ok), 400 (user error), 500 (internal error).
//
// If Expire is defined, the "Expires" header is set.
//
// The ErrFunc function allows to handler internal errors.
//
// The HeaderFunc function allows to set custom headers.
type Server struct {
	Parser      Parser
	ImageServer *imageserver.Server

	Expire time.Duration // optional

	RequestFunc  func(request *http.Request) error                                         // optional
	HeaderFunc   func(header http.Header, request *http.Request, err error)                // optional
	ErrorFunc    func(err error, request *http.Request)                                    // optional
	ResponseFunc func(request *http.Request, statusCode int, contentSize int64, err error) // optional
}

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
	inmHeader := request.Header.Get("If-None-Match")
	if len(inmHeader) == 0 {
		return false
	}

	matches := inmHeaderRegexp.FindStringSubmatch(inmHeader)
	if matches == nil {
		return false
	}

	inm := matches[1]
	if inm != parameters.Hash() {
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

	header.Set("ETag", fmt.Sprintf("\"%s\"", parameters.Hash()))

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
