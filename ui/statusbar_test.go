package ui

import (
	"testing"
	
	"fyne.io/fyne/v2/test"
)

func TestStatusBarCreation(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	
	// Verify status bar was created
	if editor.StatusBar == nil {
		t.Fatal("StatusBar should not be nil")
	}
	
	// Verify status bar components are initialized
	if editor.StatusBar.positionLabel == nil {
		t.Error("Position label should not be nil")
	}
	if editor.StatusBar.fileInfoLabel == nil {
		t.Error("File info label should not be nil")
	}
	if editor.StatusBar.modifiedLabel == nil {
		t.Error("Modified label should not be nil")
	}
	if editor.StatusBar.encodingLabel == nil {
		t.Error("Encoding label should not be nil")
	}
	if editor.StatusBar.languageLabel == nil {
		t.Error("Language label should not be nil")
	}
	if editor.StatusBar.selectionLabel == nil {
		t.Error("Selection label should not be nil")
	}
}

func TestStatusBarPositionUpdate(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	statusBar := editor.StatusBar
	
	// Test position update
	statusBar.UpdatePosition(10, 25)
	
	expected := "Ln 10, Col 25"
	if statusBar.positionLabel.Text != expected {
		t.Errorf("Expected position label to be '%s', got '%s'", expected, statusBar.positionLabel.Text)
	}
}

func TestStatusBarFileInfoUpdate(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	statusBar := editor.StatusBar
	
	// Test file info update with path
	statusBar.UpdateFileInfo("/path/to/test.go", 1024)
	
	expected := "test.go (1.0 KB)"
	if statusBar.fileInfoLabel.Text != expected {
		t.Errorf("Expected file info label to be '%s', got '%s'", expected, statusBar.fileInfoLabel.Text)
	}
	
	// Test with empty path
	statusBar.UpdateFileInfo("", 0)
	
	expected = "Untitled"
	if statusBar.fileInfoLabel.Text != expected {
		t.Errorf("Expected file info label to be '%s', got '%s'", expected, statusBar.fileInfoLabel.Text)
	}
}

func TestStatusBarModifiedStatus(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	statusBar := editor.StatusBar
	
	// Test modified status
	statusBar.SetModified(true)
	if statusBar.modifiedLabel.Text != "●" {
		t.Errorf("Expected modified label to be '●', got '%s'", statusBar.modifiedLabel.Text)
	}
	
	statusBar.SetModified(false)
	if statusBar.modifiedLabel.Text != "" {
		t.Errorf("Expected modified label to be empty, got '%s'", statusBar.modifiedLabel.Text)
	}
}

func TestStatusBarLanguageUpdate(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	statusBar := editor.StatusBar
	
	// Test language update
	statusBar.SetLanguage("Go")
	if statusBar.languageLabel.Text != "Go" {
		t.Errorf("Expected language label to be 'Go', got '%s'", statusBar.languageLabel.Text)
	}
	
	// Test with empty language
	statusBar.SetLanguage("")
	if statusBar.languageLabel.Text != "Plain Text" {
		t.Errorf("Expected language label to be 'Plain Text', got '%s'", statusBar.languageLabel.Text)
	}
}

func TestStatusBarEncodingUpdate(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	statusBar := editor.StatusBar
	
	// Test encoding update
	statusBar.SetEncoding("UTF-16")
	if statusBar.encodingLabel.Text != "UTF-16" {
		t.Errorf("Expected encoding label to be 'UTF-16', got '%s'", statusBar.encodingLabel.Text)
	}
	
	// Test with empty encoding
	statusBar.SetEncoding("")
	if statusBar.encodingLabel.Text != "UTF-8" {
		t.Errorf("Expected encoding label to be 'UTF-8', got '%s'", statusBar.encodingLabel.Text)
	}
}

func TestStatusBarSelectionUpdate(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	statusBar := editor.StatusBar
	
	// Test selection update with single line
	statusBar.UpdateSelection(true, "hello world")
	expected := "(11 chars selected)"
	if statusBar.selectionLabel.Text != expected {
		t.Errorf("Expected selection label to be '%s', got '%s'", expected, statusBar.selectionLabel.Text)
	}
	
	// Test selection update with multiple lines
	statusBar.UpdateSelection(true, "line 1\nline 2\nline 3")
	expected = "(20 chars, 3 lines selected)"
	if statusBar.selectionLabel.Text != expected {
		t.Errorf("Expected selection label to be '%s', got '%s'", expected, statusBar.selectionLabel.Text)
	}
	
	// Test no selection
	statusBar.UpdateSelection(false, "")
	if statusBar.selectionLabel.Text != "" {
		t.Errorf("Expected selection label to be empty, got '%s'", statusBar.selectionLabel.Text)
	}
}

func TestStatusBarRefresh(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	statusBar := editor.StatusBar
	
	// Set up some state in the editor
	editor.State.SetCurrentFile("/test/file.go", 2048, "Go")
	editor.State.SetModified(true)
	editor.State.SetCursorPosition(15, 30)
	
	// Refresh the status bar
	statusBar.Refresh()
	
	// Verify all components were updated
	if statusBar.positionLabel.Text != "Ln 15, Col 30" {
		t.Errorf("Position not updated correctly: %s", statusBar.positionLabel.Text)
	}
	
	if statusBar.fileInfoLabel.Text != "file.go (2.0 KB)" {
		t.Errorf("File info not updated correctly: %s", statusBar.fileInfoLabel.Text)
	}
	
	if statusBar.modifiedLabel.Text != "●" {
		t.Errorf("Modified status not updated correctly: %s", statusBar.modifiedLabel.Text)
	}
	
	if statusBar.languageLabel.Text != "Go" {
		t.Errorf("Language not updated correctly: %s", statusBar.languageLabel.Text)
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 bytes"},
		{512, "512 bytes"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}
	
	for _, test := range tests {
		result := formatFileSize(test.size)
		if result != test.expected {
			t.Errorf("formatFileSize(%d) = %s, expected %s", test.size, result, test.expected)
		}
	}
}

func TestStatusBarIntegrationWithEditor(t *testing.T) {
	// Create a test app
	testApp := test.NewApp()
	defer testApp.Quit()
	
	// Create an editor instance
	editor := NewEditor()
	
	// Simulate loading a file
	editor.State.SetCurrentFile("/test/example.go", 1024, "Go")
	editor.State.SetModified(false)
	editor.State.SetCursorPosition(5, 10)
	
	// Trigger status bar update
	if editor.OnFileChanged != nil {
		editor.OnFileChanged("/test/example.go")
	}
	
	// Verify status bar reflects the state
	statusBar := editor.StatusBar
	if statusBar.positionLabel.Text != "Ln 5, Col 10" {
		t.Errorf("Position not synchronized: %s", statusBar.positionLabel.Text)
	}
	
	if statusBar.fileInfoLabel.Text != "example.go (1.0 KB)" {
		t.Errorf("File info not synchronized: %s", statusBar.fileInfoLabel.Text)
	}
	
	if statusBar.modifiedLabel.Text != "" {
		t.Errorf("Modified status not synchronized: %s", statusBar.modifiedLabel.Text)
	}
	
	// Simulate modification
	editor.State.SetModified(true)
	if editor.OnModified != nil {
		editor.OnModified(true)
	}
	
	if statusBar.modifiedLabel.Text != "●" {
		t.Errorf("Modified status not updated after change: %s", statusBar.modifiedLabel.Text)
	}
}