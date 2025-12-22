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

	// Create temp directory with feed/index.xml (no index.html)
	dir := t.TempDir()
	feedDir := filepath.Join(dir, "feed")
	if err := os.MkdirAll(feedDir, 0755); err != nil {
		t.Fatal(err)
	}

	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Feed</title>
  </channel>
</rss>`

	if err := os.WriteFile(filepath.Join(feedDir, "index.xml"), []byte(xmlContent), 0644); err != nil {
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

	// Make request to /feed/ (trailing slash for directory)
	addr := srv.Addr()
	if addr == "" {
		t.Fatal("Addr() returned empty string")
	}

	resp, err := http.Get("http://" + addr + "/feed/")
	if err != nil {
		t.Fatalf("GET /feed/ failed: %v", err)
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
	// http.FileServer sets application/xml for .xml files
	if contentType != "application/xml" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/xml")
	}
}
