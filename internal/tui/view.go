// Package tui handles the terminal user interface
package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderHeader renders the application header
func (m *Model) renderHeader() string {
	title := "📁 organize"
	subtitle := fmt.Sprintf(" %s | Files: %d | Selected: %d",
		m.rootPath,
		len(m.files),
		m.selectedCount(),
	)

	if m.dryRun {
		subtitle += " | [DRY RUN]"
	}

	style := m.styles.Header
	width := m.width - lipgloss.Width(title) - lipgloss.Width(subtitle) - 2
	if width < 0 {
		width = 0
	}
	padding := strings.Repeat(" ", width)

	return style.Render(title + padding + subtitle)
}

// renderFooter renders the footer with keybindings
func (m *Model) renderFooter() string {
	km := DefaultKeyMap()

	// Build keybinding help based on current mode
	var bindings []string

	if m.mode == ModeNormal {
		bindings = append(bindings, m.renderKey("↑↓", "navigate"))
		bindings = append(bindings, m.renderKey("tab", "switch panel"))
		bindings = append(bindings, m.renderKey("space/enter", "toggle"))
		bindings = append(bindings, m.renderKey("a", "select all"))
		bindings = append(bindings, m.renderKey("o", "organize"))
		bindings = append(bindings, m.renderKey("d", "dry run"))
		bindings = append(bindings, m.renderKey("u", "undo"))
		bindings = append(bindings, m.renderKey("c", "clear"))
		bindings = append(bindings, m.renderKey("q", "quit"))
	} else if m.mode == ModeConfirm {
		bindings = append(bindings, m.renderKey("enter", "confirm"))
		bindings = append(bindings, m.renderKey("c/esc", "cancel"))
	}

	return m.styles.Footer.Render(strings.Join(bindings, ""))
}

// renderKey renders a keybinding
func (m *Model) renderKey(key, desc string) string {
	return m.styles.KeyBinding.Render(
		m.styles.Key.Render(key) + " " + m.styles.Description.Render(desc),
	)
}

// renderMain renders the main content area
func (m *Model) renderMain() string {
	availableHeight := m.height - 2 // Subtract header and footer
	
	// Calculate panel widths
	totalWidth := m.width - 4 // Account for borders
	leftWidth := int(float64(totalWidth) * 0.6)
	rightWidth := totalWidth - leftWidth

	// Render panels
	leftPanel := m.renderFileList(leftWidth, availableHeight-1)
	rightPanel := m.renderRightPanel(rightWidth, availableHeight-1)

	// Combine side by side
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	// Add status bar at bottom
	statusBar := m.renderStatusBar(m.width - 4)

	return lipgloss.JoinVertical(lipgloss.Top, mainContent, statusBar)
}

// renderFileList renders the file list panel
func (m *Model) renderFileList(width, height int) string {
	files := m.getFilteredFiles()
	
	// Adjust styles to fit
	m.styles.FileList = m.styles.FileList.Width(width).Height(height)
	m.styles.FileItem = m.styles.FileItem.Width(width - 2)
	m.styles.FileSelected = m.styles.FileSelected.Width(width - 2)
	m.styles.FileCursor = m.styles.FileCursor.Width(width - 2)

	var items []string

	// Title
	title := "Files"
	if m.filterCategory != "" {
		title += fmt.Sprintf(" [%s]", m.filterCategory)
	}
	items = append(items, m.styles.PanelTitle.Render(title))

	// File items
	visibleCount := height - 3 // Account for title and border
	startIdx := 0
	if m.cursor > visibleCount/2 && len(files) > visibleCount {
		startIdx = m.cursor - visibleCount/2
		if startIdx > len(files)-visibleCount {
			startIdx = len(files) - visibleCount
		}
	}

	endIdx := startIdx + visibleCount
	if endIdx > len(files) {
		endIdx = len(files)
	}

	for i := startIdx; i < endIdx; i++ {
		file := files[i]
		items = append(items, m.renderFileItem(file, i == m.cursor, i))
	}

	// Fill empty space if needed
	for i := len(items) - 1; i < visibleCount+1; i++ {
		items = append(items, m.styles.FileItem.Render(""))
	}

	panel := lipgloss.JoinVertical(lipgloss.Left, items...)
	
	// Highlight border if focused
	panelStyle := m.styles.Panel
	if m.panel == PanelFiles {
		panelStyle = panelStyle.BorderForeground(ColorPrimary)
	}
	panelStyle = panelStyle.Width(width).Height(height)
	
	return panelStyle.Render(panel)
}

