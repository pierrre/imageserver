package groupcache

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

func BenchmarkHTTPPoolContext(b *testing.B) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		b.Fatal(err)
	}
	ctx := &Context{
		Params: imageserver.Params{
			"foo": "bar",
		},
	}
	err = setContext(req, ctx)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := HTTPPoolContext(req)
		if ctx == nil {
			b.Fatal("context is nil")
		}
	}
}

func BenchmarkNewHTTPPoolTransport(b *testing.B) {
	resp := &http.Response{}
	ctx := &Context{
		Params: imageserver.Params{
			"foo": "bar",
		},
	}
	rt := NewHTTPPoolTransport(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return resp, nil
	}))(ctx)
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rt.RoundTrip(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
