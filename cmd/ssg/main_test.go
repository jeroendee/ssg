package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
