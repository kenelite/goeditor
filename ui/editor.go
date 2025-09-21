package ui

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kenelite/goeditor/backend"
	"github.com/kenelite/goeditor/ui/syntax"
	"github.com/kenelite/goeditor/ui/dialogs"
	"os"
	"strings"
)

// Editor represents the main text editor component
type Editor struct {
	TextWidget         *widget.Entry
	LineNumberWidget   *LineNumberWidget
	ScrollContainer    *container.Scroll
	EditorContainer    *fyne.Container
	StatusBar          *StatusBar
	State              *backend.EditorState
	FileManager        *backend.FileManager
	ConfigManager      *backend.ConfigManager
	History            *backend.History
	SearchManager      *backend.SearchManager
	IndentationManager *IndentationManager
	
	// Dialogs
	FindDialog     *dialogs.FindDialog
	ReplaceDialog  *dialogs.ReplaceDialog
	GoToLineDialog *dialogs.GoToLineDialog
	
	// Callbacks for state changes
	OnFileChanged      func(path string)
	OnModified         func(modified bool)
	OnCursorChanged    func(line, col int)
	OnSelectionChanged func(hasSelection bool)
}

// NewEditor creates a new editor instance
func NewEditor() *Editor {
	e := &Editor{
		TextWidget:         widget.NewMultiLineEntry(),
		State:              backend.NewEditorState(),
		FileManager:        backend.NewFileManager(),
		ConfigManager:      backend.NewConfigManager(),
		History:            backend.NewHistory(),
		SearchManager:      backend.NewSearchManager(),
		IndentationManager: NewIndentationManager(),
	}
	
	// Create line number widget (but don't add to UI yet to avoid crashes)
	e.LineNumberWidget = NewLineNumberWidget(e)
	
	// Create status bar
	e.StatusBar = NewStatusBar(e)
	
	// Create editor container
	e.createEditorContainer()
	
	// Load configuration
	if err := e.ConfigManager.Load(); err != nil {
		// Log error but continue with defaults
		// TODO: Add proper logging
	}
	
	// Set up text widget callbacks
	e.setupTextWidgetCallbacks()
	
	return e
}

// createEditorContainer creates the container with line numbers and text editor
func (e *Editor) createEditorContainer() {
	// For now, create a simple container with just the text widget
	// Line numbers will be added in a future iteration when Fyne supports it better
	e.EditorContainer = container.NewMax(e.TextWidget)
	
	// Wrap in scroll container
	e.ScrollContainer = container.NewScroll(e.EditorContainer)
}

// GetEditorContainer returns the main editor container for embedding in the UI
func (e *Editor) GetEditorContainer() *container.Scroll {
	return e.ScrollContainer
}

// GetCompleteLayout returns the complete editor layout including status bar
func (e *Editor) GetCompleteLayout() *fyne.Container {
	return container.NewBorder(
		nil, e.StatusBar.GetContainer(), // top, bottom
		nil, nil, // left, right
		e.ScrollContainer, // center
	)
}

// GetStatusBar returns the status bar component
func (e *Editor) GetStatusBar() *StatusBar {
	return e.StatusBar
}

// EnableLineNumbers enables line number display (call after UI is initialized)
func (e *Editor) EnableLineNumbers() {
	if e.LineNumberWidget != nil && e.EditorContainer != nil {
		// Recreate container with line numbers
		e.EditorContainer.Objects = []fyne.CanvasObject{
			container.NewHBox(
				e.LineNumberWidget,
				widget.NewSeparator(),
				e.TextWidget,
			),
		}
		e.EditorContainer.Refresh()
	}
}

// InitializeDialogs initializes the search dialogs (call this after window is available)
func (e *Editor) InitializeDialogs(window fyne.Window) {
	e.FindDialog = dialogs.NewFindDialog(e, e.SearchManager, window)
	e.ReplaceDialog = dialogs.NewReplaceDialog(e, e.SearchManager, window)
	e.GoToLineDialog = dialogs.NewGoToLineDialog(e, window)
}

// setupTextWidgetCallbacks sets up callbacks for the text widget
func (e *Editor) setupTextWidgetCallbacks() {
	// Track text changes for modified state and line count
	e.TextWidget.OnChanged = func(content string) {
		if !e.State.IsModified {
			e.State.SetModified(true)
			if e.OnModified != nil {
				e.OnModified(true)
			}
		}
		
		// Update line number widget
		e.updateLineNumbers()
		
		// Simulate cursor position update
		e.SimulateCursorMovement()
		
		// Update status bar
		if e.StatusBar != nil {
			e.StatusBar.Refresh()
		}
	}
	
	// Handle key events for indentation
	e.TextWidget.OnSubmitted = func(content string) {
		// This is called when Enter is pressed
		e.handleEnterKey()
		
		// Update cursor position in status bar (simulate cursor movement)
		if e.OnCursorChanged != nil {
			lines := strings.Split(content, "\n")
			e.OnCursorChanged(len(lines), 1) // Move to start of new line
		}
	}
	
	// TODO: Add cursor position tracking when Fyne supports it
	// TODO: Add selection tracking when Fyne supports it
}