// renderFileItem renders a single file item
func (m *Model) renderFileItem(file *scanner.FileInfo, isCursor, idx int) string {
	// Find actual index
	actualIdx := -1
	for i, f := range m.files {
		if f == file {
			actualIdx = i
			break
		}
	}

	isSelected := false
	if actualIdx >= 0 {
		_, isSelected = m.selectedFiles[actualIdx]
	}

	// Build display
	var prefix string
	if isSelected {
		prefix = "☑ "
	} else {
		prefix = "☐ "
	}

	name := file.Name
	if len(name) > 25 {
		name = name[:22] + "..."
	}

	ext := file.Ext
	if ext == "" {
		ext = "no ext"
	}

	size := m.formatBytes(file.Size)

	// Category color
	catColor := CategoryColor(file.Category)
	category := lipgloss.NewStyle().Foreground(catColor).Render(file.Category)

	// Format: [✓] filename.ext (Size) Category
	line := fmt.Sprintf("%s%-30s %6s %s",
		prefix,
		name,
		size,
		category,
	)

	// Apply style based on state
	if isCursor {
		return m.styles.FileCursor.Render(line)
	}
	if isSelected {
		return m.styles.FileSelected.Render(line)
	}
	return m.styles.FileItem.Render(line)
}

// renderRightPanel renders the right panel (categories + preview)
func (m *Model) renderRightPanel(width, height int) string {
	// Split right panel into categories and preview
	catHeight := int(float64(height) * 0.4)
	previewHeight := height - catHeight - 1

	categories := m.renderCategories(width, catHeight)
	preview := m.renderPreview(width, previewHeight)

	return lipgloss.JoinVertical(lipgloss.Top, categories, preview)
}

// renderCategories renders the categories panel
func (m *Model) renderCategories(width, height int) string {
	m.styles.CategoryList = m.styles.CategoryList.Width(width - 2)
	m.styles.CategoryItem = m.styles.CategoryItem.Width(width - 4)
	m.styles.CategorySelected = m.styles.CategorySelected.Width(width - 4)

	var items []string
	items = append(items, m.styles.PanelTitle.Render("Categories"))

	// All category
	allCount := len(m.files)
	allSize := int64(0)
	for _, f := range m.files {
		allSize += f.Size
	}
	allLine := m.renderCategoryLine("All", allCount, allSize, m.filterCategory == "")
	items = append(items, allLine)

	// Individual categories
	for i, cat := range m.categories {
		isSelected := m.filterCategory == cat.Name
		isCursor := m.catCursor == i && m.panel == PanelCategories
		line := m.renderCategoryLineWithCursor(cat.Name, cat.Count, cat.Size, isSelected, isCursor)
		items = append(items, line)
	}

	// Clear filter option
	isCursor := m.catCursor == len(m.categories) && m.panel == PanelCategories
	clearLine := m.renderCategoryLineWithCursor("Clear Filter", 0, 0, m.filterCategory == "", isCursor)
	items = append(items, clearLine)

	// Fill empty space
	visibleCount := height - 3
	for i := len(items) - 1; i < visibleCount+1; i++ {
		items = append(items, "")
	}

	panel := lipgloss.JoinVertical(lipgloss.Left, items...)
	
	panelStyle := m.styles.Panel
	if m.panel == PanelCategories {
		panelStyle = panelStyle.BorderForeground(ColorPrimary)
	}
	panelStyle = panelStyle.Width(width).Height(height)
	
	return panelStyle.Render(panel)
}

// renderCategoryLine renders a category line
func (m *Model) renderCategoryLine(name string, count int, size int64, isSelected bool) string {
	return m.renderCategoryLineWithCursor(name, count, size, isSelected, false)
}

