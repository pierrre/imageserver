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

var msgInternalError = "Internal Server Error"

// Http image server
//
// Only GET method is supported.
//
// Supports ETag/If-None-Match (status code 304).
// It doesn't check if the image really exists.
//
// Status codes: 200 (everything is ok), 400 (user error), 500 (internal error).
//
// If Expire is defined, the "Expires" header is set.
//
// The HeaderFunc function allows to set custom headers.
type Server struct {
	Parser      Parser
	ImageServer *imageserver.Server

	Expire time.Duration // optional

	ErrFunc    func(error, *http.Request)                               //optional
	HeaderFunc func(http.Header, *http.Request, imageserver.Parameters) // optional
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" && request.Method != "HEAD" {
		server.sendError(writer, request, imageserver.NewError("Invalid request method"))
		return
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

	server.setHeader(writer, request, parameters)
	writer.WriteHeader(http.StatusNotModified)
	return true
}

func (server *Server) sendImage(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters, image *imageserver.Image) error {
	server.setHeader(writer, request, parameters)

	if len(image.Type) > 0 {
		writer.Header().Set("Content-Type", "image/"+image.Type)
	}

	writer.Header().Set("Content-Length", strconv.Itoa(len(image.Data)))

	if request.Method == "GET" {
		if _, err := writer.Write(image.Data); err != nil {
			return err
		}
	}

	return nil
}

func (server *Server) setHeader(writer http.ResponseWriter, request *http.Request, parameters imageserver.Parameters) {
	header := writer.Header()

	server.setHeaderCache(header, parameters)

	if server.HeaderFunc != nil {
		server.HeaderFunc(header, request, parameters)
	}
}

func (server *Server) setHeaderCache(header http.Header, parameters imageserver.Parameters) {
	header.Set("Cache-Control", "public")

	header.Set("ETag", fmt.Sprintf("\"%s\"", parameters.Hash()))

	if server.Expire != 0 {
		t := time.Now()
		t = t.Add(server.Expire)
		t = t.In(expiresHeaderLocation)
		header.Set("Expires", t.Format(time.RFC1123))
	}
}

func (server *Server) sendError(writer http.ResponseWriter, request *http.Request, err error) {
	var message string
	var status int
	var internalErr error

	switch err := err.(type) {
	case *imageserver.Error:
		message = err.Error()
		internalErr = err.Previous
	default:
		message = msgInternalError
		internalErr = err
	}

	if internalErr != nil {
		server.callErrFunc(internalErr, request)
		status = http.StatusInternalServerError
	} else {
		status = http.StatusBadRequest
	}

	http.Error(writer, message, status)
}

func (server *Server) callErrFunc(err error, request *http.Request) {
	if server.ErrFunc != nil {
		server.ErrFunc(err, request)
	}
}
