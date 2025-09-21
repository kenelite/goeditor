package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StatusBar represents the status bar component at the bottom of the editor
type StatusBar struct {
	*fyne.Container
	
	// Status bar components
	positionLabel  *widget.Label  // Shows cursor position (line:column)
	fileInfoLabel  *widget.Label  // Shows file path and size
	modifiedLabel  *widget.Label  // Shows modification status
	encodingLabel  *widget.Label  // Shows file encoding
	languageLabel  *widget.Label  // Shows file language/type
	selectionLabel *widget.Label  // Shows selection statistics
	
	// Reference to editor for state access
	editor *Editor
}

// StatusBarInterface defines the interface for status bar operations
type StatusBarInterface interface {
	UpdatePosition(line, col int)
	UpdateFileInfo(path string, size int64)
	SetModified(modified bool)
	SetLanguage(language string)
	SetEncoding(encoding string)
	UpdateSelection(hasSelection bool, selectedText string)
	Refresh()
}

// NewStatusBar creates a new status bar component
func NewStatusBar(editor *Editor) *StatusBar {
	sb := &StatusBar{
		editor: editor,
	}
	
	// Initialize labels
	sb.positionLabel = widget.NewLabel("Ln 1, Col 1")
	sb.fileInfoLabel = widget.NewLabel("Ready")
	sb.modifiedLabel = widget.NewLabel("")
	sb.encodingLabel = widget.NewLabel("UTF-8")
	sb.languageLabel = widget.NewLabel("Plain Text")
	sb.selectionLabel = widget.NewLabel("")
	
	// Style the labels
	sb.styleLabels()
	
	// Create the container layout
	sb.createLayout()
	
	// Set up editor callbacks to update status bar
	sb.setupCallbacks()
	
	return sb
}

// styleLabels applies consistent styling to all status bar labels
func (sb *StatusBar) styleLabels() {
	// Set text style for all labels
	labels := []*widget.Label{
		sb.positionLabel,
		sb.fileInfoLabel,
		sb.modifiedLabel,
		sb.encodingLabel,
		sb.languageLabel,
		sb.selectionLabel,
	}
	
	for _, label := range labels {
		label.TextStyle = fyne.TextStyle{}
		label.Alignment = fyne.TextAlignLeading
	}
	
	// Special styling for modified indicator
	sb.modifiedLabel.TextStyle = fyne.TextStyle{Bold: true}
}

// createLayout creates the status bar layout with all components
func (sb *StatusBar) createLayout() {
	// Create separators
	separator1 := widget.NewSeparator()
	separator2 := widget.NewSeparator()
	
	// Create left section (file info and modified status)
	leftSection := container.NewHBox(
		sb.fileInfoLabel,
		sb.modifiedLabel,
	)
	
	// Create center section (selection info)
	centerSection := container.NewHBox(
		sb.selectionLabel,
	)
	
	// Create right section (position, encoding, language)
	rightSection := container.NewHBox(
		sb.positionLabel,
		separator1,
		sb.encodingLabel,
		separator2,
		sb.languageLabel,
	)
	
	// Create the main container using border layout
	sb.Container = container.NewBorder(
		nil, nil, // top, bottom
		leftSection, rightSection, // left, right
		centerSection, // center
	)
}

// setupCallbacks sets up callbacks to automatically update the status bar
func (sb *StatusBar) setupCallbacks() {
	if sb.editor == nil {
		return
	}
	
	// Store original callbacks
	originalOnFileChanged := sb.editor.OnFileChanged
	originalOnModified := sb.editor.OnModified
	originalOnCursorChanged := sb.editor.OnCursorChanged
	originalOnSelectionChanged := sb.editor.OnSelectionChanged
	
	// Set up file change callback
	sb.editor.OnFileChanged = func(path string) {
		if originalOnFileChanged != nil {
			originalOnFileChanged(path)
		}
		sb.updateFromEditorState()
	}
	
	// Set up modified callback
	sb.editor.OnModified = func(modified bool) {
		if originalOnModified != nil {
			originalOnModified(modified)
		}
		sb.SetModified(modified)
	}
	
	// Set up cursor change callback
	sb.editor.OnCursorChanged = func(line, col int) {
		if originalOnCursorChanged != nil {
			originalOnCursorChanged(line, col)
		}
		sb.UpdatePosition(line, col)
	}
	
	// Set up selection change callback
	sb.editor.OnSelectionChanged = func(hasSelection bool) {
		if originalOnSelectionChanged != nil {
			originalOnSelectionChanged(hasSelection)
		}
		sb.updateSelectionFromEditor()
	}
}

// UpdatePosition updates the cursor position display
func (sb *StatusBar) UpdatePosition(line, col int) {
	sb.positionLabel.SetText(fmt.Sprintf("Ln %d, Col %d", line, col))
}