// updateLineNumbers updates the line number widget based on current content
func (e *Editor) updateLineNumbers() {
	if e.LineNumberWidget != nil {
		e.LineNumberWidget.UpdateFromEditor()
	}
}

// handleEnterKey handles Enter key press for auto-indentation
func (e *Editor) handleEnterKey() {
	if !e.IndentationManager.GetAutoIndent() {
		return
	}
	
	content := e.GetContent()
	// For now, we'll implement basic auto-indentation
	// In a full implementation, we would need cursor position tracking
	lines := strings.Split(content, "\n")
	if len(lines) > 1 {
		lastLine := lines[len(lines)-2] // Previous line before the new one
		indentation := e.IndentationManager.GetLineIndentation(lastLine)
		
		// Add extra indentation for certain patterns
		trimmedLine := strings.TrimSpace(lastLine)
		if strings.HasSuffix(trimmedLine, "{") || strings.HasSuffix(trimmedLine, ":") {
			indentation += e.IndentationManager.GetIndentString()
		}
		
		// Update the current line with proper indentation
		if indentation != "" {
			lines[len(lines)-1] = indentation + strings.TrimLeft(lines[len(lines)-1], " \t")
			newContent := strings.Join(lines, "\n")
			e.SetContent(newContent)
		}
	}
}

// LoadFile loads a file into the editor
func (e *Editor) LoadFile(path string) error {
	content, fileInfo, err := e.FileManager.ReadFileWithInfo(path)
	if err != nil {
		return err
	}
	
	// Update editor content
	e.TextWidget.SetText(content)
	
	// Clear history when loading a new file
	e.History.Clear()
	
	// Update state
	fileType := e.FileManager.GetFileType(path)
	e.State.SetCurrentFile(path, fileInfo.Size, fileType.Name)
	e.State.SetModified(false)
	
	// Reset cursor position
	e.State.SetCursorPosition(1, 1)
	
	// Notify callbacks
	if e.OnFileChanged != nil {
		e.OnFileChanged(path)
	}
	if e.OnModified != nil {
		e.OnModified(false)
	}
	if e.OnCursorChanged != nil {
		e.OnCursorChanged(1, 1)
	}
	
	// Update status bar
	if e.StatusBar != nil {
		e.StatusBar.Refresh()
	}
	
	return nil
}

// SaveFile saves the current content to a file
func (e *Editor) SaveFile(path string) error {
	content := e.TextWidget.Text
	
	if err := e.FileManager.SaveFileWithBackup(path, content); err != nil {
		return err
	}
	
	// Update state
	fileInfo, _ := e.FileManager.GetFileInfo(path)
	if fileInfo != nil {
		fileType := e.FileManager.GetFileType(path)
		e.State.SetCurrentFile(path, fileInfo.Size, fileType.Name)
	}
	e.State.SetModified(false)
	
	// Notify callbacks
	if e.OnModified != nil {
		e.OnModified(false)
	}
	if e.OnFileChanged != nil {
		e.OnFileChanged(path)
	}
	
	// Update status bar
	if e.StatusBar != nil {
		e.StatusBar.Refresh()
	}
	
	return nil
}

// NewFile creates a new empty file
func (e *Editor) NewFile() {
	e.TextWidget.SetText("")
	e.State = backend.NewEditorState()
	
	// Clear history when creating a new file
	e.History.Clear()
	
	// Notify callbacks
	if e.OnFileChanged != nil {
		e.OnFileChanged("")
	}
	if e.OnModified != nil {
		e.OnModified(false)
	}
	if e.OnCursorChanged != nil {
		e.OnCursorChanged(1, 1)
	}
	
	// Update status bar
	if e.StatusBar != nil {
		e.StatusBar.Refresh()
	}
}

// GetContent returns the current editor content
func (e *Editor) GetContent() string {
	return e.TextWidget.Text
}

// SetContent sets the editor content
func (e *Editor) SetContent(content string) {
	e.TextWidget.SetText(content)
}

// IsModified returns whether the file has been modified
func (e *Editor) IsModified() bool {
	return e.State.IsModified
}

