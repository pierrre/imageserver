package http

import (
	"net/http"
	"time"
)

var expiresHeaderLocation, _ = time.LoadLocation("GMT")

// ExpiresHandler adds "Expires" header
//
// It only adds the header if the status code is OK (200) or NotModified (304)
type ExpiresHandler struct {
	http.Handler
	Expires time.Duration
}

func (eh *ExpiresHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hrw := &headerResponseWriter{
		ResponseWriter: w,
		OnWriteHeaderFunc: func(code int) {
			if code != http.StatusOK && code != http.StatusNotModified {
				return
			}
			t := time.Now()
			t = t.Add(eh.Expires)
			t = t.In(expiresHeaderLocation)
			w.Header().Set("Expires", t.Format(time.RFC1123))
		},
	}
	eh.Handler.ServeHTTP(hrw, r)
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
