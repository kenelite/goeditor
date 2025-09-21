package backend

import (
	"testing"
)

func TestSearchManager_Find(t *testing.T) {
	sm := NewSearchManager()
	text := "Hello World\nThis is a test\nHello again"

	// Test basic search
	matches := sm.Find(text, "Hello")
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}

	// Verify first match
	if matches[0].Start.Line != 1 || matches[0].Start.Column != 1 {
		t.Errorf("First match position incorrect: Line %d, Column %d", matches[0].Start.Line, matches[0].Start.Column)
	}

	// Verify second match
	if matches[1].Start.Line != 3 || matches[1].Start.Column != 1 {
		t.Errorf("Second match position incorrect: Line %d, Column %d", matches[1].Start.Line, matches[1].Start.Column)
	}
}

func TestSearchManager_CaseSensitive(t *testing.T) {
	sm := NewSearchManager()
	text := "Hello HELLO hello"

	// Case insensitive (default)
	matches := sm.Find(text, "hello")
	if len(matches) != 3 {
		t.Errorf("Case insensitive: Expected 3 matches, got %d", len(matches))
	}

	// Case sensitive
	options := sm.GetOptions()
	options.CaseSensitive = true
	sm.SetOptions(options)
	
	matches = sm.Find(text, "hello")
	if len(matches) != 1 {
		t.Errorf("Case sensitive: Expected 1 match, got %d", len(matches))
	}
}

func TestSearchManager_WholeWord(t *testing.T) {
	sm := NewSearchManager()
	text := "test testing tested"

	// Without whole word
	matches := sm.Find(text, "test")
	if len(matches) != 3 {
		t.Errorf("Without whole word: Expected 3 matches, got %d", len(matches))
	}

	// With whole word
	options := sm.GetOptions()
	options.WholeWord = true
	sm.SetOptions(options)
	
	matches = sm.Find(text, "test")
	if len(matches) != 1 {
		t.Errorf("With whole word: Expected 1 match, got %d", len(matches))
	}
}

func TestSearchManager_RegularExpression(t *testing.T) {
	sm := NewSearchManager()
	text := "test123 test456 testABC"

	// Enable regex
	options := sm.GetOptions()
	options.RegularExpression = true
	sm.SetOptions(options)

	// Test regex pattern
	matches := sm.Find(text, `test\d+`)
	if len(matches) != 2 {
		t.Errorf("Regex search: Expected 2 matches, got %d", len(matches))
	}
}

func TestSearchManager_Navigation(t *testing.T) {
	sm := NewSearchManager()
	text := "test test test"

	matches := sm.Find(text, "test")
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matches))
	}

	// Test next navigation
	match := sm.NextMatch()
	if match == nil || sm.GetCurrentIndex() != 0 {
		t.Errorf("NextMatch failed, index: %d", sm.GetCurrentIndex())
	}

	match = sm.NextMatch()
	if match == nil || sm.GetCurrentIndex() != 1 {
		t.Errorf("Second NextMatch failed, index: %d", sm.GetCurrentIndex())
	}

	// Test previous navigation
	match = sm.PreviousMatch()
	if match == nil || sm.GetCurrentIndex() != 0 {
		t.Errorf("PreviousMatch failed, index: %d", sm.GetCurrentIndex())
	}
}

func TestSearchManager_WrapAround(t *testing.T) {
	sm := NewSearchManager()
	text := "test test test"

	matches := sm.Find(text, "test")
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matches))
	}

	// Go to last match
	sm.SetCurrentMatch(2)
	
	// Test wrap around forward
	match := sm.NextMatch()
	if match == nil || sm.GetCurrentIndex() != 0 {
		t.Errorf("Wrap around forward failed, index: %d", sm.GetCurrentIndex())
	}

	// Test wrap around backward
	match = sm.PreviousMatch()
	if match == nil || sm.GetCurrentIndex() != 2 {
		t.Errorf("Wrap around backward failed, index: %d", sm.GetCurrentIndex())
	}
}

func TestSearchManager_Replace(t *testing.T) {
	sm := NewSearchManager()
	text := "Hello World\nHello Universe"

	// Test replace all
	options := ReplaceOptions{
		SearchOptions: SearchOptions{
			CaseSensitive:     false,
			WholeWord:         false,
			RegularExpression: false,
			WrapAround:        true,
		},
		ReplaceAll: true,
	}

	result, count := sm.Replace(text, "Hello", "Hi", options)
	if count != 2 {
		t.Errorf("Replace all: Expected 2 replacements, got %d", count)
	}

	expected := "Hi World\nHi Universe"
	if result != expected {
		t.Errorf("Replace all result incorrect:\nExpected: %s\nGot: %s", expected, result)
	}
}

func TestSearchManager_ReplaceCurrent(t *testing.T) {
	sm := NewSearchManager()
	text := "Hello World\nHello Universe"

	// Find matches first
	sm.Find(text, "Hello")
	sm.SetCurrentMatch(0)

	// Test replace current
	options := ReplaceOptions{
		SearchOptions: SearchOptions{
			CaseSensitive:     false,
			WholeWord:         false,
			RegularExpression: false,
			WrapAround:        true,
		},
		ReplaceAll: false,
	}

	result, count := sm.Replace(text, "Hello", "Hi", options)
	if count != 1 {
		t.Errorf("Replace current: Expected 1 replacement, got %d", count)
	}

	expected := "Hi World\nHello Universe"
	if result != expected {
		t.Errorf("Replace current result incorrect:\nExpected: %s\nGot: %s", expected, result)
	}
}

func TestSearchManager_RegexReplace(t *testing.T) {
	sm := NewSearchManager()
	text := "test123 test456"

	options := ReplaceOptions{
		SearchOptions: SearchOptions{
			CaseSensitive:     false,
			WholeWord:         false,
			RegularExpression: true,
			WrapAround:        true,
		},
		ReplaceAll: true,
	}

	result, count := sm.Replace(text, `test(\d+)`, "number$1", options)
	if count != 2 {
		t.Errorf("Regex replace: Expected 2 replacements, got %d", count)
	}

	expected := "number123 number456"
	if result != expected {
		t.Errorf("Regex replace result incorrect:\nExpected: %s\nGot: %s", expected, result)
	}
}

func TestSearchManager_EmptyPattern(t *testing.T) {
	sm := NewSearchManager()
	text := "Hello World"

	matches := sm.Find(text, "")
	if len(matches) != 0 {
		t.Errorf("Empty pattern should return no matches, got %d", len(matches))
	}
}

func TestSearchManager_Clear(t *testing.T) {
	sm := NewSearchManager()
	text := "Hello World"

	sm.Find(text, "Hello")
	if len(sm.GetMatches()) == 0 {
		t.Error("Should have matches before clear")
	}

	sm.Clear()
	if len(sm.GetMatches()) != 0 {
		t.Error("Should have no matches after clear")
	}

	if sm.GetCurrentIndex() != -1 {
		t.Error("Current index should be -1 after clear")
	}
}