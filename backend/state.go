package backend

import (
	"time"
)

// Position represents a cursor or selection position in the editor
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// EditorState holds the complete state of the editor
type EditorState struct {
	CurrentFile    string    `json:"currentFile"`
	IsModified     bool      `json:"isModified"`
	CursorLine     int       `json:"cursorLine"`
	CursorColumn   int       `json:"cursorColumn"`
	SelectionStart Position  `json:"selectionStart"`
	SelectionEnd   Position  `json:"selectionEnd"`
	ScrollPosition Position  `json:"scrollPosition"`
	Language       string    `json:"language"`
	Encoding       string    `json:"encoding"`
	FileSize       int64     `json:"fileSize"`
	LastModified   time.Time `json:"lastModified"`
}

// NewEditorState creates a new editor state with default values
func NewEditorState() *EditorState {
	return &EditorState{
		CurrentFile:    "",
		IsModified:     false,
		CursorLine:     1,
		CursorColumn:   1,
		SelectionStart: Position{Line: 1, Column: 1},
		SelectionEnd:   Position{Line: 1, Column: 1},
		ScrollPosition: Position{Line: 1, Column: 1},
		Language:       "text",
		Encoding:       "UTF-8",
		FileSize:       0,
		LastModified:   time.Now(),
	}
}

// SetCurrentFile updates the current file and related metadata
func (s *EditorState) SetCurrentFile(path string, size int64, language string) {
	s.CurrentFile = path
	s.FileSize = size
	s.Language = language
	s.IsModified = false
	s.LastModified = time.Now()
}

// SetModified marks the file as modified or unmodified
func (s *EditorState) SetModified(modified bool) {
	s.IsModified = modified
	if !modified {
		s.LastModified = time.Now()
	}
}

// SetCursorPosition updates the cursor position
func (s *EditorState) SetCursorPosition(line, column int) {
	s.CursorLine = line
	s.CursorColumn = column
}

// GetCursorPosition returns the current cursor position
func (s *EditorState) GetCursorPosition() (int, int) {
	return s.CursorLine, s.CursorColumn
}

// SetSelection updates the selection range
func (s *EditorState) SetSelection(start, end Position) {
	s.SelectionStart = start
	s.SelectionEnd = end
}

// GetSelection returns the current selection range
func (s *EditorState) GetSelection() (Position, Position) {
	return s.SelectionStart, s.SelectionEnd
}

// HasSelection returns true if there is an active selection
func (s *EditorState) HasSelection() bool {
	return s.SelectionStart != s.SelectionEnd
}

// ClearSelection clears the current selection
func (s *EditorState) ClearSelection() {
	pos := Position{Line: s.CursorLine, Column: s.CursorColumn}
	s.SelectionStart = pos
	s.SelectionEnd = pos
}

// Legacy State struct for backward compatibility
type State struct {
	CurrentFile string
}

// NewState creates a new legacy state (for backward compatibility)
func NewState() *State {
	return &State{}
}
