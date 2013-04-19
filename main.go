package main

import (
	"bytes"
	"encoding/gob"
	"github.com/bradfitz/gomemcache/memcache"
	"net/http"
)

func main() {
	server := &Server{
		Cache: &MemcacheImageCache{
			memcache: memcache.New("localhost"),
		},
	}

	http.HandleFunc("/", server.HandleHttpRequest)
	http.ListenAndServe(":8080", nil)
}

type Server struct {
	Cache ImageCache
}

func (server *Server) HandleHttpRequest(w http.ResponseWriter, r *http.Request) {

}

type Image struct {
	Type string
	Data []byte
}

func (image *Image) serialize() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(image)
	return buffer.Bytes()
}

func (image *Image) unserialize(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(image)
}

type ImageCache interface {
	Get(key string) *Image
	Set(key string, image *Image)
}

type MemcacheImageCache struct {
	prefix   string
	memcache *memcache.Client
}

func (c *MemcacheImageCache) Get(key string) *Image {
	return nil
}

func (c *MemcacheImageCache) Set(key string, image *Image) {
}

type ImageConverter interface {
	Convert(image *Image, parameters *ImageConverterParameters) *Image
}

type ImageConverterParameters struct {
	Width  int
	Height int
}

type GraphicsMagickImageConverter struct {
	executable string
}

func (converter *GraphicsMagickImageConverter) Convert(image *Image, parameters *ImageConverterParameters) *Image {
	return nil
}
