package backend

import (
	"testing"
)

func TestFileTypeDetection(t *testing.T) {
	ftm := NewFileTypeManager()
	
	testCases := []struct {
		filename     string
		expectedName string
		expectedLexer string
	}{
		{"main.go", "Go", "go"},
		{"script.js", "JavaScript", "javascript"},
		{"app.tsx", "TypeScript", "typescript"},
		{"script.py", "Python", "python"},
		{"Main.java", "Java", "java"},
		{"index.html", "HTML", "html"},
		{"style.css", "CSS", "css"},
		{"data.json", "JSON", "json"},
		{"config.xml", "XML", "xml"},
		{"docker-compose.yml", "YAML", "yaml"},
		{"README.md", "Markdown", "markdown"},
		{"script.sh", "Shell", "bash"},
		{"unknown.xyz", "Plain Text", "text"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			fileType := ftm.DetectFileType(tc.filename)
			
			if fileType.Name != tc.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tc.expectedName, fileType.Name)
			}
			
			if fileType.LexerName != tc.expectedLexer {
				t.Errorf("Expected lexer '%s', got '%s'", tc.expectedLexer, fileType.LexerName)
			}
		})
	}
}

func TestFileTypeByExtension(t *testing.T) {
	ftm := NewFileTypeManager()
	
	// Test with dot
	fileType, exists := ftm.GetFileTypeByExtension(".go")
	if !exists {
		t.Error("Expected .go extension to exist")
	}
	if fileType.Name != "Go" {
		t.Errorf("Expected Go, got %s", fileType.Name)
	}
	
	// Test without dot
	fileType, exists = ftm.GetFileTypeByExtension("js")
	if !exists {
		t.Error("Expected js extension to exist")
	}
	if fileType.Name != "JavaScript" {
		t.Errorf("Expected JavaScript, got %s", fileType.Name)
	}
	
	// Test non-existent extension
	_, exists = ftm.GetFileTypeByExtension(".nonexistent")
	if exists {
		t.Error("Expected .nonexistent extension to not exist")
	}
}

func TestIsSupported(t *testing.T) {
	ftm := NewFileTypeManager()
	
	supportedFiles := []string{
		"main.go",
		"script.js",
		"app.py",
		"style.css",
		"data.json",
	}
	
	for _, filename := range supportedFiles {
		if !ftm.IsSupported(filename) {
			t.Errorf("Expected %s to be supported", filename)
		}
	}
	
	unsupportedFiles := []string{
		"binary.exe",
		"image.png",
		"document.pdf",
	}
	
	for _, filename := range unsupportedFiles {
		if ftm.IsSupported(filename) {
			t.Errorf("Expected %s to not be supported", filename)
		}
	}
}

func TestGetLexerName(t *testing.T) {
	ftm := NewFileTypeManager()
	
	testCases := []struct {
		filename string
		expected string
	}{
		{"main.go", "go"},
		{"script.js", "javascript"},
		{"app.py", "python"},
		{"unknown.xyz", "text"},
	}
	
	for _, tc := range testCases {
		lexer := ftm.GetLexerName(tc.filename)
		if lexer != tc.expected {
			t.Errorf("Expected lexer '%s' for %s, got '%s'", tc.expected, tc.filename, lexer)
		}
	}
}

func TestGetAllFileTypes(t *testing.T) {
	ftm := NewFileTypeManager()
	
	fileTypes := ftm.GetAllFileTypes()
	
	// Should have at least the default types
	if len(fileTypes) < 10 {
		t.Errorf("Expected at least 10 file types, got %d", len(fileTypes))
	}
	
	// Check for some expected types
	expectedTypes := []string{"Go", "JavaScript", "Python", "Java", "HTML", "CSS"}
	typeNames := make(map[string]bool)
	
	for _, ft := range fileTypes {
		typeNames[ft.Name] = true
	}
	
	for _, expected := range expectedTypes {
		if !typeNames[expected] {
			t.Errorf("Expected file type '%s' to be present", expected)
		}
	}
}