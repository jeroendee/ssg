package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeroendee/ssg/internal/config"
	"github.com/jeroendee/ssg/internal/model"
)

func TestLoad_ValidConfig(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Quality Shepherd"
  baseURL: "https://www.qualityshepherd.nl"
  author: "Jeroen"

build:
  content: "content"
  output: "public"
  assets: "assets"

navigation:
  - title: "Home"
    url: "/"
  - title: "About"
    url: "/about/"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Title != "Quality Shepherd" {
		t.Errorf("Title = %q, want %q", cfg.Title, "Quality Shepherd")
	}
	if cfg.BaseURL != "https://www.qualityshepherd.nl" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://www.qualityshepherd.nl")
	}
	if cfg.Author != "Jeroen" {
		t.Errorf("Author = %q, want %q", cfg.Author, "Jeroen")
	}
	if cfg.ContentDir != "content" {
		t.Errorf("ContentDir = %q, want %q", cfg.ContentDir, "content")
	}
	if cfg.OutputDir != "public" {
		t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, "public")
	}
	if cfg.AssetsDir != "assets" {
		t.Errorf("AssetsDir = %q, want %q", cfg.AssetsDir, "assets")
	}
	if len(cfg.Navigation) != 2 {
		t.Errorf("len(Navigation) = %d, want %d", len(cfg.Navigation), 2)
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.ContentDir != "content" {
		t.Errorf("ContentDir = %q, want default %q", cfg.ContentDir, "content")
	}
	if cfg.OutputDir != "public" {
		t.Errorf("OutputDir = %q, want default %q", cfg.OutputDir, "public")
	}
	if cfg.AssetsDir != "assets" {
		t.Errorf("AssetsDir = %q, want default %q", cfg.AssetsDir, "assets")
	}
}

func TestLoad_MissingTitle(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  baseURL: "https://example.com"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := config.Load(cfgFile)
	if err == nil {
		t.Error("Load() expected error for missing title")
	}
}

func TestLoad_MissingBaseURL(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := config.Load(cfgFile)
	if err == nil {
		t.Error("Load() expected error for missing baseURL")
	}
}

func TestLoad_NonExistentFile(t *testing.T) {
	t.Parallel()
	_, err := config.Load("/nonexistent/path/ssg.yaml")
	if err == nil {
		t.Error("Load() expected error for non-existent file")
	}
}

func TestLoad_WithOverrides(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
build:
  content: "content"
  output: "public"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	opts := config.Options{
		ContentDir: "custom-content",
		OutputDir:  "custom-output",
		AssetsDir:  "custom-assets",
	}
	cfg, err := config.LoadWithOptions(cfgFile, opts)
	if err != nil {
		t.Fatalf("LoadWithOptions() error = %v", err)
	}
	if cfg.ContentDir != "custom-content" {
		t.Errorf("ContentDir = %q, want %q", cfg.ContentDir, "custom-content")
	}
	if cfg.OutputDir != "custom-output" {
		t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, "custom-output")
	}
	if cfg.AssetsDir != "custom-assets" {
		t.Errorf("AssetsDir = %q, want %q", cfg.AssetsDir, "custom-assets")
	}
}

func TestNavItemUnmarshal(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test"
  baseURL: "https://example.com"
navigation:
  - title: "Blog"
    url: "/blog/"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := model.NavItem{Title: "Blog", URL: "/blog/"}
	if len(cfg.Navigation) != 1 || cfg.Navigation[0] != want {
		t.Errorf("Navigation = %v, want %v", cfg.Navigation, []model.NavItem{want})
	}
}

func TestLoad_WithLogo(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
  logo: "/logo.svg"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Logo != "/logo.svg" {
		t.Errorf("Logo = %q, want %q", cfg.Logo, "/logo.svg")
	}
}

func TestLoad_WithoutLogo(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Logo != "" {
		t.Errorf("Logo = %q, want empty string", cfg.Logo)
	}
}

func TestLoad_WithDescription(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
  description: "A blog about testing"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Description != "A blog about testing" {
		t.Errorf("Description = %q, want %q", cfg.Description, "A blog about testing")
	}
}

func TestLoad_WithoutDescription(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Description != "" {
		t.Errorf("Description = %q, want empty string", cfg.Description)
	}
}

func TestLoad_WithAssetsDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
build:
  assets: "static"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AssetsDir != "static" {
		t.Errorf("AssetsDir = %q, want %q", cfg.AssetsDir, "static")
	}
}

func TestLoad_DefaultAssetsDir(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AssetsDir != "assets" {
		t.Errorf("AssetsDir = %q, want default %q", cfg.AssetsDir, "assets")
	}
}

func TestLoad_AssetsDirOverride(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
build:
  assets: "static"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	opts := config.Options{
		AssetsDir: "cli-assets",
	}
	cfg, err := config.LoadWithOptions(cfgFile, opts)
	if err != nil {
		t.Fatalf("LoadWithOptions() error = %v", err)
	}
	if cfg.AssetsDir != "cli-assets" {
		t.Errorf("AssetsDir = %q, want %q", cfg.AssetsDir, "cli-assets")
	}
}

func TestLoad_WithGoatCounter(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
analytics:
  goatcounter: "aishepherd"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Analytics.GoatCounter != "aishepherd" {
		t.Errorf("Analytics.GoatCounter = %q, want %q", cfg.Analytics.GoatCounter, "aishepherd")
	}
}

func TestLoad_WithoutGoatCounter(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Analytics.GoatCounter != "" {
		t.Errorf("Analytics.GoatCounter = %q, want empty string", cfg.Analytics.GoatCounter)
	}
}

func TestLoad_WithOGImage(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
  ogImage: "/social-image.png"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.OGImage != "/social-image.png" {
		t.Errorf("OGImage = %q, want %q", cfg.OGImage, "/social-image.png")
	}
}

func TestLoad_WithoutOGImage(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "ssg.yaml")
	content := `site:
  title: "Test Site"
  baseURL: "https://example.com"
`
	if err := os.WriteFile(cfgFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.OGImage != "" {
		t.Errorf("OGImage = %q, want empty string", cfg.OGImage)
	}
}
