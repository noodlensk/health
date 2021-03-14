// Package http contains http server
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server defines the HTTP server
type Server struct {
	server *http.Server
	router *mux.Router
}

// Serve is blocking serving of HTTP requests
func (s *Server) Serve() error {
	s.registerHandlers()

	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server from HTTP connections.
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)

	return s.server.Shutdown(ctx)
}

func (s *Server) registerHandlers() {
	var next http.Handler = s.router

	s.server.Handler = next
}

// Handle the method and path with the handler
func (s *Server) Handle(method, path string, handler http.Handler) {
	s.router.Handle(path, handler).Methods(method)
}

// NewServer returns new HTTP server
func NewServer(address string) *Server {
	router := mux.NewRouter().StrictSlash(true)
	httpServer := &http.Server{
		Addr: address,
	}
	server := &Server{
		router: router,
		server: httpServer,
	}

	return server
}
