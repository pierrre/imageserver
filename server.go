package imageproxy

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
	cache      Cache
}

func NewServer(httpServer *http.Server, cache Cache) *Server {
	return &Server{
		httpServer: httpServer,
		cache:      cache,
	}
}

func (server *Server) Run() {
	serveMux := http.NewServeMux();
	serveMux.HandleFunc("/", server.handleHttpRequest);
	server.httpServer.Handler = serveMux
	server.httpServer.ListenAndServe()
}

func (server *Server) handleHttpRequest(w http.ResponseWriter, r *http.Request) {
}
