// Package file provides a imageserver.Server implementation that get the Image from a file.
package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime"
	"path"
	"path/filepath"
	"strings"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
)

// Server is a imageserver.Server implementation that get the Image from a file.
//
// It takes the "source" param and loads it from the Root directory.
// It expects a slash separated path.
type Server struct {
	// Root is the directory where images are loaded from.
	Root string

	// Identify identifies the Image format.
	// By default, it uses IdentifyMime().
	Identify func(pth string, data []byte) (format string, err error)
}

// Get implements imageserver.Server.
func (srv *Server) Get(ctx context.Context, params imageserver.Params) (*imageserver.Image, error) {
	pth, err := srv.getPath(params)
	if err != nil {
		return nil, err
	}
	data, err := loadFile(pth)
	if err != nil {
		return nil, err
	}
	format, err := srv.identify(pth, data)
	if err != nil {
		return nil, err
	}
	return &imageserver.Image{
		Format: format,
		Data:   data,
	}, nil
}

func (srv *Server) getPath(params imageserver.Params) (string, error) {
	src, err := params.GetString(imageserver_source.Param)
	if err != nil {
		return "", err
	}
	// This trick comes from net/http.Dir.Open().
	// It allows to "jail" the path inside the root.
	return filepath.Join(srv.Root, filepath.FromSlash(path.Clean("/"+src))), nil
}

func loadFile(pth string) ([]byte, error) {
	data, err := ioutil.ReadFile(pth)
	if err != nil {
		return nil, newSourceError(fmt.Sprintf("error while reading file: %s: %s", pth, err.Error()))
	}
	return data, nil
}

func (srv *Server) identify(pth string, data []byte) (format string, err error) {
	idf := srv.Identify
	if idf == nil {
		idf = IdentifyMime
	}
	format, err = idf(pth, data)
	if err != nil {
		return "", newSourceError(fmt.Sprintf("unable to identify image format: %s", err.Error()))
	}
	return format, nil
}

func newSourceError(msg string) error {
	return &imageserver.ParamError{
		Param:   imageserver_source.Param,
		Message: msg,
	}
}

// IdentifyMime identifies the Image format with the "mime" package.
func IdentifyMime(pth string, data []byte) (format string, err error) {
	ext := filepath.Ext(pth)
	if ext == "" {
		return "", fmt.Errorf("no file extension: %s", pth)
	}
	typ := mime.TypeByExtension(ext)
	if typ == "" {
		return "", fmt.Errorf("unkwnon file type for extension %s", ext)
	}
	const pref = "image/"
	if !strings.HasPrefix(typ, pref) {
		return "", fmt.Errorf("file type does not begin with \"%s\": %s", pref, typ)
	}
	return typ[len(pref):], nil
}
