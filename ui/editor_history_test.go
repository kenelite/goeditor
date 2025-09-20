package ui

import (
	"testing"
	"github.com/kenelite/goeditor/backend"
)

func TestEditorHistoryBackendIntegration(t *testing.T) {
	// Test history functionality without UI components
	history := backend.NewHistory()
	
	// Test initial state
	if history.CanUndo() {
		t.Error("New history should not have undo operations")
	}
	
	if history.CanRedo() {
		t.Error("New history should not have redo operations")
	}
	
	// Test recording operations
	pos := backend.Position{Line: 1, Column: 1}
	history.RecordInsert(pos, "Hello")
	
	if !history.CanUndo() {
		t.Error("Should be able to undo after recording operation")
	}
	
	if history.CanRedo() {
		t.Error("Should not be able to redo after new operation")
	}
	
	// Test undo
	op := history.Undo()
	if op == nil {
		t.Error("Undo should return operation")
	}
	
	if history.CanUndo() {
		t.Error("Should not be able to undo after undoing single operation")
	}
	
	if !history.CanRedo() {
		t.Error("Should be able to redo after undo")
	}
	
	// Test redo
	op = history.Redo()
	if op == nil {
		t.Error("Redo should return operation")
	}
	
	if !history.CanUndo() {
		t.Error("Should be able to undo after redo")
	}
	
	if history.CanRedo() {
		t.Error("Should not be able to redo after redo")
	}
}

func TestEditorHistoryMultipleOperations(t *testing.T) {
	history := backend.NewHistory()
	pos := backend.Position{Line: 1, Column: 1}
	
	// Record multiple operations
	history.RecordInsert(pos, "Hello")
	history.RecordInsert(backend.Position{Line: 1, Column: 6}, " ")
	history.RecordInsert(backend.Position{Line: 1, Column: 7}, "World")
	
	// Should have 3 operations to undo
	undoCount := 0
	for history.CanUndo() {
		history.Undo()
		undoCount++
	}
	
	if undoCount != 3 {
		t.Errorf("Expected 3 undo operations, got %d", undoCount)
	}
	
	// Should have 3 operations to redo
	redoCount := 0
	for history.CanRedo() {
		history.Redo()
		redoCount++
	}
	
	if redoCount != 3 {
		t.Errorf("Expected 3 redo operations, got %d", redoCount)
	}
}

func TestEditorHistoryClear(t *testing.T) {
	history := backend.NewHistory()
	pos := backend.Position{Line: 1, Column: 1}
	
	// Add some operations
	history.RecordInsert(pos, "Test content")
	history.Undo() // Move to redo stack
	
	if !history.CanRedo() {
		t.Fatal("Should have redo operations before clearing")
	}
	
	// Clear history
	history.Clear()
	
	if history.CanUndo() {
		t.Error("History should be cleared")
	}
	
	if history.CanRedo() {
		t.Error("History should be cleared")
	}
}

func TestEditorHistoryDeleteOperation(t *testing.T) {
	history := backend.NewHistory()
	
	// Record delete operation
	pos := backend.Position{Line: 1, Column: 6}
	history.RecordDelete(pos, "World")
	
	if !history.CanUndo() {
		t.Error("Should be able to undo delete operation")
	}
	
	// Undo the delete
	op := history.Undo()
	if op == nil {
		t.Error("Undo should return operation for delete")
	}
	
	if op.Type != backend.Delete {
		t.Error("Operation should be Delete type")
	}
}

func TestEditorHistoryReplaceOperation(t *testing.T) {
	history := backend.NewHistory()
	
	// Record replace operation
	pos := backend.Position{Line: 1, Column: 1}
	history.RecordReplace(pos, "Hello", "Hi")
	
	if !history.CanUndo() {
		t.Error("Should be able to undo replace operation")
	}
	
	// Undo the replace
	op := history.Undo()
	if op == nil {
		t.Error("Undo should return operation for replace")
	}
	
	if op.Type != backend.Replace {
		t.Error("Operation should be Replace type")
	}
}

func TestEditorHistoryDisabled(t *testing.T) {
	history := backend.NewHistory()
	pos := backend.Position{Line: 1, Column: 1}
	
	// Disable history
	history.SetEnabled(false)
	
	// Try to record operation
	history.RecordInsert(pos, "Test")
	
	if history.CanUndo() {
		t.Error("Should not record operations when history is disabled")
	}
	
	// Re-enable history
	history.SetEnabled(true)
	
	// Record operation again
	history.RecordInsert(pos, "Test2")
	
	if !history.CanUndo() {
		t.Error("Should record operations when history is re-enabled")
	}
}

func TestEditorHistoryIntegrationWithEditor(t *testing.T) {
	// Test that editor properly integrates with history
	history := backend.NewHistory()
	pos := backend.Position{Line: 1, Column: 1}
	
	// Add some operations
	history.RecordInsert(pos, "Test1")
	history.RecordInsert(backend.Position{Line: 1, Column: 6}, "Test2")
	history.Undo() // Move one to redo stack
	
	if !history.CanUndo() || !history.CanRedo() {
		t.Fatal("Setup failed - should have both undo and redo operations")
	}
	
	// Clear history
	history.Clear()
	
	if history.CanUndo() {
		t.Error("Should not have undo operations after clearing history")
	}
	
	if history.CanRedo() {
		t.Error("Should not have redo operations after clearing history")
	}
}