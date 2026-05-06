// Package tui handles the terminal user interface
package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/OmkarVijayBagade/forter/internal/config"
	"github.com/OmkarVijayBagade/forter/internal/organizer"
	"github.com/OmkarVijayBagade/forter/internal/scanner"
)

// Panel represents the currently active panel
type Panel int

const (
	PanelFiles Panel = iota
	PanelCategories
	PanelPreview
)

// Mode represents the current application mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeConfirm
	ModeProcessing
	ModeDone
	ModeError
)

// Model is the main TUI model
type Model struct {
	// Core components
	config    *config.Config
	organizer *organizer.Organizer
	scanner   *scanner.Scanner

	// Data
	files      []*scanner.FileInfo
	categories []CategoryCount
	operations []organizer.Operation
	logs       []LogEntry

	// State
	width      int
	height     int
	styles     Styles
	panel      Panel
	mode       Mode
	cursor     int // Current cursor position in file list
	catCursor  int // Current cursor position in categories

	// Selection
	selectedFiles map[int]struct{}
	allSelected   bool

	// Paths
	rootPath string

	// Processing
	processing   bool
	stats        organizer.Stats
	dryRun       bool
	recursive    bool

	// Modal
	modalMessage string
	modalAction  func()

	// Filter
	filterCategory string
}

// CategoryCount tracks files per category
type CategoryCount struct {
	Name  string
	Count int
	Size  int64
}

// LogEntry represents a log message
type LogEntry struct {
	Message   string
	Status    string
	Timestamp time.Time
}

// KeyMap defines all keybindings
type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Tab       key.Binding
	Enter     key.Binding
	Space     key.Binding
	SelectAll key.Binding
	Quit      key.Binding
	Help      key.Binding
	Confirm   key.Binding
	Cancel    key.Binding
	Undo      key.Binding
	DryRun    key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left panel"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right panel"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch panel"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/confirm"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle file"),
		),
		SelectAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "select all"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q/esc", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "organize"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "clear selection"),
		),
		Undo: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "undo last"),
		),
		DryRun: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "dry run"),
		),
	}
}

// NewModel creates a new TUI model
func NewModel(cfg *config.Config, rootPath string, dryRun, recursive bool) (*Model, error) {
	// Ensure path exists and is accessible
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", rootPath)
	}

	scan := scanner.NewScanner(cfg, recursive)
	org := organizer.NewOrganizer(cfg, dryRun)

	m := &Model{
		config:        cfg,
		organizer:     org,
		scanner:       scan,
		files:         make([]*scanner.FileInfo, 0),
		categories:    make([]CategoryCount, 0),
		logs:          make([]LogEntry, 0),
		selectedFiles: make(map[int]struct{}),
		rootPath:      rootPath,
		dryRun:        dryRun,
		recursive:     recursive,
		panel:         PanelFiles,
		mode:          ModeNormal,
	}

	// Set up progress callback
	org.SetProgressCallback(func(op organizer.Operation) {
		// This will be called during processing
	})

	return m, nil
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.scanFiles(),
		tea.EnterAltScreen,
	)
}

// scanFiles performs the initial directory scan
func (m *Model) scanFiles() tea.Cmd {
	return func() tea.Msg {
		result, err := m.scanner.Scan(m.rootPath)
		if err != nil {
			return errMsg{err: err}
		}
		return scanCompleteMsg{result: result}
	}
}

// Messages
type errMsg struct {
	err error
}

type scanCompleteMsg struct {
	result *scanner.Result
}

type organizeCompleteMsg struct {
	operations []organizer.Operation
}

type progressMsg struct {
	operation organizer.Operation
}

type tickMsg struct {
	time time.Time
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.styles = NewStyles(m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case errMsg:
		m.mode = ModeError
		m.modalMessage = fmt.Sprintf("Error: %v", msg.err)
		return m, nil

	case scanCompleteMsg:
		m.files = msg.result.Files
		m.updateCategories()
		return m, nil

	case organizeCompleteMsg:
		m.mode = ModeDone
		m.operations = msg.operations
		m.processing = false
		m.addLog(fmt.Sprintf("Organized %d files", len(msg.operations)), "success")
		return m, nil

	case progressMsg:
		if msg.operation.Status == organizer.OpSuccess {
			m.addLog(fmt.Sprintf("Moved: %s", filepath.Base(msg.operation.Source)), "success")
		} else if msg.operation.Status == organizer.OpFailed {
			m.addLog(fmt.Sprintf("Failed: %s - %v", filepath.Base(msg.operation.Source), msg.operation.Error), "error")
		} else if msg.operation.Status == organizer.OpSkipped && m.dryRun {
			m.addLog(fmt.Sprintf("Would move: %s -> %s", filepath.Base(msg.operation.Source), msg.operation.Destination), "pending")
		}
		return m, nil

	case tickMsg:
		return m, nil
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Build layout - logo always at top
	var sections []string

	// Logo (always visible at top)
	sections = append(sections, m.renderLogo())

	// Header
	sections = append(sections, m.renderHeader())

	// Main content (height adjusted for logo)
	mainContent := m.renderMain()
	sections = append(sections, mainContent)

	// Footer
	sections = append(sections, m.renderFooter())

	// Join all sections
	return lipgloss.JoinVertical(lipgloss.Top, sections...)
}

// Helper methods

func (m *Model) updateCategories() {
	catMap := make(map[string]*CategoryCount)

	for _, file := range m.files {
		cat, exists := catMap[file.Category]
		if !exists {
			cat = &CategoryCount{Name: file.Category}
			catMap[file.Category] = cat
		}
		cat.Count++
		cat.Size += file.Size
	}

	// Convert to slice
	m.categories = make([]CategoryCount, 0, len(catMap))
	for _, cat := range catMap {
		m.categories = append(m.categories, *cat)
	}

	// Sort by predefined order, then alphabetically
	order := []string{"Documents", "Images", "Videos", "Audio", "Archives", "Code", "Others"}
	sort.Slice(m.categories, func(i, j int) bool {
		iIdx := -1
		jIdx := -1
		for idx, name := range order {
			if m.categories[i].Name == name {
				iIdx = idx
			}
			if m.categories[j].Name == name {
				jIdx = idx
			}
		}
		if iIdx != -1 && jIdx != -1 {
			return iIdx < jIdx
		}
		if iIdx != -1 {
			return true
		}
		if jIdx != -1 {
			return false
		}
		return strings.Compare(m.categories[i].Name, m.categories[j].Name) < 0
	})
}

func (m *Model) getSelectedFiles() []*scanner.FileInfo {
	selected := make([]*scanner.FileInfo, 0)
	for idx := range m.selectedFiles {
		if idx < len(m.files) {
			file := m.files[idx]
			file.IsSelected = true
			selected = append(selected, file)
		}
	}
	return selected
}

func (m *Model) getFilteredFiles() []*scanner.FileInfo {
	if m.filterCategory == "" {
		return m.files
	}

	filtered := make([]*scanner.FileInfo, 0)
	for _, f := range m.files {
		if f.Category == m.filterCategory {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func (m *Model) addLog(message, status string) {
	entry := LogEntry{
		Message:   message,
		Status:    status,
		Timestamp: time.Now(),
	}
	m.logs = append(m.logs, entry)
	// Keep only last 100 entries
	if len(m.logs) > 100 {
		m.logs = m.logs[len(m.logs)-100:]
	}
}

func (m *Model) selectedCount() int {
	return len(m.selectedFiles)
}

func (m *Model) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