// renderCategoryLineWithCursor renders a category line with cursor
func (m *Model) renderCategoryLineWithCursor(name string, count int, size int64, isSelected, isCursor bool) string {
	catColor := CategoryColor(name)
	
	var nameStyle lipgloss.Style
	if isCursor {
		nameStyle = m.styles.CategorySelected
	} else if isSelected {
		nameStyle = lipgloss.NewStyle().Foreground(catColor).Bold(true)
	} else {
		nameStyle = lipgloss.NewStyle().Foreground(ColorText)
	}

	countStr := fmt.Sprintf("%d", count)
	sizeStr := m.formatBytes(size)
	
	if name == "Clear Filter" {
		return nameStyle.Render("  [ Clear Filter ]")
	}

	line := fmt.Sprintf("  %s %s (%s)",
		nameStyle.Render(name),
		lipgloss.NewStyle().Foreground(ColorMuted).Render(countStr),
		lipgloss.NewStyle().Foreground(ColorInfo).Render(sizeStr),
	)

	if isSelected {
		line = "▸" + line[1:]
	}

	return m.styles.CategoryItem.Render(line)
}

// renderPreview renders the preview/status panel
func (m *Model) renderPreview(width, height int) string {
	files := m.getFilteredFiles()
	
	m.styles.Preview = m.styles.Preview.Width(width - 2).Height(height - 2)

	var content []string
	content = append(content, m.styles.PanelTitle.Render("Preview"))

	if m.cursor < len(files) {
		file := files[m.cursor]
		targetPath := m.organizer.PreviewDestination(file, m.rootPath)
		
		info := []string{
			fmt.Sprintf("Name: %s", file.Name),
			fmt.Sprintf("Size: %s", m.formatBytes(file.Size)),
			fmt.Sprintf("Ext:  %s", file.Ext),
			fmt.Sprintf("Cat:  %s", file.Category),
			"",
			"Target:",
			filepath.Dir(targetPath),
		}
		
		content = append(content, strings.Join(info, "\n"))
	} else if len(files) == 0 {
		content = append(content, "No files found")
	} else {
		content = append(content, "Select a file to preview")
	}

	// Recent activity
	if len(m.logs) > 0 {
		content = append(content, "")
		content = append(content, m.styles.PanelTitle.Render("Activity"))
		
		logCount := 3
		if len(m.logs) < logCount {
			logCount = len(m.logs)
		}
		
		startIdx := len(m.logs) - logCount
		for i := len(m.logs) - 1; i >= startIdx; i-- {
			log := m.logs[i]
			color := StatusColor(log.Status)
			entry := lipgloss.NewStyle().Foreground(color).Render("• " + log.Message)
			content = append(content, entry)
		}
	}

	panel := lipgloss.JoinVertical(lipgloss.Left, content...)
	return m.styles.Preview.Render(panel)
}

// renderStatusBar renders the status bar at the bottom
func (m *Model) renderStatusBar(width int) string {
	m.styles.StatusBar = m.styles.StatusBar.Width(width)

	var status string
	
	if m.mode == ModeProcessing {
		status = "Processing..."
	} else if m.mode == ModeConfirm {
		status = m.modalMessage
	} else if m.mode == ModeDone {
		total, success, failed, skipped := m.organizer.GetStats()
		status = fmt.Sprintf("Complete: %d total, %d success, %d failed, %d skipped", 
			total, success, failed, skipped)
	} else if m.mode == ModeError {
		status = m.modalMessage
	} else {
		// Show selected summary
		if m.selectedCount() > 0 {
			totalSize := int64(0)
			for idx := range m.selectedFiles {
				if idx < len(m.files) {
					totalSize += m.files[idx].Size
				}
			}
			status = fmt.Sprintf("Selected: %d files (%s)", m.selectedCount(), m.formatBytes(totalSize))
		} else {
			status = fmt.Sprintf("Ready | %d files (%s)", len(m.files), m.formatBytes(m.getTotalSize()))
		}
	}

	return m.styles.StatusBar.Render(status)
}

// getTotalSize returns total size of all files
func (m *Model) getTotalSize() int64 {
	var total int64
	for _, f := range m.files {
		total += f.Size
	}
	return total
}
