package imageproxy

import (
	"net/http"
)

type Server struct {
	cache Cache
}

func NewServer(cache Cache) *Server {
	return &Server{
		cache: cache,
	}
}

func (server *Server) HandleHttpRequest(w http.ResponseWriter, r *http.Request) {

}
