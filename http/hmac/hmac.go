package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pierrre/imageserver"
)

// Error constants retured by various functions.
const (
	NotValidError   = Error("signature not valid")
	ExpiredError    = Error("signature expired")
	BadRequestError = Error("bad request")
)

// Error is an error type
type Error string

// Error returns the stringified version of Error
func (e Error) Error() string {
	return string(e)
}

// Verifier is an imageserver.Server that verifies the signature and possibly the timestamp of a request URL
type Verifier struct {
	Next       imageserver.Server
	key        []byte
	expiration bool
}

// NewVerifier returns an initialized Verifier. If expiration is 0, expirations are not enforced.
func NewVerifier(srv imageserver.Server, key string, expiration time.Duration) imageserver.Server {
	return &Verifier{
		Next:       srv,
		key:        []byte(key),
		expiration: expiration > 0,
	}
}

// Get does the HMAC verification, and possibly expiration calculation, of the request
func (srv *Verifier) Get(params imageserver.Params) (*imageserver.Image, error) {
	source, err := params.GetString("source")
	if err != nil {
		return nil, err
	}

	var parts []string
	if parts = strings.SplitN(source, "/", 2); len(parts) != 2 {
		return nil, BadRequestError
	}
	params.Set("source", parts[1]) // set source without the hmac

	if srv.expiration {
		exp, perr := params.GetInt64(param)
		if perr != nil {
			// Couldn't convert the expiration stamp to an int!
			return nil, perr
		}

		expTime := time.UnixMilli(exp)
		if !time.Now().Before(expTime) {
			return nil, ExpiredError
		}
	}

	ok, err := verifyHMAC([]byte(params.String()), srv.key, parts[0])
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, NotValidError
	}

	// carry on!
	return srv.Next.Get(params)
}

// Signer is an http.Handler that signs the request URL + query parameters and redirects to the signed URL.
type Signer struct {
	key     []byte
	expires time.Duration
}

// NewSigner returns an initialized HMACSigner. If expires is 0, then expiration is not computed.
func NewSigner(key []byte, expires time.Duration) *Signer {
	return &Signer{
		key:     key,
		expires: expires,
	}
}

// ServeHTTP handles the request
func (s *Signer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := make(imageserver.Params)
	params.Set("source", req.URL.Path)
	for k, v := range req.URL.Query() {
		params.Set(k, v[0])
	}
	var exp time.Time
	if s.expires > 0 {
		exp = time.Now().Add(s.expires)
		params.Set(param, exp.UnixMilli())
		req.URL.RawQuery += fmt.Sprintf("&%s=%d", param, exp.UnixMilli())
	}
	hash := signHMAC([]byte(params.String()), s.key)
	http.Redirect(w, req, fmt.Sprintf("/%s/%s?%s", hash, req.URL.Path, req.URL.Query().Encode()), http.StatusTemporaryRedirect)
}

// signHMAC is the primitive signer, using sha256 and returning a base64 URL-encoded string
func signHMAC(msg, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	macSum := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(macSum)
}

// verifyHMAC returns true iff the base64-decoded hash matches the msg sum
func verifyHMAC(msg, key []byte, hash string) (bool, error) {
	sig, err := base64.RawURLEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(msg)

	return hmac.Equal(sig, mac.Sum(nil)), nil
}
