package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	mu     sync.RWMutex
}

// New creates a new Server with the given configuration.
func New(cfg Config) *Server {
	return &Server{cfg: cfg}
}

// Start begins serving files and blocks until the server is stopped.
// It returns http.ErrServerClosed on graceful shutdown.
func (s *Server) Start() error {
	fs := http.FileServer(http.Dir(s.cfg.Dir))
	handler := &xmlIndexHandler{dir: s.cfg.Dir, fs: fs}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.addr = ln.Addr().String()
	s.server = &http.Server{Handler: handler}
	s.mu.Unlock()

	return s.server.Serve(ln)
}

// Addr returns the address the server is listening on.
// Returns empty string if server has not started.
func (s *Server) Addr() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.addr
}

// Shutdown gracefully stops the server, allowing in-flight requests to complete.
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.RLock()
	srv := s.server
	s.mu.RUnlock()

	if srv == nil {
		return nil
	}
	return srv.Shutdown(ctx)
}

// xmlIndexHandler wraps http.FileServer to serve index.xml for directories without index.html.
type xmlIndexHandler struct {
	dir string
	fs  http.Handler
}

// ServeHTTP serves index.xml with proper Content-Type when a directory has no index.html.
func (h *xmlIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only handle directory requests (ending with /)
	if !strings.HasSuffix(r.URL.Path, "/") {
		h.fs.ServeHTTP(w, r)
		return
	}

	dirPath := filepath.Join(h.dir, r.URL.Path)

	// Check if index.html exists - delegate to FileServer
	htmlPath := filepath.Join(dirPath, "index.html")
	if _, err := os.Stat(htmlPath); err == nil {
		h.fs.ServeHTTP(w, r)
		return
	}

	// Check if index.xml exists - serve it directly
	xmlPath := filepath.Join(dirPath, "index.xml")
	if _, err := os.Stat(xmlPath); err == nil {
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		http.ServeFile(w, r, xmlPath)
		return
	}

	// No index file found - delegate to FileServer (shows directory listing)
	h.fs.ServeHTTP(w, r)
}
