package http

import (
	"encoding/hex"
	"github.com/pierrre/imageserver"
	"hash"
	"io"
)

type ETagProvider interface {
	Get(parameters imageserver.Parameters) string
}

type HashParametersETagProvider struct {
	HashFunc func() hash.Hash
}

func (etagProvider *HashParametersETagProvider) Get(parameters imageserver.Parameters) string {
	hash := etagProvider.HashFunc()
	io.WriteString(hash, parameters.String())
	data := hash.Sum(nil)
	return hex.EncodeToString(data)
}
