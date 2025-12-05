// Ssg is a static site generator that converts markdown files into a blog-style website.
//
// Usage:
//
//	ssg build [flags]
//
// The flags are:
//
//	-c, --config string
//		Path to config file (default "ssg.yaml")
//	-o, --output string
//		Output directory (overrides config)
//	--content string
//		Content directory (overrides config)
//	--assets string
//		Assets directory (default "assets")
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeroendee/ssg/internal/builder"
	"github.com/jeroendee/ssg/internal/config"
	"github.com/jeroendee/ssg/internal/server"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ssg",
		Short:   "A static site generator for blogs",
		Long:    "ssg is a static site generator that converts markdown files into a complete website.",
		Version: version,
	}

	rootCmd.AddCommand(newBuildCmd())
	rootCmd.AddCommand(newServeCmd())

	return rootCmd
}

func newBuildCmd() *cobra.Command {
	var (
		configPath string
		outputDir  string
		contentDir string
		assetsDir  string
	)

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build the static site",
		Long:  "Build generates the static site by processing markdown files and producing HTML output.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBuild(configPath, outputDir, contentDir, assetsDir)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "ssg.yaml", "path to config file")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory (overrides config)")
	cmd.Flags().StringVar(&contentDir, "content", "", "content directory (overrides config)")
	cmd.Flags().StringVar(&assetsDir, "assets", "assets", "assets directory")

	return cmd
}

func runBuild(configPath, outputDir, contentDir, assetsDir string) error {
	// Load config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Apply overrides
	if outputDir != "" {
		cfg.OutputDir = outputDir
	}
	if contentDir != "" {
		cfg.ContentDir = contentDir
	}

	// Create and run builder
	b := builder.New(cfg)
	b.SetAssetsDir(assetsDir)

	if err := b.Build(); err != nil {
		return fmt.Errorf("building site: %w", err)
	}

	fmt.Println("Site built successfully!")
	return nil
}

func newServeCmd() *cobra.Command {
	var (
		configPath string
		port       int
		dir        string
		build      bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start a development server",
		Long:  "Serve starts a local HTTP server to preview your site during development.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(configPath, port, dir, build)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "ssg.yaml", "path to config file")
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "port to serve on")
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to serve (overrides config output_dir)")
	cmd.Flags().BoolVarP(&build, "build", "b", false, "build site before serving")

	return cmd
}

func runServe(configPath string, port int, dir string, doBuild bool) error {
	// Create context that cancels on interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	return runServeWithContext(ctx, configPath, port, dir, doBuild, nil)
}

func runServeWithContext(ctx context.Context, configPath string, port int, dir string, doBuild bool, addrCh chan<- string) error {
	// Load config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Determine serve directory
	serveDir := cfg.OutputDir
	if dir != "" {
		serveDir = dir
	}

	// Validate port range (0 is allowed for auto-assignment in tests)
	if port < 0 || port > 65535 {
		return fmt.Errorf("invalid port %d: must be between 0 and 65535", port)
	}

	// Optionally build first
	if doBuild {
		b := builder.New(cfg)
		b.SetAssetsDir("assets")
		if err := b.Build(); err != nil {
			return fmt.Errorf("building site: %w", err)
		}
		fmt.Println("Site built successfully!")
	}

	// Validate serve directory exists (skip if building, as build creates it)
	if !doBuild {
		if _, err := os.Stat(serveDir); os.IsNotExist(err) {
			return fmt.Errorf("serve directory does not exist: %s", serveDir)
		}
	}

	// Create and start server
	srv := server.New(server.Config{
		Port: port,
		Dir:  serveDir,
	})

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Wait briefly for server to start, then report address
	time.Sleep(50 * time.Millisecond)
	addr := srv.Addr()
	if addrCh != nil {
		addrCh <- addr
	}
	fmt.Printf("Serving %s at http://%s\n", serveDir, addr)
	fmt.Println("Press Ctrl+C to stop")

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		fmt.Println("\nShutting down...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
