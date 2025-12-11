package main

import (
	"github.com/spf13/cobra"
)

// Version variables set via -ldflags at build time
var (
	Version   = "dev"     // Git SHA or version tag
	BuildDate = "unknown" // Build timestamp
)

var version = Version // For backward compatibility with cobra

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "ssg",
		Short:   "A static site generator for blogs",
		Long:    "ssg is a static site generator that converts markdown files into a complete website.",
		Version: version,
	}

	rootCmd.AddCommand(newBuildCmd())
	rootCmd.AddCommand(newServeCmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}
