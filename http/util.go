package http

import (
	"net/http"
	"net/url"
	"time"
)

var expiresHeaderLocation, _ = time.LoadLocation("GMT")

// ExpiresHandler is a HTTP Handler that adds "Expires" header.
//
// It only adds the header if the status code is OK (200) or NotModified (304)
type ExpiresHandler struct {
	http.Handler
	Expires time.Duration
}

// ServeHTTP implements http.Handler.
func (eh *ExpiresHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hrw := &headerResponseWriter{
		ResponseWriter: rw,
		OnWriteHeaderFunc: func(code int) {
			if code != http.StatusOK && code != http.StatusNotModified {
				return
			}
			t := time.Now()
			t = t.Add(eh.Expires)
			t = t.In(expiresHeaderLocation)
			rw.Header().Set("Expires", t.Format(time.RFC1123))
		},
	}
	eh.Handler.ServeHTTP(hrw, req)
}

type headerResponseWriter struct {
	http.ResponseWriter
	OnWriteHeaderFunc func(code int)
	wroteHeader       bool
}

func (hrw *headerResponseWriter) Write(data []byte) (int, error) {
	if !hrw.wroteHeader {
		hrw.WriteHeader(http.StatusOK)
	}
	return hrw.ResponseWriter.Write(data)
}

func (hrw *headerResponseWriter) WriteHeader(code int) {
	if hrw.wroteHeader {
		return
	}
	hrw.OnWriteHeaderFunc(code)
	hrw.ResponseWriter.WriteHeader(code)
	hrw.wroteHeader = true
}

func copyURL(u *url.URL) *url.URL {
	uTmp := *u
	uCopy := &uTmp
	if u.User != nil {
		usTmp := *u.User
		uCopy.User = &usTmp
	}
	return uCopy
}
