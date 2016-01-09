package http

import (
	"net/http"
	"time"
)

// CacheControlPublicHandler is a net/http.Handler implementation that sets the "Cache-Control" header to "public".
//
// It only sets the header if the status code is StatusOK/204 or StatusNotModified/304.
type CacheControlPublicHandler struct {
	http.Handler
}

// ServeHTTP implements net/http.Handler
func (h *CacheControlPublicHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hrw := &headerResponseWriter{
		ResponseWriter: rw,
		OnWriteHeaderFunc: func(code int) {
			if code == http.StatusOK || code == http.StatusNotModified {
				rw.Header().Set("Cache-Control", "public")
			}
		},
	}
	h.Handler.ServeHTTP(hrw, req)
}

var expiresHeaderLocation = getTimeLocation("GMT")

func getTimeLocation(name string) *time.Location {
	l, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return l
}

// ExpiresHandler is a net/http.Handler implementation that sets the "Expires" header.
//
// It only sets the header if the status code is StatusOK/204 or StatusNotModified/304.
type ExpiresHandler struct {
	http.Handler
	Expires time.Duration
}

// ServeHTTP implements net/http.Handler.
func (eh *ExpiresHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hrw := &headerResponseWriter{
		ResponseWriter: rw,
		OnWriteHeaderFunc: func(code int) {
			if code == http.StatusOK || code == http.StatusNotModified {
				t := time.Now()
				t = t.Add(eh.Expires)
				t = t.In(expiresHeaderLocation)
				rw.Header().Set("Expires", t.Format(time.RFC1123))
			}
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