// UpdateFileInfo updates the file information display
func (sb *StatusBar) UpdateFileInfo(path string, size int64) {
	if path == "" {
		sb.fileInfoLabel.SetText("Untitled")
		return
	}
	
	filename := filepath.Base(path)
	sizeStr := formatFileSize(size)
	sb.fileInfoLabel.SetText(fmt.Sprintf("%s (%s)", filename, sizeStr))
}

// SetModified updates the modification status display
func (sb *StatusBar) SetModified(modified bool) {
	if modified {
		sb.modifiedLabel.SetText("●")
	} else {
		sb.modifiedLabel.SetText("")
	}
}

// SetLanguage updates the language/file type display
func (sb *StatusBar) SetLanguage(language string) {
	if language == "" {
		language = "Plain Text"
	}
	sb.languageLabel.SetText(language)
}

// SetEncoding updates the encoding display
func (sb *StatusBar) SetEncoding(encoding string) {
	if encoding == "" {
		encoding = "UTF-8"
	}
	sb.encodingLabel.SetText(encoding)
}

// UpdateSelection updates the selection statistics display
func (sb *StatusBar) UpdateSelection(hasSelection bool, selectedText string) {
	if !hasSelection || selectedText == "" {
		sb.selectionLabel.SetText("")
		return
	}
	
	// Count characters and lines in selection
	charCount := len(selectedText)
	lineCount := strings.Count(selectedText, "\n") + 1
	
	if lineCount > 1 {
		sb.selectionLabel.SetText(fmt.Sprintf("(%d chars, %d lines selected)", charCount, lineCount))
	} else {
		sb.selectionLabel.SetText(fmt.Sprintf("(%d chars selected)", charCount))
	}
}

// Refresh updates all status bar information from the editor state
func (sb *StatusBar) Refresh() {
	sb.updateFromEditorState()
}

// updateFromEditorState updates all status bar components from the current editor state
func (sb *StatusBar) updateFromEditorState() {
	if sb.editor == nil || sb.editor.State == nil {
		return
	}
	
	state := sb.editor.State
	
	// Update position
	sb.UpdatePosition(state.CursorLine, state.CursorColumn)
	
	// Update file info
	sb.UpdateFileInfo(state.CurrentFile, state.FileSize)
	
	// Update modified status
	sb.SetModified(state.IsModified)
	
	// Update language
	sb.SetLanguage(state.Language)
	
	// Update encoding
	sb.SetEncoding(state.Encoding)
	
	// Update selection
	sb.updateSelectionFromEditor()
}

// updateSelectionFromEditor updates selection info from the editor
func (sb *StatusBar) updateSelectionFromEditor() {
	if sb.editor == nil || sb.editor.State == nil {
		sb.UpdateSelection(false, "")
		return
	}
	
	state := sb.editor.State
	hasSelection := state.HasSelection()
	
	if hasSelection {
		// In a real implementation, we would get the actual selected text
		// For now, we'll simulate it based on selection positions
		selectedText := sb.getSelectedTextFromEditor()
		sb.UpdateSelection(true, selectedText)
	} else {
		sb.UpdateSelection(false, "")
	}
}

// getSelectedTextFromEditor gets the currently selected text from the editor
func (sb *StatusBar) getSelectedTextFromEditor() string {
	if sb.editor == nil {
		return ""
	}
	
	// In a real implementation with proper cursor/selection tracking,
	// we would extract the actual selected text from the editor content
	// For now, we'll return a placeholder
	
	// This is a simplified implementation - in reality we would need
	// to track selection start/end positions and extract the text between them
	content := sb.editor.GetContent()
	if content == "" {
		return ""
	}
	
	// For demonstration, assume some text is selected if the state indicates selection
	if sb.editor.State != nil && sb.editor.State.HasSelection() {
		// Return a sample selection for demonstration
		lines := strings.Split(content, "\n")
		if len(lines) > 0 {
			return lines[0] // Return first line as "selected" text for demo
		}
	}
	
	return ""
}

// formatFileSize formats a file size in bytes to a human-readable string
func formatFileSize(size int64) string {
	if size == 0 {
		return "0 bytes"
	}
	
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d bytes", size)
	}
	
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(size)/float64(div), units[exp])
}

// GetContainer returns the status bar container for embedding in the main window
func (sb *StatusBar) GetContainer() *fyne.Container {
	return sb.Container
}

// Hide hides the status bar
func (sb *StatusBar) Hide() {
	sb.Container.Hide()
}

// Show shows the status bar
func (sb *StatusBar) Show() {
	sb.Container.Show()
}

// IsVisible returns whether the status bar is currently visible
func (sb *StatusBar) IsVisible() bool {
	return sb.Container.Visible()
}