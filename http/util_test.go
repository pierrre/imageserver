package http

import (
	"net/http"
)

type nopResponseWriter struct {
}

func (nrw *nopResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (nrw *nopResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (nrw *nopResponseWriter) WriteHeader(int) {
}
