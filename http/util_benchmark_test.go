package http

import (
	"net/http"
	"testing"
	"time"
)

func BenchmarkExpiresHandler(b *testing.B) {
	eh := &ExpiresHandler{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			_, _ = rw.Write(nil)
		}),
		Expires: 1 * time.Hour,
	}
	nrw := new(nopResponseWriter)
	req := new(http.Request)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eh.ServeHTTP(nrw, req)
	}
}