// GetCurrentFile returns the current file path
func (e *Editor) GetCurrentFile() string {
	return e.State.CurrentFile
}

// GetFileType returns the current file type
func (e *Editor) GetFileType() backend.FileType {
	if e.State.CurrentFile == "" {
		return backend.FileType{Name: "Plain Text", LexerName: "text"}
	}
	return e.FileManager.GetFileType(e.State.CurrentFile)
}

// ApplyConfiguration applies configuration settings to the editor
func (e *Editor) ApplyConfiguration() {
	config := e.ConfigManager.GetEditorConfig()
	
	// Apply font size (when Fyne supports it)
	// TODO: Implement font size setting
	
	// Apply word wrap
	e.TextWidget.Wrapping = fyne.TextWrapWord
	if !config.WordWrap {
		e.TextWidget.Wrapping = fyne.TextWrapOff
	}
	
	// TODO: Apply other configuration settings as Fyne supports them
}

// Undo undoes the last operation
func (e *Editor) Undo() bool {
	if !e.History.CanUndo() {
		return false
	}

	op := e.History.Undo()
	if op == nil {
		return false
	}

	// Apply the reverse operation
	e.applyReverseOperation(*op)
	return true
}

// Redo redoes the last undone operation
func (e *Editor) Redo() bool {
	if !e.History.CanRedo() {
		return false
	}

	op := e.History.Redo()
	if op == nil {
		return false
	}

	// Apply the operation
	e.applyOperation(*op)
	return true
}

// CanUndo returns true if undo is available
func (e *Editor) CanUndo() bool {
	return e.History.CanUndo()
}

// CanRedo returns true if redo is available
func (e *Editor) CanRedo() bool {
	return e.History.CanRedo()
}

// InsertText inserts text at the current cursor position and records the operation
func (e *Editor) InsertText(text string) {
	if text == "" {
		return
	}

	// Get current cursor position (simplified - in real implementation would need proper cursor tracking)
	currentPos := backend.Position{Line: 1, Column: 1} // TODO: Get actual cursor position
	
	// Record the operation
	e.History.RecordInsert(currentPos, text)
	
	// Insert the text
	currentText := e.TextWidget.Text
	e.TextWidget.SetText(currentText + text)
	
	// Mark as modified
	e.State.SetModified(true)
	if e.OnModified != nil {
		e.OnModified(true)
	}
}

// DeleteText deletes text and records the operation
func (e *Editor) DeleteText(position backend.Position, length int) {
	content := e.TextWidget.Text
	if len(content) == 0 {
		return
	}

	// For simplification, delete from end (in real implementation would use position)
	if length > len(content) {
		length = len(content)
	}
	
	deletedText := content[len(content)-length:]
	newContent := content[:len(content)-length]
	
	// Record the operation
	e.History.RecordDelete(position, deletedText)
	
	// Apply the deletion
	e.TextWidget.SetText(newContent)
	
	// Mark as modified
	e.State.SetModified(true)
	if e.OnModified != nil {
		e.OnModified(true)
	}
}

// ReplaceText replaces text and records the operation
func (e *Editor) ReplaceText(position backend.Position, oldText, newText string) {
	// Record the operation
	e.History.RecordReplace(position, oldText, newText)
	
	// Apply the replacement (simplified implementation)
	content := e.TextWidget.Text
	newContent := content // In real implementation, would perform actual replacement
	if oldText != "" {
		// Simple replacement for demonstration
		newContent = content + newText
	}
	
	e.TextWidget.SetText(newContent)
	
	// Mark as modified
	e.State.SetModified(true)
	if e.OnModified != nil {
		e.OnModified(true)
	}
}

// applyOperation applies an operation to the editor content
func (e *Editor) applyOperation(op backend.Operation) {
	// Temporarily disable history recording to avoid recording undo/redo operations
	e.History.SetEnabled(false)
	defer e.History.SetEnabled(true)

	switch op.Type {
	case backend.Insert:
		// Insert the new text
		content := e.TextWidget.Text
		e.TextWidget.SetText(content + op.NewText) // Simplified implementation
	case backend.Delete:
		// Remove the text that was deleted
		content := e.TextWidget.Text
		if len(content) >= len(op.OldText) {
			newContent := content[:len(content)-len(op.OldText)]
			e.TextWidget.SetText(newContent)
		}
	case backend.Replace:
		// Replace with new text
		content := e.TextWidget.Text
		e.TextWidget.SetText(content + op.NewText) // Simplified implementation
	}
}

