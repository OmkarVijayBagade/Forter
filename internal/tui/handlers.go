// Package tui handles the terminal user interface
package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/OmkarVijayBagade/forter/internal/organizer"
)

// handleKeyMsg processes keyboard input
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	km := DefaultKeyMap()

	// Global quit
	if key.Matches(msg, km.Quit) {
		if m.mode == ModeConfirm {
			m.mode = ModeNormal
			return m, nil
		}
		return m, tea.Quit
	}

	// Handle modal modes
	switch m.mode {
	case ModeConfirm:
		if key.Matches(msg, km.Enter) {
			m.mode = ModeProcessing
			return m, m.startOrganize()
		}
		if key.Matches(msg, km.Cancel) || key.Matches(msg, km.Quit) {
			m.mode = ModeNormal
			return m, nil
		}
		return m, nil

	case ModeProcessing:
		// Block input during processing
		return m, nil

	case ModeDone:
		if key.Matches(msg, km.Enter) || key.Matches(msg, km.Quit) {
			m.mode = ModeNormal
			// Rescan to refresh
			return m, m.scanFiles()
		}
		return m, nil

	case ModeError:
		m.mode = ModeNormal
		return m, nil
	}

	// Normal mode navigation
	// Panel switching
	if key.Matches(msg, km.Tab) || key.Matches(msg, km.Right) {
		m.switchPanelNext()
		return m, nil
	}

	if key.Matches(msg, km.Left) {
		m.switchPanelPrev()
		return m, nil
	}

	// Actions
	if key.Matches(msg, km.Confirm) {
		if m.selectedCount() > 0 {
			m.mode = ModeConfirm
			if m.dryRun {
				m.modalMessage = fmt.Sprintf("Dry run: Preview %d files?", m.selectedCount())
			} else {
				m.modalMessage = fmt.Sprintf("Organize %d files? (Enter to confirm)", m.selectedCount())
			}
		}
		return m, nil
	}

	if key.Matches(msg, km.SelectAll) {
		m.toggleSelectAll()
		return m, nil
	}

	if key.Matches(msg, km.Cancel) {
		m.clearSelection()
		return m, nil
	}

	if key.Matches(msg, km.DryRun) {
		m.dryRun = !m.dryRun
		m.organizer = organizer.NewOrganizer(m.config, m.dryRun)
		m.addLog(fmt.Sprintf("Dry run mode: %v", m.dryRun), "info")
		return m, nil
	}

	if key.Matches(msg, km.Undo) {
		return m, m.undoLast()
	}

	// Panel-specific navigation
	switch m.panel {
	case PanelFiles:
		return m.handleFilesPanelKey(msg, km)
	case PanelCategories:
		return m.handleCategoriesPanelKey(msg, km)
	}

	return m, nil
}

// handleFilesPanelKey handles keys when file list is focused
func (m *Model) handleFilesPanelKey(msg tea.KeyMsg, km KeyMap) (tea.Model, tea.Cmd) {
	files := m.getFilteredFiles()

	if key.Matches(msg, km.Up) {
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	}

	if key.Matches(msg, km.Down) {
		if m.cursor < len(files)-1 {
			m.cursor++
		}
		return m, nil
	}

	if key.Matches(msg, km.Enter) || key.Matches(msg, km.Space) {
		m.toggleFileSelection()
		return m, nil
	}

	return m, nil
}

// handleCategoriesPanelKey handles keys when categories panel is focused
func (m *Model) handleCategoriesPanelKey(msg tea.KeyMsg, km KeyMap) (tea.Model, tea.Cmd) {
	if key.Matches(msg, km.Up) {
		if m.catCursor > 0 {
			m.catCursor--
		}
		return m, nil
	}

	if key.Matches(msg, km.Down) {
		if m.catCursor < len(m.categories) {
			m.catCursor++
		}
		return m, nil
	}

	if key.Matches(msg, km.Enter) || key.Matches(msg, km.Space) {
		m.toggleCategoryFilter()
		return m, nil
	}

	return m, nil
}

// switchPanelNext switches to the next panel
func (m *Model) switchPanelNext() {
	switch m.panel {
	case PanelFiles:
		m.panel = PanelCategories
	case PanelCategories:
		m.panel = PanelFiles
	}
}

// switchPanelPrev switches to the previous panel
func (m *Model) switchPanelPrev() {
	switch m.panel {
	case PanelFiles:
		m.panel = PanelCategories
	case PanelCategories:
		m.panel = PanelFiles
	}
}

// toggleFileSelection toggles selection of the current file
func (m *Model) toggleFileSelection() {
	files := m.getFilteredFiles()
	if m.cursor >= len(files) {
		return
	}

	// Find actual index in original files slice
	actualIdx := -1
	for i, f := range m.files {
		if f == files[m.cursor] {
			actualIdx = i
			break
		}
	}

	if actualIdx == -1 {
		return
	}

	if _, exists := m.selectedFiles[actualIdx]; exists {
		delete(m.selectedFiles, actualIdx)
		m.files[actualIdx].IsSelected = false
	} else {
		m.selectedFiles[actualIdx] = struct{}{}
		m.files[actualIdx].IsSelected = true
	}
}

// toggleSelectAll selects or deselects all files
func (m *Model) toggleSelectAll() {
	if m.allSelected {
		m.clearSelection()
	} else {
		for i := range m.files {
			m.selectedFiles[i] = struct{}{}
			m.files[i].IsSelected = true
		}
		m.allSelected = true
		m.addLog(fmt.Sprintf("Selected all %d files", len(m.files)), "info")
	}
}

// clearSelection clears all selections
func (m *Model) clearSelection() {
	m.selectedFiles = make(map[int]struct{})
	for _, f := range m.files {
		f.IsSelected = false
	}
	m.allSelected = false
	m.addLog("Cleared selection", "info")
}

// toggleCategoryFilter filters files by the selected category
func (m *Model) toggleCategoryFilter() {
	if m.catCursor >= len(m.categories) {
		m.filterCategory = ""
		m.cursor = 0
		m.addLog("Showing all files", "info")
		return
	}

	cat := m.categories[m.catCursor].Name
	if m.filterCategory == cat {
		m.filterCategory = "" // Clear filter
		m.addLog("Showing all files", "info")
	} else {
		m.filterCategory = cat
		m.cursor = 0
		m.addLog(fmt.Sprintf("Filtered by: %s", cat), "info")
	}
}

// startOrganize begins the organization process
func (m *Model) startOrganize() tea.Cmd {
	return func() tea.Msg {
		selected := m.getSelectedFiles()
		
		// Update progress callback
		m.organizer.SetProgressCallback(func(op organizer.Operation) {
			// Send progress message through program
		})

		if err := m.organizer.Execute(selected, m.rootPath); err != nil {
			return errMsg{err: err}
		}

		ops := m.organizer.GetOperations()
		return organizeCompleteMsg{operations: ops}
	}
}

// undoLast undoes the last organization operation
func (m *Model) undoLast() tea.Cmd {
	return func() tea.Msg {
		if err := m.organizer.Undo(); err != nil {
			m.addLog(fmt.Sprintf("Undo failed: %v", err), "error")
			return errMsg{err: err}
		}
		m.addLog("Undo completed", "success")
		return m.scanFiles()
	}
}
