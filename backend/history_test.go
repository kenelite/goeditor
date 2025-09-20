package backend

import (
	"testing"
	"time"
)

func TestNewHistory(t *testing.T) {
	h := NewHistory()
	
	if h == nil {
		t.Fatal("NewHistory() returned nil")
	}
	
	if h.GetMaxSize() != 100 {
		t.Errorf("Expected default max size 100, got %d", h.GetMaxSize())
	}
	
	if !h.IsEnabled() {
		t.Error("Expected history to be enabled by default")
	}
	
	if h.CanUndo() {
		t.Error("New history should not have undo operations")
	}
	
	if h.CanRedo() {
		t.Error("New history should not have redo operations")
	}
}

func TestNewHistoryWithSize(t *testing.T) {
	h := NewHistoryWithSize(50)
	
	if h.GetMaxSize() != 50 {
		t.Errorf("Expected max size 50, got %d", h.GetMaxSize())
	}
	
	// Test with invalid size
	h2 := NewHistoryWithSize(-1)
	if h2.GetMaxSize() != 100 {
		t.Errorf("Expected default max size 100 for invalid input, got %d", h2.GetMaxSize())
	}
}

func TestOperationTypeString(t *testing.T) {
	tests := []struct {
		op       OperationType
		expected string
	}{
		{Insert, "Insert"},
		{Delete, "Delete"},
		{Replace, "Replace"},
		{OperationType(999), "Unknown"},
	}
	
	for _, test := range tests {
		if test.op.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.op.String())
		}
	}
}

func TestRecordOperation(t *testing.T) {
	h := NewHistory()
	
	op := Operation{
		Type:     Insert,
		Position: Position{Line: 1, Column: 1},
		NewText:  "Hello",
	}
	
	h.RecordOperation(op)
	
	if !h.CanUndo() {
		t.Error("Should be able to undo after recording operation")
	}
	
	if h.CanRedo() {
		t.Error("Should not be able to redo after recording new operation")
	}
	
	if h.GetUndoCount() != 1 {
		t.Errorf("Expected undo count 1, got %d", h.GetUndoCount())
	}
}

func TestRecordInsert(t *testing.T) {
	h := NewHistory()
	pos := Position{Line: 1, Column: 5}
	text := "world"
	
	h.RecordInsert(pos, text)
	
	if !h.CanUndo() {
		t.Error("Should be able to undo after insert")
	}
	
	op := h.GetLastOperation()
	if op == nil {
		t.Fatal("Expected operation to be recorded")
	}
	
	if op.Type != Insert {
		t.Errorf("Expected Insert operation, got %s", op.Type.String())
	}
	
	if op.Position != pos {
		t.Errorf("Expected position %+v, got %+v", pos, op.Position)
	}
	
	if op.NewText != text {
		t.Errorf("Expected new text '%s', got '%s'", text, op.NewText)
	}
	
	if op.OldText != "" {
		t.Errorf("Expected empty old text, got '%s'", op.OldText)
	}
}

func TestRecordDelete(t *testing.T) {
	h := NewHistory()
	pos := Position{Line: 2, Column: 3}
	text := "deleted"
	
	h.RecordDelete(pos, text)
	
	op := h.GetLastOperation()
	if op == nil {
		t.Fatal("Expected operation to be recorded")
	}
	
	if op.Type != Delete {
		t.Errorf("Expected Delete operation, got %s", op.Type.String())
	}
	
	if op.OldText != text {
		t.Errorf("Expected old text '%s', got '%s'", text, op.OldText)
	}
	
	if op.NewText != "" {
		t.Errorf("Expected empty new text, got '%s'", op.NewText)
	}
}

