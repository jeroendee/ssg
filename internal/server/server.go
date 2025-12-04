package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// Config holds server configuration settings.
type Config struct {
	// Port is the TCP port number to listen on.
	Port int

	// Dir is the directory to serve files from.
	Dir string
}

// Server serves static files over HTTP.
type Server struct {
	cfg    Config
	server *http.Server
	addr   string
}

// New creates a new Server with the given configuration.
func New(cfg Config) *Server {
	return &Server{cfg: cfg}
}

// Start begins serving files and blocks until the server is stopped.
// It returns http.ErrServerClosed on graceful shutdown.
func (s *Server) Start() error {
	handler := http.FileServer(http.Dir(s.cfg.Dir))

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}
	s.addr = ln.Addr().String()

	s.server = &http.Server{Handler: handler}
	return s.server.Serve(ln)
}

// Addr returns the address the server is listening on.
// Returns empty string if server has not started.
func (s *Server) Addr() string {
	return s.addr
}

// Shutdown gracefully stops the server, allowing in-flight requests to complete.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}
