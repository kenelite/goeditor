package backend

import (
	"time"
)

// OperationType represents the type of operation performed
type OperationType int

const (
	// Insert represents text insertion operation
	Insert OperationType = iota
	// Delete represents text deletion operation
	Delete
	// Replace represents text replacement operation
	Replace
)

// String returns the string representation of OperationType
func (ot OperationType) String() string {
	switch ot {
	case Insert:
		return "Insert"
	case Delete:
		return "Delete"
	case Replace:
		return "Replace"
	default:
		return "Unknown"
	}
}

// Operation represents a single edit operation that can be undone/redone
type Operation struct {
	Type      OperationType `json:"type"`
	Position  Position      `json:"position"`
	OldText   string        `json:"oldText"`
	NewText   string        `json:"newText"`
	Timestamp time.Time     `json:"timestamp"`
}

// History manages the undo/redo history for the editor
type History struct {
	undoStack []Operation
	redoStack []Operation
	maxSize   int
	enabled   bool
}

// NewHistory creates a new history manager with default settings
func NewHistory() *History {
	return &History{
		undoStack: make([]Operation, 0),
		redoStack: make([]Operation, 0),
		maxSize:   100, // Default maximum history size
		enabled:   true,
	}
}

// NewHistoryWithSize creates a new history manager with specified maximum size
func NewHistoryWithSize(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{
		undoStack: make([]Operation, 0),
		redoStack: make([]Operation, 0),
		maxSize:   maxSize,
		enabled:   true,
	}
}

// RecordOperation records a new operation in the history
func (h *History) RecordOperation(op Operation) {
	if !h.enabled {
		return
	}

	// Set timestamp if not already set
	if op.Timestamp.IsZero() {
		op.Timestamp = time.Now()
	}

	// Add operation to undo stack
	h.undoStack = append(h.undoStack, op)

	// Clear redo stack when new operation is recorded
	h.redoStack = h.redoStack[:0]

	// Maintain maximum size
	if len(h.undoStack) > h.maxSize {
		// Remove oldest operations
		copy(h.undoStack, h.undoStack[1:])
		h.undoStack = h.undoStack[:len(h.undoStack)-1]
	}
}

// RecordInsert records an insert operation
func (h *History) RecordInsert(position Position, text string) {
	op := Operation{
		Type:      Insert,
		Position:  position,
		OldText:   "",
		NewText:   text,
		Timestamp: time.Now(),
	}
	h.RecordOperation(op)
}

// RecordDelete records a delete operation
func (h *History) RecordDelete(position Position, text string) {
	op := Operation{
		Type:      Delete,
		Position:  position,
		OldText:   text,
		NewText:   "",
		Timestamp: time.Now(),
	}
	h.RecordOperation(op)
}

// RecordReplace records a replace operation
func (h *History) RecordReplace(position Position, oldText, newText string) {
	op := Operation{
		Type:      Replace,
		Position:  position,
		OldText:   oldText,
		NewText:   newText,
		Timestamp: time.Now(),
	}
	h.RecordOperation(op)
}

// CanUndo returns true if there are operations that can be undone
func (h *History) CanUndo() bool {
	return len(h.undoStack) > 0
}

// CanRedo returns true if there are operations that can be redone
func (h *History) CanRedo() bool {
	return len(h.redoStack) > 0
}

// Undo undoes the last operation and returns it, or nil if nothing to undo
func (h *History) Undo() *Operation {
	if !h.CanUndo() {
		return nil
	}

	// Pop operation from undo stack
	lastIndex := len(h.undoStack) - 1
	op := h.undoStack[lastIndex]
	h.undoStack = h.undoStack[:lastIndex]

	// Push to redo stack
	h.redoStack = append(h.redoStack, op)

	return &op
}

// Redo redoes the last undone operation and returns it, or nil if nothing to redo
func (h *History) Redo() *Operation {
	if !h.CanRedo() {
		return nil
	}

	// Pop operation from redo stack
	lastIndex := len(h.redoStack) - 1
	op := h.redoStack[lastIndex]
	h.redoStack = h.redoStack[:lastIndex]

	// Push to undo stack
	h.undoStack = append(h.undoStack, op)

	return &op
}

// Clear clears all history
func (h *History) Clear() {
	h.undoStack = h.undoStack[:0]
	h.redoStack = h.redoStack[:0]
}

// GetUndoCount returns the number of operations that can be undone
func (h *History) GetUndoCount() int {
	return len(h.undoStack)
}

// GetRedoCount returns the number of operations that can be redone
func (h *History) GetRedoCount() int {
	return len(h.redoStack)
}

// SetEnabled enables or disables history recording
func (h *History) SetEnabled(enabled bool) {
	h.enabled = enabled
}

// IsEnabled returns whether history recording is enabled
func (h *History) IsEnabled() bool {
	return h.enabled
}

// SetMaxSize sets the maximum number of operations to keep in history
func (h *History) SetMaxSize(maxSize int) {
	if maxSize <= 0 {
		return
	}

	h.maxSize = maxSize

	// Trim existing history if necessary
	if len(h.undoStack) > maxSize {
		excess := len(h.undoStack) - maxSize
		copy(h.undoStack, h.undoStack[excess:])
		h.undoStack = h.undoStack[:maxSize]
	}
}

// GetMaxSize returns the maximum history size
func (h *History) GetMaxSize() int {
	return h.maxSize
}

// GetLastOperation returns the last operation without removing it from history
func (h *History) GetLastOperation() *Operation {
	if len(h.undoStack) == 0 {
		return nil
	}
	op := h.undoStack[len(h.undoStack)-1]
	return &op
}

// GetOperationHistory returns a copy of the undo stack for inspection
func (h *History) GetOperationHistory() []Operation {
	history := make([]Operation, len(h.undoStack))
	copy(history, h.undoStack)
	return history
}