func TestRecordReplace(t *testing.T) {
	h := NewHistory()
	pos := Position{Line: 3, Column: 7}
	oldText := "old"
	newText := "new"
	
	h.RecordReplace(pos, oldText, newText)
	
	op := h.GetLastOperation()
	if op == nil {
		t.Fatal("Expected operation to be recorded")
	}
	
	if op.Type != Replace {
		t.Errorf("Expected Replace operation, got %s", op.Type.String())
	}
	
	if op.OldText != oldText {
		t.Errorf("Expected old text '%s', got '%s'", oldText, op.OldText)
	}
	
	if op.NewText != newText {
		t.Errorf("Expected new text '%s', got '%s'", newText, op.NewText)
	}
}

func TestUndoRedo(t *testing.T) {
	h := NewHistory()
	
	// Record some operations
	h.RecordInsert(Position{Line: 1, Column: 1}, "Hello")
	h.RecordInsert(Position{Line: 1, Column: 6}, " World")
	
	if h.GetUndoCount() != 2 {
		t.Errorf("Expected 2 undo operations, got %d", h.GetUndoCount())
	}
	
	// Test undo
	op1 := h.Undo()
	if op1 == nil {
		t.Fatal("Expected undo operation")
	}
	
	if op1.NewText != " World" {
		t.Errorf("Expected ' World', got '%s'", op1.NewText)
	}
	
	if h.GetUndoCount() != 1 {
		t.Errorf("Expected 1 undo operation after undo, got %d", h.GetUndoCount())
	}
	
	if h.GetRedoCount() != 1 {
		t.Errorf("Expected 1 redo operation after undo, got %d", h.GetRedoCount())
	}
	
	// Test redo
	op2 := h.Redo()
	if op2 == nil {
		t.Fatal("Expected redo operation")
	}
	
	if op2.NewText != " World" {
		t.Errorf("Expected ' World', got '%s'", op2.NewText)
	}
	
	if h.GetUndoCount() != 2 {
		t.Errorf("Expected 2 undo operations after redo, got %d", h.GetUndoCount())
	}
	
	if h.GetRedoCount() != 0 {
		t.Errorf("Expected 0 redo operations after redo, got %d", h.GetRedoCount())
	}
}

func TestUndoRedoEmpty(t *testing.T) {
	h := NewHistory()
	
	// Test undo on empty history
	op := h.Undo()
	if op != nil {
		t.Error("Expected nil when undoing empty history")
	}
	
	// Test redo on empty history
	op = h.Redo()
	if op != nil {
		t.Error("Expected nil when redoing empty history")
	}
}

func TestClearHistory(t *testing.T) {
	h := NewHistory()
	
	// Add some operations
	h.RecordInsert(Position{Line: 1, Column: 1}, "test1")
	h.RecordInsert(Position{Line: 1, Column: 6}, "test2")
	h.Undo() // Move one to redo stack
	
	if h.GetUndoCount() == 0 || h.GetRedoCount() == 0 {
		t.Fatal("Setup failed - should have operations in both stacks")
	}
	
	// Clear history
	h.Clear()
	
	if h.GetUndoCount() != 0 {
		t.Errorf("Expected 0 undo operations after clear, got %d", h.GetUndoCount())
	}
	
	if h.GetRedoCount() != 0 {
		t.Errorf("Expected 0 redo operations after clear, got %d", h.GetRedoCount())
	}
	
	if h.CanUndo() {
		t.Error("Should not be able to undo after clear")
	}
	
	if h.CanRedo() {
		t.Error("Should not be able to redo after clear")
	}
}

func TestMaxSizeLimit(t *testing.T) {
	h := NewHistoryWithSize(3)
	
	// Add more operations than max size
	for i := 0; i < 5; i++ {
		h.RecordInsert(Position{Line: 1, Column: i}, string(rune('a'+i)))
	}
	
	if h.GetUndoCount() != 3 {
		t.Errorf("Expected max 3 operations, got %d", h.GetUndoCount())
	}
	
	// Check that oldest operations were removed
	history := h.GetOperationHistory()
	if len(history) != 3 {
		t.Errorf("Expected 3 operations in history, got %d", len(history))
	}
	
	// Should have operations 'c', 'd', 'e' (last 3)
	if history[0].NewText != "c" {
		t.Errorf("Expected first operation to be 'c', got '%s'", history[0].NewText)
	}
}

