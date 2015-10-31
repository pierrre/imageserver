package groupcache

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/golang/groupcache"
	"github.com/pierrre/imageserver"
)

// HTTPPoolContextHeader is the header used to store the Context in HTTP requests.
const HTTPPoolContextHeader = "X-Imageserver-Groupcache-Context"

// HTTPPoolContext must be used in groupcache.HTTPPool.Context.
func HTTPPoolContext(req *http.Request) groupcache.Context {
	ctx, err := getContext(req)
	if err != nil {
		return nil
	}
	return ctx
}

func getContext(req *http.Request) (*Context, error) {
	h := req.Header.Get(HTTPPoolContextHeader)
	if h == "" {
		return nil, fmt.Errorf("header is not set")
	}
	return decodeContext(h)
}

func decodeContext(s string) (*Context, error) {
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var ctx *Context
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(&ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

// NewHTTPPoolTransport returns a function that must be used in groupcache.HTTPPool.Transport.
//
// rt is optional, http.DefaultTransport is used by default.
func NewHTTPPoolTransport(rt http.RoundTripper) func(groupcache.Context) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return func(ctx groupcache.Context) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if ctx, ok := ctx.(*Context); ok && ctx != nil {
				err := setContext(req, ctx)
				if err != nil {
					return nil, err
				}
			}
			return rt.RoundTrip(req)
		})
	}
}

func setContext(req *http.Request, ctx *Context) error {
	h, err := encodeContext(ctx)
	if err != nil {
		return err
	}
	req.Header.Set(HTTPPoolContextHeader, h)
	return nil
}

func encodeContext(ctx *Context) (string, error) {
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(ctx)
	if err != nil {
		return "", err
	}
	s := base64.RawURLEncoding.EncodeToString(buf.Bytes())
	return s, nil
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func init() {
	gob.Register(new(Context))
	gob.Register(imageserver.Params{})
}
