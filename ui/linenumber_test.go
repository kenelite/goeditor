package ui

import (
	"testing"
)

func TestLineNumberWidget(t *testing.T) {
	// Create a simple line number widget without full editor
	ln := &LineNumberWidget{
		lineCount:  1,
		lineHeight: 20.0,
		fontHeight: 14.0,
		padding:    5.0,
	}
	
	// Test initial state
	if ln.GetLineCount() != 1 {
		t.Errorf("Expected initial line count to be 1, got %d", ln.GetLineCount())
	}
	
	// Test updating line count
	ln.UpdateLineCount(3)
	if ln.GetLineCount() != 3 {
		t.Errorf("Expected line count to be 3, got %d", ln.GetLineCount())
	}
	
	// Test preferred width calculation
	width := ln.GetPreferredWidth()
	if width <= 0 {
		t.Errorf("Expected positive preferred width, got %f", width)
	}
	
	// Test line height setting
	ln.SetLineHeight(25.0)
	if ln.lineHeight != 25.0 {
		t.Errorf("Expected line height to be 25.0, got %f", ln.lineHeight)
	}
	
	// Test font height setting
	ln.SetFontHeight(16.0)
	if ln.fontHeight != 16.0 {
		t.Errorf("Expected font height to be 16.0, got %f", ln.fontHeight)
	}
}