func TestSetMaxSize(t *testing.T) {
	h := NewHistory()
	
	// Add some operations
	for i := 0; i < 10; i++ {
		h.RecordInsert(Position{Line: 1, Column: i}, string(rune('a'+i)))
	}
	
	// Reduce max size
	h.SetMaxSize(5)
	
	if h.GetMaxSize() != 5 {
		t.Errorf("Expected max size 5, got %d", h.GetMaxSize())
	}
	
	if h.GetUndoCount() != 5 {
		t.Errorf("Expected 5 operations after reducing max size, got %d", h.GetUndoCount())
	}
	
	// Test invalid max size
	h.SetMaxSize(-1)
	if h.GetMaxSize() != 5 {
		t.Error("Max size should not change for invalid input")
	}
}

func TestEnabledDisabled(t *testing.T) {
	h := NewHistory()
	
	// Disable history
	h.SetEnabled(false)
	
	if h.IsEnabled() {
		t.Error("History should be disabled")
	}
	
	// Try to record operation while disabled
	h.RecordInsert(Position{Line: 1, Column: 1}, "test")
	
	if h.CanUndo() {
		t.Error("Should not record operations when disabled")
	}
	
	// Re-enable and test
	h.SetEnabled(true)
	h.RecordInsert(Position{Line: 1, Column: 1}, "test")
	
	if !h.CanUndo() {
		t.Error("Should record operations when enabled")
	}
}

func TestTimestamp(t *testing.T) {
	h := NewHistory()
	
	before := time.Now()
	h.RecordInsert(Position{Line: 1, Column: 1}, "test")
	after := time.Now()
	
	op := h.GetLastOperation()
	if op == nil {
		t.Fatal("Expected operation to be recorded")
	}
	
	if op.Timestamp.Before(before) || op.Timestamp.After(after) {
		t.Error("Operation timestamp should be between before and after times")
	}
}

func TestRedoStackClearedOnNewOperation(t *testing.T) {
	h := NewHistory()
	
	// Record operations and undo some
	h.RecordInsert(Position{Line: 1, Column: 1}, "first")
	h.RecordInsert(Position{Line: 1, Column: 6}, "second")
	h.Undo() // Move "second" to redo stack
	
	if h.GetRedoCount() != 1 {
		t.Fatal("Setup failed - should have 1 redo operation")
	}
	
	// Record new operation
	h.RecordInsert(Position{Line: 1, Column: 6}, "third")
	
	if h.GetRedoCount() != 0 {
		t.Error("Redo stack should be cleared when new operation is recorded")
	}
}

func TestGetOperationHistory(t *testing.T) {
	h := NewHistory()
	
	// Record some operations
	h.RecordInsert(Position{Line: 1, Column: 1}, "first")
	h.RecordDelete(Position{Line: 1, Column: 6}, "second")
	h.RecordReplace(Position{Line: 1, Column: 1}, "old", "new")
	
	history := h.GetOperationHistory()
	
	if len(history) != 3 {
		t.Errorf("Expected 3 operations in history, got %d", len(history))
	}
	
	// Verify operations are in correct order
	if history[0].Type != Insert || history[0].NewText != "first" {
		t.Error("First operation should be Insert 'first'")
	}
	
	if history[1].Type != Delete || history[1].OldText != "second" {
		t.Error("Second operation should be Delete 'second'")
	}
	
	if history[2].Type != Replace || history[2].OldText != "old" || history[2].NewText != "new" {
		t.Error("Third operation should be Replace 'old' -> 'new'")
	}
	
	// Verify it's a copy (modifying returned slice shouldn't affect original)
	history[0].NewText = "modified"
	originalHistory := h.GetOperationHistory()
	if originalHistory[0].NewText == "modified" {
		t.Error("GetOperationHistory should return a copy, not the original slice")
	}
}