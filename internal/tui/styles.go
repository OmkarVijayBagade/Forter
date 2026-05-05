// Package tui handles the terminal user interface
package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color scheme - inspired by modern tools like lazygit
var (
	ColorPrimary   = lipgloss.Color("#7aa2f7") // Blue
	ColorSecondary = lipgloss.Color("#bb9af7") // Purple
	ColorSuccess   = lipgloss.Color("#9ece6a") // Green
	ColorWarning   = lipgloss.Color("#e0af68") // Yellow
	ColorError     = lipgloss.Color("#f7768e") // Red
	ColorInfo      = lipgloss.Color("#73daca") // Cyan
	ColorText      = lipgloss.Color("#c0caf5") // Light text
	ColorMuted     = lipgloss.Color("#565f89") // Dark text
	ColorBg        = lipgloss.Color("#1a1b26") // Background
	ColorSurface   = lipgloss.Color("#24283b") // Surface
	ColorBorder    = lipgloss.Color("#414868") // Border
)

// Styles defines all UI styles
type Styles struct {
	// Layout
	App        lipgloss.Style
	Header     lipgloss.Style
	Footer     lipgloss.Style
	Main       lipgloss.Style
	Panel      lipgloss.Style
	PanelTitle lipgloss.Style

	// File list
	FileList       lipgloss.Style
	FileItem       lipgloss.Style
	FileSelected   lipgloss.Style
	FileDeselected lipgloss.Style
	FileCursor     lipgloss.Style
	FileCategory   lipgloss.Style

	// Categories
	CategoryList     lipgloss.Style
	CategoryItem     lipgloss.Style
	CategorySelected lipgloss.Style
	CategoryCount    lipgloss.Style

	// Preview/Status
	Preview    lipgloss.Style
	StatusBar  lipgloss.Style
	StatusItem lipgloss.Style
	LogEntry   lipgloss.Style

	// Keybindings
	KeyBinding  lipgloss.Style
	Key         lipgloss.Style
	Description lipgloss.Style

	// Modal
	Modal       lipgloss.Style
	ModalTitle  lipgloss.Style
	ModalButton lipgloss.Style
}

// NewStyles creates all UI styles
func NewStyles(width, height int) Styles {
	return Styles{
		// Layout
		App: lipgloss.NewStyle().
			Width(width).
			Height(height).
			Background(ColorBg),

		Header: lipgloss.NewStyle().
			Width(width).
			Height(1).
			Background(ColorSurface).
			Foreground(ColorPrimary).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1),

		Footer: lipgloss.NewStyle().
			Width(width).
			Height(1).
			Background(ColorSurface).
			Foreground(ColorMuted).
			PaddingLeft(1).
			PaddingRight(1),

		Main: lipgloss.NewStyle().
			Width(width).
			Height(height - 2).
			Background(ColorBg),

		Panel: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Background(ColorSurface).
			Padding(1),

		PanelTitle: lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			MarginBottom(1),

		// File list
		FileList: lipgloss.NewStyle().
			Background(ColorSurface),

		FileItem: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Height(1),

		FileSelected: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Foreground(ColorSuccess).
			Bold(true),

		FileDeselected: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Foreground(ColorMuted),

		FileCursor: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Background(ColorPrimary).
			Foreground(ColorBg).
			Bold(true),

		FileCategory: lipgloss.NewStyle().
			Foreground(ColorSecondary),

		// Categories
		CategoryList: lipgloss.NewStyle().
			Background(ColorSurface),

		CategoryItem: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Height(1),

		CategorySelected: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Height(1).
			Background(ColorPrimary).
			Foreground(ColorBg).
			Bold(true),

		CategoryCount: lipgloss.NewStyle().
			Foreground(ColorWarning).
			Align(lipgloss.Right),

		// Preview/Status
		Preview: lipgloss.NewStyle().
			Background(ColorSurface).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder),

		StatusBar: lipgloss.NewStyle().
			Background(ColorSurface).
			Foreground(ColorText).
			Height(5).
			Padding(1),

		StatusItem: lipgloss.NewStyle().
			PaddingLeft(1),

		LogEntry: lipgloss.NewStyle().
			Height(1).
			PaddingLeft(1),

		// Keybindings
		KeyBinding: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(2),

		Key: lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true),

		Description: lipgloss.NewStyle().
			Foreground(ColorMuted),

		// Modal
		Modal: lipgloss.NewStyle().
			Background(ColorSurface).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(2).
			Width(60),

		ModalTitle: lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			MarginBottom(1).
			Align(lipgloss.Center),

		ModalButton: lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorBg).
			Padding(0, 2).
			MarginTop(1),
	}
}

// Helper functions for status colors
func StatusColor(status string) lipgloss.Color {
	switch status {
	case "success", "moved":
		return ColorSuccess
	case "error", "failed":
		return ColorError
	case "skipped", "pending":
		return ColorWarning
	default:
		return ColorMuted
	}
}

// CategoryColor returns a color for a category
func CategoryColor(category string) lipgloss.Color {
	switch category {
	case "Documents":
		return lipgloss.Color("#7aa2f7")
	case "Images":
		return lipgloss.Color("#bb9af7")
	case "Videos":
		return lipgloss.Color("#f7768e")
	case "Audio":
		return lipgloss.Color("#9ece6a")
	case "Archives":
		return lipgloss.Color("#e0af68")
	case "Code":
		return lipgloss.Color("#73daca")
	default:
		return ColorMuted
	}
}
