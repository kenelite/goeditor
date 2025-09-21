package dialogs

import (
	"testing"
	"fyne.io/fyne/v2/test"
	"github.com/kenelite/goeditor/backend"
)

// MockEditor implements EditorInterface for testing
type MockEditor struct {
	content string
}

func (m *MockEditor) GetContent() string {
	return m.content
}

func (m *MockEditor) SetContent(content string) {
	m.content = content
}

func TestFindDialog_Creation(t *testing.T) {
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()

	editor := &MockEditor{content: "Hello World\nThis is a test"}
	searchManager := backend.NewSearchManager()
	
	dialog := NewFindDialog(editor, searchManager, window)
	
	if dialog == nil {
		t.Error("FindDialog should not be nil")
	}
	
	if dialog.searchEntry == nil {
		t.Error("Search entry should not be nil")
	}
	
	if dialog.resultLabel == nil {
		t.Error("Result label should not be nil")
	}
}

func TestFindDialog_Search(t *testing.T) {
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()

	editor := &MockEditor{content: "Hello World\nHello Universe"}
	searchManager := backend.NewSearchManager()
	
	dialog := NewFindDialog(editor, searchManager, window)
	
	// Test search functionality
	dialog.SetSearchText("Hello")
	
	matches := searchManager.GetMatches()
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}
	
	// Test navigation
	found := dialog.FindNext()
	if !found {
		t.Error("FindNext should return true when matches exist")
	}
	
	currentIndex := searchManager.GetCurrentIndex()
	if currentIndex != 1 {
		t.Errorf("Expected current index 1 after FindNext, got %d", currentIndex)
	}
}

func TestReplaceDialog_Creation(t *testing.T) {
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()

	editor := &MockEditor{content: "Hello World\nThis is a test"}
	searchManager := backend.NewSearchManager()
	
	dialog := NewReplaceDialog(editor, searchManager, window)
	
	if dialog == nil {
		t.Error("ReplaceDialog should not be nil")
	}
	
	if dialog.replaceEntry == nil {
		t.Error("Replace entry should not be nil")
	}
	
	if dialog.replaceButton == nil {
		t.Error("Replace button should not be nil")
	}
	
	if dialog.replaceAllButton == nil {
		t.Error("Replace all button should not be nil")
	}
}

func TestReplaceDialog_Replace(t *testing.T) {
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()

	editor := &MockEditor{content: "Hello World\nHello Universe"}
	searchManager := backend.NewSearchManager()
	
	dialog := NewReplaceDialog(editor, searchManager, window)
	
	// Set search and replace text
	dialog.SetSearchText("Hello")
	dialog.SetReplaceText("Hi")
	
	// Test replace all
	dialog.ReplaceAll()
	
	expected := "Hi World\nHi Universe"
	if editor.GetContent() != expected {
		t.Errorf("Expected content '%s', got '%s'", expected, editor.GetContent())
	}
}

func TestReplaceDialog_ReplaceCurrent(t *testing.T) {
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()

	editor := &MockEditor{content: "Hello World\nHello Universe"}
	searchManager := backend.NewSearchManager()
	
	dialog := NewReplaceDialog(editor, searchManager, window)
	
	// Set search and replace text
	dialog.SetSearchText("Hello")
	dialog.SetReplaceText("Hi")
	
	// Find first match
	dialog.FindNext()
	
	// Replace current match
	dialog.ReplaceCurrent()
	
	expected := "Hi World\nHello Universe"
	if editor.GetContent() != expected {
		t.Errorf("Expected content '%s', got '%s'", expected, editor.GetContent())
	}
}

func TestFindDialog_Options(t *testing.T) {
	app := test.NewApp()
	window := test.NewWindow(nil)
	defer app.Quit()

	editor := &MockEditor{content: "Hello HELLO hello"}
	searchManager := backend.NewSearchManager()
	
	dialog := NewFindDialog(editor, searchManager, window)
	
	// Test case insensitive (default)
	dialog.SetSearchText("hello")
	matches := searchManager.GetMatches()
	if len(matches) != 3 {
		t.Errorf("Case insensitive: Expected 3 matches, got %d", len(matches))
	}
	
	// Test case sensitive
	dialog.optionsCheck["caseSensitive"].SetChecked(true)
	dialog.updateSearchOptions()
	
	matches = searchManager.GetMatches()
	if len(matches) != 1 {
		t.Errorf("Case sensitive: Expected 1 match, got %d", len(matches))
	}
}