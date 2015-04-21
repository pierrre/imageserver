package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ Server = ServerFunc(nil)

func TestServerFunc(t *testing.T) {
	called := false
	server := ServerFunc(func(params Params) (*Image, error) {
		called = true
		return testdata.Medium, nil
	})
	server.Get(Params{})
	if !called {
		t.Fatal("not called")
	}
}

var _ Server = &SourceServer{}

func TestSourceServer(t *testing.T) {
	server := &SourceServer{
		Server: testdata.Server,
	}
	im, err := server.Get(Params{SourceParam: testdata.MediumFileName})
	if err != nil {
		t.Fatal(err)
	}
	if im != testdata.Medium {
		t.Fatal("not equal")
	}
}

func TestSourceServerParam(t *testing.T) {
	server := &SourceServer{
		Server: ServerFunc(func(params Params) (*Image, error) {
			if !params.Has(SourceParam) {
				t.Fatal("no source param")
			}
			if params.Has("foo") {
				t.Fatal("unexpected param")
			}
			return testdata.Medium, nil
		}),
	}
	server.Get(Params{
		SourceParam: testdata.MediumFileName,
		"foo":       "bar",
	})
}

func TestSourceServerError(t *testing.T) {
	server := &SourceServer{
		Server: testdata.Server,
	}
	_, err := server.Get(Params{SourceParam: "foobar"})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSourceServerErrorNoSource(t *testing.T) {
	server := &SourceServer{
		Server: testdata.Server,
	}
	_, err := server.Get(Params{})
	if err == nil {
		t.Fatal("no error")
	}
}
