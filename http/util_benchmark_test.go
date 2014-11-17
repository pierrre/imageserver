package http

import (
	"net/http"
	"testing"
	"time"
)

func BenchmarkExpiresHandler(b *testing.B) {
	eh := &ExpiresHandler{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write(nil)
		}),
		Expires: 1 * time.Hour,
	}
	nrw := new(nopResponseWriter)
	req := new(http.Request)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			eh.ServeHTTP(nrw, req)
		}
	})
}
