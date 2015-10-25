package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

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

func TestCopyURL(t *testing.T) {
	u, err := url.Parse("http://foo:bar@foo.bar/foobar?foo=bar#foobar")
	if err != nil {
		t.Fatal(err)
	}
	uCopy := copyURL(u)
	if u == uCopy {
		t.Fatal("same pointer")
	}
	if !reflect.DeepEqual(u, uCopy) {
		t.Fatal("not equals")
	}
}

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
