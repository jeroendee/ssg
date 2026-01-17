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
	cfg     Config
	server  *http.Server
	addr    string
	handler http.Handler
	mu      sync.RWMutex
}

// New creates a new Server with the given configuration.
func New(cfg Config) *Server {
	fs := http.FileServer(http.Dir(cfg.Dir))
	handler := &xmlIndexHandler{dir: cfg.Dir, fs: fs}
	return &Server{
		cfg:     cfg,
		handler: handler,
	}
}

// Start begins serving files and blocks until the server is stopped.
// It returns http.ErrServerClosed on graceful shutdown.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.addr = ln.Addr().String()
	s.server = &http.Server{Handler: s.handler}
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

// Handler returns the HTTP handler used by the server.
func (s *Server) Handler() http.Handler {
	return s.handler
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
	// Validate path stays within h.dir to prevent traversal attacks
	if !h.isPathSafe(r.URL.Path) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

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

// isPathSafe validates that the resolved path stays within h.dir.
func (h *xmlIndexHandler) isPathSafe(urlPath string) bool {
	// Reject any path containing .. as defense-in-depth
	if strings.Contains(urlPath, "..") {
		return false
	}

	// Clean and resolve the path
	cleanPath := filepath.Clean(urlPath)
	fullPath := filepath.Join(h.dir, cleanPath)

	// Get absolute paths for comparison
	absDir, err := filepath.Abs(h.dir)
	if err != nil {
		return false
	}
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return false
	}

	// Path is safe if it starts with the base directory (belt and suspenders)
	return strings.HasPrefix(absPath, absDir)
}
