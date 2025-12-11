package main

import (
	"fmt"

	"github.com/jeroendee/ssg/internal/builder"
	"github.com/jeroendee/ssg/internal/config"
	"github.com/spf13/cobra"
)

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
	b.SetVersion(Version)
	b.SetAssetsDir(assetsDir)

	if err := b.Build(); err != nil {
		return fmt.Errorf("building site: %w", err)
	}

	fmt.Println("Site built successfully!")
	return nil
}
