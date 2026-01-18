package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jeroendee/ssg/internal/server"
)

func TestVersionVariables_DefaultValues(t *testing.T) {
	t.Parallel()

	// Version should default to "dev" when not set via ldflags
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if Version != "dev" {
		t.Logf("Version = %q (may be set via ldflags)", Version)
	}

	// BuildDate should default to "unknown" when not set via ldflags
	if BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}
	if BuildDate != "unknown" {
		t.Logf("BuildDate = %q (may be set via ldflags)", BuildDate)
	}
}

func TestVersionVariables_AreExported(t *testing.T) {
	t.Parallel()

	// Verify variables are exported (accessible from test package)
	// This ensures they can be set via -ldflags
	_ = Version
	_ = BuildDate

	// Verify the version used by cobra matches Version
	cmd := newRootCmd()
	if cmd.Version != Version {
		t.Errorf("cobra Version = %q, want %q", cmd.Version, Version)
	}
}

func TestRootCommand_Version(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--version"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ssg") {
		t.Error("version output should contain 'ssg'")
	}
}

func TestVersionCommand_OutputFormat(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--version"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	// Should contain version info
	if !strings.Contains(output, "ssg") {
		t.Error("version output should contain 'ssg'")
	}
	// Version string should be included (whether "dev" or a git SHA)
	if !strings.Contains(output, "version") && !strings.Contains(output, "ssg version") {
		t.Logf("version output: %q", output)
	}
}

func TestRootCommand_Help(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"--help"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "static site generator") {
		t.Error("help should describe static site generator")
	}
}

func TestBuildCommand_Help(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"build", "--help"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "--config") {
		t.Error("help should document --config flag")
	}
	if !strings.Contains(output, "--output") {
		t.Error("help should document --output flag")
	}
	if !strings.Contains(output, "--content") {
		t.Error("help should document --content flag")
	}
}

func TestBuildCommand_WithValidConfig(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(contentDir, 0755)

	// Create required home.md
	os.WriteFile(filepath.Join(contentDir, "home.md"), []byte("---\ntitle: Home\n---\nWelcome"), 0644)

	// Create config file with correct YAML structure
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + outputDir + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"build", "--config", configPath})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Check output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("build should create output directory")
	}
}

func TestBuildCommand_WithOverrideFlags(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "custom-output")
	os.MkdirAll(contentDir, 0755)

	// Create required home.md
	os.WriteFile(filepath.Join(contentDir, "home.md"), []byte("---\ntitle: Home\n---\nWelcome"), 0644)

	// Create config file with different output dir
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + filepath.Join(tmpDir, "default-output") + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"build", "--config", configPath, "--output", outputDir})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Check custom output directory was used
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("build should use --output flag override")
	}
}

func TestBuildCommand_MissingConfig(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"build", "--config", "/nonexistent/config.yaml"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("Execute() should error with missing config file")
	}
}

func TestRootCommand_DefaultsToHelp(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Without args, should show usage/help
	output := buf.String()
	if output == "" {
		t.Error("root command without args should produce output")
	}
}

func TestServeCommand_Help(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"serve", "--help"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "--port") {
		t.Error("help should document --port flag")
	}
	if !strings.Contains(output, "--config") {
		t.Error("help should document --config flag")
	}
	if !strings.Contains(output, "--dir") {
		t.Error("help should document --dir flag")
	}
	if !strings.Contains(output, "--build") {
		t.Error("help should document --build flag")
	}
}

func TestServeCommand_MissingConfig(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"serve", "--config", "/nonexistent/config.yaml"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("Execute() should error with missing config file")
	}
}

// startTestServer starts a test server in a goroutine and returns its address and cleanup function.
// The cleanup function cancels the context and waits for the server to stop.
func startTestServer(t *testing.T, configPath string, doBuild bool) (addr string, cleanup func()) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	addrCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		errCh <- runServeWithContext(ctx, configPath, 0, "", doBuild, addrCh)
	}()

	// Wait for server to start or error
	select {
	case addr = <-addrCh:
		// Server started successfully
	case err := <-errCh:
		t.Fatalf("server failed to start: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for server to start")
	}

	cleanup = func() {
		cancel()

		// Wait for server to stop
		select {
		case <-errCh:
			// Server stopped
		case <-time.After(2 * time.Second):
			t.Error("timeout waiting for server to stop")
		}
	}

	return addr, cleanup
}

func TestServeCommand_WithBuildFlag(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(contentDir, 0755)

	// Create required home.md
	os.WriteFile(filepath.Join(contentDir, "home.md"), []byte("---\ntitle: Home\n---\nWelcome"), 0644)

	// Create config file
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + outputDir + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Use cancellable context to stop server after build
	ctx, cancel := context.WithCancel(context.Background())
	addrCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		errCh <- runServeWithContext(ctx, configPath, 0, "", true, addrCh)
	}()

	// Wait for server to start (indicates build completed)
	select {
	case <-addrCh:
		// Server started, build completed
	case err := <-errCh:
		t.Fatalf("runServeWithContext() error = %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for server to start")
	}

	// Verify build occurred by checking output dir exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("--build flag should trigger site build")
	}

	// Shutdown server
	cancel()

	// Wait for clean shutdown
	select {
	case <-errCh:
		// Server stopped
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for server to stop")
	}
}

