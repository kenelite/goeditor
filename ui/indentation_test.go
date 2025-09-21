package ui

import (
	"testing"
)

func TestIndentationManager(t *testing.T) {
	im := NewIndentationManager()
	
	// Test default settings
	if im.GetTabSize() != 4 {
		t.Errorf("Expected default tab size to be 4, got %d", im.GetTabSize())
	}
	
	if !im.GetUseSpaces() {
		t.Error("Expected default to use spaces")
	}
	
	if !im.GetAutoIndent() {
		t.Error("Expected default auto-indent to be enabled")
	}
	
	// Test indent string
	indentStr := im.GetIndentString()
	if indentStr != "    " { // 4 spaces
		t.Errorf("Expected 4 spaces, got %q", indentStr)
	}
	
	// Test with tabs
	im.SetUseSpaces(false)
	indentStr = im.GetIndentString()
	if indentStr != "\t" {
		t.Errorf("Expected tab character, got %q", indentStr)
	}
	
	// Test tab size change
	im.SetTabSize(2)
	im.SetUseSpaces(true)
	indentStr = im.GetIndentString()
	if indentStr != "  " { // 2 spaces
		t.Errorf("Expected 2 spaces, got %q", indentStr)
	}
}

func TestIndentationOperations(t *testing.T) {
	im := NewIndentationManager()
	
	// Test line indentation detection
	line := "    hello world"
	indentation := im.GetLineIndentation(line)
	if indentation != "    " {
		t.Errorf("Expected 4 spaces indentation, got %q", indentation)
	}
	
	// Test indentation level
	level := im.GetIndentationLevel(line)
	if level != 1 {
		t.Errorf("Expected indentation level 1, got %d", level)
	}
	
	// Test removing indentation
	newLine, removed := im.removeIndentation(line)
	if newLine != "hello world" || removed != 4 {
		t.Errorf("Expected 'hello world' and 4 removed chars, got %q and %d", newLine, removed)
	}
	
	// Test whitespace-only line
	if !im.IsWhitespaceOnly("   \t  ") {
		t.Error("Expected whitespace-only line to be detected")
	}
	
	if im.IsWhitespaceOnly("  hello  ") {
		t.Error("Expected non-whitespace line to not be detected as whitespace-only")
	}
}

func TestIndentLines(t *testing.T) {
	im := NewIndentationManager()
	
	content := "line 1\nline 2\nline 3"
	expected := "    line 1\n    line 2\n    line 3"
	
	result := im.IndentLines(content, 0, 2)
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestUnindentLines(t *testing.T) {
	im := NewIndentationManager()
	
	content := "    line 1\n    line 2\n    line 3"
	expected := "line 1\nline 2\nline 3"
	
	result := im.UnindentLines(content, 0, 2)
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestTabConversion(t *testing.T) {
	im := NewIndentationManager()
	
	// Test tabs to spaces
	content := "\tline 1\n\t\tline 2"
	expected := "    line 1\n        line 2"
	
	result := im.ConvertTabsToSpaces(content)
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
	
	// Test spaces to tabs
	content = "    line 1\n        line 2"
	expected = "\tline 1\n\t\tline 2"
	
	result = im.ConvertSpacesToTabs(content)
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}