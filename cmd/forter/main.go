// Package main is the entry point for the forter CLI tool
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/OmkarVijayBagade/forter/internal/config"
	"github.com/OmkarVijayBagade/forter/internal/tui"
)

var (
	// CLI flags
	dryRun    bool
	recursive bool
	version   bool

	// Version info (set by ldflags during build)
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "forter [path]",
		Short: "A high-performance TUI file organizer - forter",
		Long: `forter is a fast, interactive terminal file organizer built with Go and Bubble Tea.

It scans directories and helps you organize files into categories based on their extensions.
Features include:
  • Interactive TUI with vim-like navigation
  • Automatic category detection (Documents, Images, Videos, etc.)
  • Bulk selection and preview
  • Dry-run mode for safe testing
  • Duplicate handling with auto-renaming
  • Undo support

Examples:
  forter                    # Organize current directory
  forter ~/Downloads        # Organize Downloads folder
  forter --dry-run          # Preview changes without moving files
  forter -r ~/Desktop       # Recursively organize Desktop`,
		Args: cobra.MaximumNArgs(1),
		RunE: run,
	}

	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview changes without moving files")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Scan directories recursively")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Show version information")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if version {
		printVersion()
		return nil
	}

	// Determine target path
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Resolve to absolute path
	if p, err := os.Getwd(); err == nil && path == "." {
		path = p
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Create TUI model
	model, err := tui.NewModel(cfg, path, dryRun, recursive)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Run the TUI
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

func printVersion() {
	fmt.Printf("forter version %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
	fmt.Println("A high-performance TUI file organizer built with Go and Bubble Tea")
	fmt.Println("https://github.com/OmkarVijayBagade/forter")
}