// applyReverseOperation applies the reverse of an operation
func (e *Editor) applyReverseOperation(op backend.Operation) {
	// Temporarily disable history recording to avoid recording undo/redo operations
	e.History.SetEnabled(false)
	defer e.History.SetEnabled(true)

	switch op.Type {
	case backend.Insert:
		// Remove the inserted text
		content := e.TextWidget.Text
		if len(content) >= len(op.NewText) {
			newContent := content[:len(content)-len(op.NewText)]
			e.TextWidget.SetText(newContent)
		}
	case backend.Delete:
		// Restore the deleted text
		content := e.TextWidget.Text
		e.TextWidget.SetText(content + op.OldText) // Simplified implementation
	case backend.Replace:
		// Replace with old text
		content := e.TextWidget.Text
		e.TextWidget.SetText(content + op.OldText) // Simplified implementation
	}
}

// ClearHistory clears the undo/redo history
func (e *Editor) ClearHistory() {
	e.History.Clear()
}

// Search and Replace Methods

// ShowFindDialog shows the find dialog
func (e *Editor) ShowFindDialog() {
	if e.FindDialog != nil {
		e.FindDialog.Show()
	}
}

// ShowReplaceDialog shows the find and replace dialog
func (e *Editor) ShowReplaceDialog() {
	if e.ReplaceDialog != nil {
		e.ReplaceDialog.Show()
	}
}

// HideFindDialog hides the find dialog
func (e *Editor) HideFindDialog() {
	if e.FindDialog != nil {
		e.FindDialog.Hide()
	}
}

// HideReplaceDialog hides the replace dialog
func (e *Editor) HideReplaceDialog() {
	if e.ReplaceDialog != nil {
		e.ReplaceDialog.Hide()
	}
}

// FindNext finds the next occurrence of the current search pattern
func (e *Editor) FindNext() bool {
	if e.FindDialog != nil && e.FindDialog.IsVisible() {
		return e.FindDialog.FindNext()
	}
	if e.ReplaceDialog != nil && e.ReplaceDialog.IsVisible() {
		return e.ReplaceDialog.FindNext()
	}
	return false
}

// FindPrevious finds the previous occurrence of the current search pattern
func (e *Editor) FindPrevious() bool {
	if e.FindDialog != nil && e.FindDialog.IsVisible() {
		return e.FindDialog.FindPrevious()
	}
	if e.ReplaceDialog != nil && e.ReplaceDialog.IsVisible() {
		return e.ReplaceDialog.FindPrevious()
	}
	return false
}

// Find searches for a pattern in the editor content
func (e *Editor) Find(pattern string, options backend.SearchOptions) []backend.Match {
	e.SearchManager.SetOptions(options)
	return e.SearchManager.Find(e.GetContent(), pattern)
}

// Replace performs text replacement in the editor
func (e *Editor) Replace(pattern, replacement string, options backend.ReplaceOptions) int {
	content := e.GetContent()
	newContent, count := e.SearchManager.Replace(content, pattern, replacement, options)
	
	if count > 0 {
		e.SetContent(newContent)
		// Mark as modified
		e.State.SetModified(true)
		if e.OnModified != nil {
			e.OnModified(true)
		}
	}
	
	return count
}

// GetSearchManager returns the search manager
func (e *Editor) GetSearchManager() *backend.SearchManager {
	return e.SearchManager
}

// Go to Line Methods

// ShowGoToLineDialog shows the "Go to Line" dialog
func (e *Editor) ShowGoToLineDialog() {
	if e.GoToLineDialog != nil {
		e.GoToLineDialog.Show()
	}
}

// HideGoToLineDialog hides the "Go to Line" dialog
func (e *Editor) HideGoToLineDialog() {
	if e.GoToLineDialog != nil {
		e.GoToLineDialog.Hide()
	}
}

// GoToLine navigates to the specified line number
func (e *Editor) GoToLine(lineNumber int) bool {
	content := e.GetContent()
	lines := strings.Split(content, "\n")
	
	if lineNumber < 1 || lineNumber > len(lines) {
		return false
	}
	
	// Update state
	e.State.SetCursorPosition(lineNumber, 1)
	
	// Scroll line number widget to the line
	if e.LineNumberWidget != nil {
		e.LineNumberWidget.ScrollToLine(lineNumber)
	}
	
	// Notify callback
	if e.OnCursorChanged != nil {
		e.OnCursorChanged(lineNumber, 1)
	}
	
	return true
}

// Indentation Methods

