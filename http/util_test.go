package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pierrre/imageserver"
)

var _ http.Handler = &CacheControlPublicHandler{}

func TestCacheControlPublicHandler(t *testing.T) {
	h := &CacheControlPublicHandler{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte("foobar"))
		}),
	}
	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rw, req)
	cc := rw.Header().Get("Cache-Control")
	if cc != "public" {
		t.Fatal("not equal")
	}
}

func TestCacheControlPublicHandlerUndefined(t *testing.T) {
	h := &CacheControlPublicHandler{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		}),
	}
	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rw, req)
	cc := rw.Header().Get("Cache-Control")
	if cc != "" {
		t.Fatal("should not be set")
	}
}

func TestGetTimeLocation(t *testing.T) {
	getTimeLocation("GMT")
}

func TestGetTimeLocationError(t *testing.T) {
	defer func() {
		v := recover()
		if v == nil {
			t.Fatal("no panic")
		}
		_, ok := v.(error)
		if !ok {
			t.Fatalf("unexpected recover type: got %T, want error", v)
		}
	}()
	getTimeLocation("invalid")
}

var _ http.Handler = &ExpiresHandler{}

func TestExpiresHandler(t *testing.T) {
	h := &ExpiresHandler{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte("foobar"))
		}),
		Expires: 1 * time.Hour,
	}
	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rw, req)
	expires := rw.Header().Get("Expires")
	if expires == "" {
		t.Fatal("not set")
	}
}

func TestExpiresHandlerUndefined(t *testing.T) {
	h := &ExpiresHandler{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		}),
		Expires: 1 * time.Hour,
	}
	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rw, req)
	expires := rw.Header().Get("Expires")
	if expires != "" {
		t.Fatal("should not be set")
	}
}

func TestHeaderResponseWriter(t *testing.T) {
	called := false
	rw := &headerResponseWriter{
		ResponseWriter: &nopResponseWriter{},
		OnWriteHeaderFunc: func(code int) {
			called = true
		},
	}
	rw.Write([]byte("foobar"))
	if !called {
		t.Fatal("not called")
	}
	rw.WriteHeader(200)
}

type nopResponseWriter struct{}

func (nrw *nopResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (nrw *nopResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (nrw *nopResponseWriter) WriteHeader(int) {
}

type nopParser struct{}

func (p *nopParser) Parse(req *http.Request, params imageserver.Params) error {
	return nil
}

func (p *nopParser) Resolve(param string) string {
	return ""
}