func TestServeCommand_StartsServer(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(outputDir, 0755)

	// Create test file in output dir
	os.WriteFile(filepath.Join(outputDir, "index.html"), []byte("hello"), 0644)

	// Create server using httptest
	srv := server.New(server.Config{Dir: outputDir})
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	// Make request to server
	resp, err := http.Get(ts.URL + "/index.html")
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

	if string(body) != "hello" {
		t.Errorf("body = %q, want %q", body, "hello")
	}
}

func TestServeCommand_DirectoryNotExists(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public") // Does not exist
	os.MkdirAll(contentDir, 0755)

	// Create config file pointing to non-existent output dir
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + outputDir + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Run serve without --build flag (directory should be validated)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := runServeWithContext(ctx, configPath, 0, "", false, nil)
	if err == nil {
		t.Fatal("expected error when directory doesn't exist")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("error should mention directory does not exist, got: %v", err)
	}
}

func TestServeCommand_WithBuildFlag_CopiesAssets(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(assetsDir, 0755)

	// Create required home.md
	os.WriteFile(filepath.Join(contentDir, "home.md"), []byte("---\ntitle: Home\n---\nWelcome"), 0644)

	// Create CSS file in assets
	cssContent := "body { color: red; }"
	os.WriteFile(filepath.Join(assetsDir, "style.css"), []byte(cssContent), 0644)

	// Create config file
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + outputDir + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Change to tmpDir so relative "assets" path works
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	addr, cleanup := startTestServer(t, configPath, true)
	defer cleanup()

	// Verify assets were copied to output directory
	copiedCSS := filepath.Join(outputDir, "style.css")
	if _, err := os.Stat(copiedCSS); os.IsNotExist(err) {
		t.Error("serve --build should copy assets to output directory")
	}

	// Verify CSS is accessible via HTTP
	resp, err := http.Get("http://" + addr + "/style.css")
	if err != nil {
		t.Fatalf("GET /style.css failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /style.css status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	if string(body) != cssContent {
		t.Errorf("GET /style.css body = %q, want %q", body, cssContent)
	}
}

func TestServeCommand_InvalidPort(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(outputDir, 0755)

	// Create config file
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + outputDir + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	tests := []struct {
		name string
		port int
	}{
		{"port negative", -1},
		{"port too high", 65536},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			err := runServeWithContext(ctx, configPath, tt.port, "", false, nil)
			if err == nil {
				t.Fatalf("expected error for port %d", tt.port)
			}

			if !strings.Contains(err.Error(), "invalid port") {
				t.Errorf("error should mention invalid port, got: %v", err)
			}
		})
	}
}

func TestVersionCommand_Exists(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	versionCmd, _, err := cmd.Find([]string{"version"})
	if err != nil {
		t.Fatalf("Find(\"version\") error = %v", err)
	}

	if versionCmd.Name() != "version" {
		t.Errorf("command name = %q, want %q", versionCmd.Name(), "version")
	}
}

func TestVersionCommand_OutputsVersionString(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	cmd.SetArgs([]string{"version"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expected := "ssg version " + Version + "\n"
	if output != expected {
		t.Errorf("version output = %q, want %q", output, expected)
	}
}

func TestBuildCommand_VersionFlowsToFooter(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(contentDir, 0755)

	// Create required home.md
	os.WriteFile(filepath.Join(contentDir, "home.md"), []byte("---\ntitle: Home\n---\nWelcome"), 0644)

	// Create a simple page
	pageContent := `---
title: Test Page
slug: test
---
# Test Content`
	os.WriteFile(filepath.Join(contentDir, "test.md"), []byte(pageContent), 0644)

	// Create config file
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + outputDir + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Run build
	err := runBuild(configPath, "", "", "")
	if err != nil {
		t.Fatalf("runBuild() error = %v", err)
	}

	// Read generated HTML
	htmlPath := filepath.Join(outputDir, "test", "index.html")
	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("reading generated HTML: %v", err)
	}

	html := string(htmlBytes)
	// Version should appear in footer
	if !strings.Contains(html, Version) {
		t.Errorf("generated HTML should contain version %q in footer, got:\n%s", Version, html)
	}
}

func TestServeCommand_WithBuild_VersionFlowsToFooter(t *testing.T) {
	t.Parallel()

	// Setup test directories
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	os.MkdirAll(contentDir, 0755)

	// Create required home.md
	os.WriteFile(filepath.Join(contentDir, "home.md"), []byte("---\ntitle: Home\n---\nWelcome"), 0644)

	// Create a simple page
	pageContent := `---
title: Test Page
slug: test
---
# Test Content`
	os.WriteFile(filepath.Join(contentDir, "test.md"), []byte(pageContent), 0644)

	// Create config file
	configPath := filepath.Join(tmpDir, "ssg.yaml")
	configContent := `site:
  title: Test Site
  baseURL: https://example.com
  author: Test Author
build:
  content: ` + contentDir + `
  output: ` + outputDir + `
`
	os.WriteFile(configPath, []byte(configContent), 0644)

	addr, cleanup := startTestServer(t, configPath, true)
	defer cleanup()

	// Read generated HTML
	htmlPath := filepath.Join(outputDir, "test", "index.html")
	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("reading generated HTML: %v", err)
	}

	html := string(htmlBytes)
	// Version should appear in footer
	if !strings.Contains(html, Version) {
		t.Errorf("generated HTML should contain version %q in footer, got:\n%s", Version, html)
	}

	// Also test via HTTP request
	resp, err := http.Get("http://" + addr + "/test/")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}

	if !strings.Contains(string(body), Version) {
		t.Errorf("HTTP response should contain version %q in footer", Version)
	}
}
