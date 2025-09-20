package syntax

import (
	"strings"
	"testing"
	
	"fyne.io/fyne/v2/widget"
)

func TestHighlightCode(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		language string
	}{
		{
			name:     "Go code",
			source:   `package main\n\nfunc main() {\n\tfmt.Println("Hello, World!")\n}`,
			language: "go",
		},
		{
			name:     "JavaScript code",
			source:   `function hello() {\n\tconsole.log("Hello, World!");\n}`,
			language: "javascript",
		},
		{
			name:     "Python code",
			source:   `def hello():\n\tprint("Hello, World!")`,
			language: "python",
		},
		{
			name:     "Java code",
			source:   `public class Hello {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println("Hello, World!");\n\t}\n}`,
			language: "java",
		},
		{
			name:     "HTML code",
			source:   `<html><body><h1>Hello, World!</h1></body></html>`,
			language: "html",
		},
		{
			name:     "CSS code",
			source:   `body { color: red; font-size: 14px; }`,
			language: "css",
		},
		{
			name:     "JSON code",
			source:   `{"name": "test", "value": 123}`,
			language: "json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			segments := HighlightCode(tc.source, tc.language)
			
			if len(segments) == 0 {
				t.Error("Expected at least one segment")
				return
			}

			// Verify that all text is preserved
			var reconstructed strings.Builder
			for _, segment := range segments {
				reconstructed.WriteString(segment.Textual())
			}

			if reconstructed.String() != tc.source {
				t.Errorf("Text not preserved. Expected:\n%s\nGot:\n%s", tc.source, reconstructed.String())
			}
		})
	}
}

func TestHighlightCodeWithTheme(t *testing.T) {
	source := `package main\n\nfunc main() {\n\tfmt.Println("Hello")\n}`
	
	// Test with different themes
	themes := []string{"github", "monokai", "vs"}
	
	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			segments := HighlightCodeWithTheme(source, "go", theme)
			
			if len(segments) == 0 {
				t.Error("Expected at least one segment")
				return
			}

			// Verify text preservation
			var reconstructed strings.Builder
			for _, segment := range segments {
				reconstructed.WriteString(segment.Textual())
			}

			if reconstructed.String() != source {
				t.Error("Text not preserved with theme")
			}
		})
	}
}

func TestHighlightEmptySource(t *testing.T) {
	segments := HighlightCode("", "go")
	if len(segments) != 0 {
		t.Errorf("Expected 0 segments for empty source, got %d", len(segments))
	}
}

func TestHighlightUnsupportedLanguage(t *testing.T) {
	source := "some text content"
	segments := HighlightCode(source, "nonexistent-language")
	
	if len(segments) == 0 {
		t.Error("Expected at least one segment for unsupported language")
		return
	}

	// Should return plain text
	if segments[0].Textual() != source {
		t.Error("Expected plain text fallback for unsupported language")
	}
}

func TestDetectAndHighlight(t *testing.T) {
	testCases := []struct {
		filename string
		source   string
	}{
		{"main.go", `package main`},
		{"script.js", `function test() {}`},
		{"app.py", `def test():`},
		{"style.css", `body { color: red; }`},
		{"data.json", `{"key": "value"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			segments := DetectAndHighlight(tc.source, tc.filename)
			
			if len(segments) == 0 {
				t.Error("Expected at least one segment")
				return
			}

			// Verify text preservation
			var reconstructed strings.Builder
			for _, segment := range segments {
				reconstructed.WriteString(segment.Textual())
			}

			if reconstructed.String() != tc.source {
				t.Error("Text not preserved in detect and highlight")
			}
		})
	}
}

func TestBackwardCompatibility(t *testing.T) {
	source := `package main\n\nfunc main() {\n\tfmt.Println("Hello")\n}`
	
	// Test that old HighlightGoCode still works
	segments := HighlightGoCode(source)
	
	if len(segments) == 0 {
		t.Error("Expected at least one segment from HighlightGoCode")
		return
	}

	// Verify text preservation
	var reconstructed strings.Builder
	for _, segment := range segments {
		reconstructed.WriteString(segment.Textual())
	}

	if reconstructed.String() != source {
		t.Error("HighlightGoCode backward compatibility broken")
	}
}

func TestLanguageSpecificFunctions(t *testing.T) {
	source := "test content"
	
	testCases := []struct {
		name string
		fn   func(string) []widget.RichTextSegment
	}{
		{"JavaScript", HighlightJavaScriptCode},
		{"Python", HighlightPythonCode},
		{"Java", HighlightJavaCode},
		{"HTML", HighlightHTMLCode},
		{"CSS", HighlightCSSCode},
		{"JSON", HighlightJSONCode},
		{"XML", HighlightXMLCode},
		{"YAML", HighlightYAMLCode},
		{"Markdown", HighlightMarkdownCode},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			segments := tc.fn(source)
			
			if len(segments) == 0 {
				t.Errorf("Expected at least one segment from %s function", tc.name)
				return
			}

			// Verify text preservation
			var reconstructed strings.Builder
			for _, segment := range segments {
				reconstructed.WriteString(segment.Textual())
			}

			if reconstructed.String() != source {
				t.Errorf("%s function did not preserve text", tc.name)
			}
		})
	}
}

func TestThemeManagement(t *testing.T) {
	// Test setting theme
	err := SetTheme("monokai")
	if err != nil {
		t.Errorf("Expected no error setting theme, got: %v", err)
	}

	// Test getting current theme
	current := GetCurrentTheme()
	if current != "monokai" {
		t.Errorf("Expected current theme to be 'monokai', got '%s'", current)
	}

	// Test getting available themes
	themes := GetAvailableThemes()
	if len(themes) == 0 {
		t.Error("Expected at least some available themes")
	}
}

func TestLanguageSupport(t *testing.T) {
	// Test getting supported languages
	languages := GetSupportedLanguages()
	if len(languages) == 0 {
		t.Error("Expected at least some supported languages")
	}

	// Test language support check
	if !IsLanguageSupported("go") {
		t.Error("Expected Go to be supported")
	}

	if IsLanguageSupported("nonexistent-language-xyz") {
		t.Error("Expected nonexistent language to not be supported")
	}
}

func TestCreatePlainTextSegments(t *testing.T) {
	source := "plain text content"
	segments := createPlainTextSegments(source)
	
	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].Textual() != source {
		t.Error("Plain text segment content mismatch")
	}
}