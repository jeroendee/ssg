package server_test

import (
	"context"
	"io"
	"net/http"
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