// HandleTabKey handles Tab key press for indentation
func (e *Editor) HandleTabKey() {
	content := e.GetContent()
	// For simplified implementation, add indentation at the end
	// In a real implementation, we would need cursor position
	newContent := content + e.IndentationManager.GetIndentString()
	e.SetContent(newContent)
}

// HandleShiftTabKey handles Shift+Tab key press for unindentation
func (e *Editor) HandleShiftTabKey() {
	content := e.GetContent()
	lines := strings.Split(content, "\n")
	
	if len(lines) > 0 {
		// Remove indentation from the last line (simplified)
		lastLineIndex := len(lines) - 1
		newLine, _ := e.IndentationManager.removeIndentation(lines[lastLineIndex])
		lines[lastLineIndex] = newLine
		
		newContent := strings.Join(lines, "\n")
		e.SetContent(newContent)
	}
}

// IndentSelectedLines indents the currently selected lines
func (e *Editor) IndentSelectedLines() {
	content := e.GetContent()
	lines := strings.Split(content, "\n")
	
	// For simplified implementation, indent all lines
	// In a real implementation, we would only indent selected lines
	newContent := e.IndentationManager.IndentLines(content, 0, len(lines)-1)
	e.SetContent(newContent)
}

// UnindentSelectedLines removes indentation from the currently selected lines
func (e *Editor) UnindentSelectedLines() {
	content := e.GetContent()
	lines := strings.Split(content, "\n")
	
	// For simplified implementation, unindent all lines
	// In a real implementation, we would only unindent selected lines
	newContent := e.IndentationManager.UnindentLines(content, 0, len(lines)-1)
	e.SetContent(newContent)
}

// GetIndentationManager returns the indentation manager
func (e *Editor) GetIndentationManager() *IndentationManager {
	return e.IndentationManager
}

// SetTabSize sets the tab size for indentation
func (e *Editor) SetTabSize(size int) {
	if e.IndentationManager != nil {
		e.IndentationManager.SetTabSize(size)
	}
}

// SetUseSpaces sets whether to use spaces instead of tabs
func (e *Editor) SetUseSpaces(useSpaces bool) {
	if e.IndentationManager != nil {
		e.IndentationManager.SetUseSpaces(useSpaces)
	}
}

// SetAutoIndent sets whether to automatically indent new lines
func (e *Editor) SetAutoIndent(autoIndent bool) {
	if e.IndentationManager != nil {
		e.IndentationManager.SetAutoIndent(autoIndent)
	}
}

// HandleKeyEvent handles keyboard events for search functionality and editor features
func (e *Editor) HandleKeyEvent(event *fyne.KeyEvent) bool {
	// Let dialogs handle their events first
	if e.FindDialog != nil && e.FindDialog.HandleKeyEvent(event) {
		return true
	}
	if e.ReplaceDialog != nil && e.ReplaceDialog.HandleKeyEvent(event) {
		return true
	}
	if e.GoToLineDialog != nil && e.GoToLineDialog.HandleKeyEvent(event) {
		return true
	}
	
	// Handle editor-specific key events
	switch event.Name {
	case fyne.KeyF3:
		e.FindNext()
		return true
	case fyne.KeyTab:
		// Note: Fyne's KeyEvent doesn't have a Modifier field in the current version
		// This would need to be handled differently in a real implementation
		e.HandleTabKey()
		return true
	}
	
	return false
}

// SimulateCursorMovement simulates cursor movement for demonstration purposes
// In a real implementation, this would be handled by Fyne's text widget events
func (e *Editor) SimulateCursorMovement() {
	content := e.GetContent()
	lines := strings.Split(content, "\n")
	
	// Simulate cursor at end of content
	if len(lines) > 0 {
		lastLine := lines[len(lines)-1]
		e.State.SetCursorPosition(len(lines), len(lastLine)+1)
		
		if e.OnCursorChanged != nil {
			e.OnCursorChanged(len(lines), len(lastLine)+1)
		}
	}
}

// UpdateStatusBarFromContent updates the status bar based on current content
// This is a helper method for demonstration purposes
func (e *Editor) UpdateStatusBarFromContent() {
	if e.StatusBar != nil {
		e.StatusBar.Refresh()
		
		// Simulate cursor position update
		e.SimulateCursorMovement()
	}
}

// EditorView creates a syntax-highlighted view of a file (legacy function)
func EditorView(path string) fyne.CanvasObject {
	code, err := os.ReadFile(path)
	if err != nil {
		return widget.NewLabel("Failed to open file: " + err.Error())
	}

	highlighted := syntax.HighlightGoCode(string(code))
	richText := widget.NewRichText(highlighted...)
	scroll := container.NewScroll(richText)

	return scroll
}
