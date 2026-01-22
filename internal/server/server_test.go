package server_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jeroendee/ssg/internal/server"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cfg  server.Config
		port int
		dir  string
	}{
		{
			name: "default values",
			cfg:  server.Config{Port: 8080, Dir: "public"},
			port: 8080,
			dir:  "public",
		},
		{
			name: "custom values",
			cfg:  server.Config{Port: 3000, Dir: "/tmp/site"},
			port: 3000,
			dir:  "/tmp/site",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.cfg.Port != tt.port {
				t.Errorf("Port = %d, want %d", tt.cfg.Port, tt.port)
			}
			if tt.cfg.Dir != tt.dir {
				t.Errorf("Dir = %q, want %q", tt.cfg.Dir, tt.dir)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	cfg := server.Config{Port: 8080, Dir: "public"}
	srv := server.New(cfg)

	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestServerStart(t *testing.T) {
	t.Parallel()

	// Create temp directory with test file
	dir := t.TempDir()
	content := []byte("hello world")
	if err := os.WriteFile(filepath.Join(dir, "index.html"), content, 0644); err != nil {
		t.Fatal(err)
	}

	// Use port 0 for automatic port assignment
	srv := server.New(server.Config{Port: 0, Dir: dir})

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Wait for server to be ready
	time.Sleep(50 * time.Millisecond)

	// Make request to server
	addr := srv.Addr()
	if addr == "" {
		t.Fatal("Addr() returned empty string")
	}

	resp, err := http.Get("http://" + addr + "/index.html")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	if string(body) != "hello world" {
		t.Errorf("body = %q, want %q", body, "hello world")
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

func TestServerServesIndexXML(t *testing.T) {
	t.Parallel()

	// Test xmlIndexHandler serves index.xml when no index.html exists
	// Note: RSS feed is now at /feed.xml, but this tests other XML index files (e.g., sitemaps)
	dir := t.TempDir()
	xmlDir := filepath.Join(dir, "sitemap")
	if err := os.MkdirAll(xmlDir, 0755); err != nil {
		t.Fatal(err)
	}

	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
  </url>
</urlset>`

	if err := os.WriteFile(filepath.Join(xmlDir, "index.xml"), []byte(xmlContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Use port 0 for automatic port assignment
	srv := server.New(server.Config{Port: 0, Dir: dir})

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Wait for server to be ready
	time.Sleep(50 * time.Millisecond)

	// Make request to /sitemap/ (trailing slash for directory)
	addr := srv.Addr()
	if addr == "" {
		t.Fatal("Addr() returned empty string")
	}

	resp, err := http.Get("http://" + addr + "/sitemap/")
	if err != nil {
		t.Fatalf("GET /sitemap/ failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Verify Content-Type contains application/xml
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || (contentType != "application/xml" && contentType != "text/xml; charset=utf-8") {
		t.Errorf("Content-Type = %q, want application/xml or text/xml", contentType)
	}

	// Verify body is XML content (not directory listing)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	if string(body) != xmlContent {
		t.Errorf("body = %q, want %q", body, xmlContent)
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}
}

func TestHandler_ServesIndexHTML(t *testing.T) {
	t.Parallel()

	// Create temp directory with index.html
	dir := t.TempDir()
	content := []byte("<h1>Hello World</h1>")
	if err := os.WriteFile(filepath.Join(dir, "index.html"), content, 0644); err != nil {
		t.Fatal(err)
	}

	srv := server.New(server.Config{Dir: dir})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	if got := rec.Body.String(); got != string(content) {
		t.Errorf("body = %q, want %q", got, content)
	}
}

func TestHandler_ServesIndexXML_WhenNoIndexHTML(t *testing.T) {
	t.Parallel()

	// Create temp directory with index.xml (no index.html)
	dir := t.TempDir()
	xmlContent := []byte(`<?xml version="1.0"?><feed></feed>`)
	if err := os.WriteFile(filepath.Join(dir, "index.xml"), xmlContent, 0644); err != nil {
		t.Fatal(err)
	}

	srv := server.New(server.Config{Dir: dir})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/xml; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", contentType, "text/xml; charset=utf-8")
	}

	if got := rec.Body.String(); got != string(xmlContent) {
		t.Errorf("body = %q, want %q", got, xmlContent)
	}
}

func TestHandler_Returns404_ForMissingFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	srv := server.New(server.Config{Dir: dir})
	req := httptest.NewRequest("GET", "/missing.html", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHandler_SetsContentType_ForCSS(t *testing.T) {
	t.Parallel()

	// Create temp directory with CSS file
	dir := t.TempDir()
	cssContent := []byte("body { color: red; }")
	if err := os.WriteFile(filepath.Join(dir, "style.css"), cssContent, 0644); err != nil {
		t.Fatal(err)
	}

	srv := server.New(server.Config{Dir: dir})
	req := httptest.NewRequest("GET", "/style.css", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/css; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", contentType, "text/css; charset=utf-8")
	}
}

func TestHandler_SetsContentType_ForXML(t *testing.T) {
	t.Parallel()

	// Create temp directory with XML file
	dir := t.TempDir()
	xmlContent := []byte(`<?xml version="1.0"?><data></data>`)
	if err := os.WriteFile(filepath.Join(dir, "data.xml"), xmlContent, 0644); err != nil {
		t.Fatal(err)
	}

	srv := server.New(server.Config{Dir: dir})
	req := httptest.NewRequest("GET", "/data.xml", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	contentType := rec.Header().Get("Content-Type")
	// http.FileServer uses system MIME database; .xml may be application/xml or text/xml
	if contentType != "application/xml" && contentType != "text/xml; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/xml or text/xml", contentType)
	}
}

func TestHandler_RejectsPathTraversal(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Create a legitimate file in the directory
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("ok"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
	}{
		{
			name: "obvious traversal",
			path: "/../../../etc/passwd",
		},
		{
			name: "traversal to tmp",
			path: "/../../../tmp/secret",
		},
		{
			name: "hidden traversal with valid prefix",
			path: "/valid/../../../etc/passwd",
		},
		{
			name: "double encoded traversal",
			path: "/..%2F..%2F..%2Fetc/passwd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := server.New(server.Config{Dir: dir})
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()

			srv.Handler().ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Errorf("path %q: status = %d, want %d (Forbidden)", tt.path, rec.Code, http.StatusForbidden)
			}
		})
	}
}

func TestHandler_AllowsLegitimatePathsWithTraversalCleanup(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Create index.html at root
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("root"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create sitemap/index.xml (testing xmlIndexHandler for other XML content)
	sitemapDir := filepath.Join(dir, "sitemap")
	if err := os.MkdirAll(sitemapDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sitemapDir, "index.xml"), []byte("<urlset/>"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantContains string
	}{
		{
			name:         "root directory",
			path:         "/",
			wantStatus:   http.StatusOK,
			wantContains: "root",
		},
		{
			name:         "sitemap directory with index.xml",
			path:         "/sitemap/",
			wantStatus:   http.StatusOK,
			wantContains: "<urlset/>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := server.New(server.Config{Dir: dir})
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()

			srv.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("path %q: status = %d, want %d", tt.path, rec.Code, tt.wantStatus)
			}

			if body := rec.Body.String(); !contains(body, tt.wantContains) {
				t.Errorf("path %q: body = %q, want to contain %q", tt.path, body, tt.wantContains)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestServerTimeouts(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	srv := server.New(server.Config{Port: 0, Dir: dir})

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Wait for server to be ready
	time.Sleep(50 * time.Millisecond)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	// Get the underlying http.Server to verify timeouts
	httpServer := srv.HTTPServer()
	if httpServer == nil {
		t.Fatal("HTTPServer() returned nil")
	}

	// Verify ReadHeaderTimeout is set (prevents slowloris attacks)
	wantReadHeaderTimeout := 10 * time.Second
	if httpServer.ReadHeaderTimeout != wantReadHeaderTimeout {
		t.Errorf("ReadHeaderTimeout = %v, want %v", httpServer.ReadHeaderTimeout, wantReadHeaderTimeout)
	}

	// Verify IdleTimeout is set (closes idle keepalive connections)
	wantIdleTimeout := 120 * time.Second
	if httpServer.IdleTimeout != wantIdleTimeout {
		t.Errorf("IdleTimeout = %v, want %v", httpServer.IdleTimeout, wantIdleTimeout)
	}
}
