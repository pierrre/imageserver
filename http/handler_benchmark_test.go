package http

import (
	"context"
	"crypto/sha256"
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkHandler(b *testing.B) {
	h := &Handler{
		Parser: &nopParser{},
		Server: imageserver.ServerFunc(func(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		ETagFunc: func(params imageserver.Params) string {
			return "foo"
		},
	}
	rw := &nopResponseWriter{}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set("If-None-Match", "\"bar\"")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ServeHTTP(rw, req)
	}
}

func BenchmarkNewParamsHashETagFunc(b *testing.B) {
	params := imageserver.Params{"foo": "bar"}
	f := NewParamsHashETagFunc(sha256.New)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(params)
	}
}
