package groupcache

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
)

func TestHTTPPoolContext(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx1 := &Context{
		Params: imageserver.Params{
			"foo": "bar",
		},
	}
	err = setContext(req, ctx1)
	if err != nil {
		t.Fatal(err)
	}
	tmpctx := HTTPPoolContext(req)
	ctx2, ok := tmpctx.(*Context)
	if !ok {
		t.Fatalf("unexpected context type: %T", tmpctx)
	}
	if ctx2 == nil {
		t.Fatal("context is nil")
	}
	if ctx1.Params.String() != ctx2.Params.String() {
		t.Fatal("not equals")
	}
}

func TestHTTPPoolContextErrorHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := HTTPPoolContext(req)
	if ctx != nil {
		t.Fatal("not nil")
	}
}

func TestHTTPPoolContextErrorBase64(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(HTTPPoolContextHeader, "@@@@@@@@")
	ctx := HTTPPoolContext(req)
	if ctx != nil {
		t.Fatal("not nil")
	}
}

func TestHTTPPoolContextErrorGob(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(HTTPPoolContextHeader, "aaaaaaaa")
	ctx := HTTPPoolContext(req)
	if ctx != nil {
		t.Fatal("not nil")
	}
}

func TestNewHTTPPoolTransport(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx1 := &Context{
		Params: imageserver.Params{
			"foo": "bar",
		},
	}
	resp, err := NewHTTPPoolTransport(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		ctx2, err := getContext(req)
		if err != nil {
			t.Fatal(err)
		}
		if ctx1.Params.String() != ctx2.Params.String() {
			t.Fatal("not equals")
		}
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	}))(ctx1).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestNewHTTPPoolTransportErrorGob(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &Context{
		Params: imageserver.Params{
			"foo": struct{}{},
		},
	}
	_, err = NewHTTPPoolTransport(nil)(ctx).RoundTrip(req)
	if err == nil {
		t.Fatal("no error")
	}
}
