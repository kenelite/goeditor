package syntax

import (
	"strings"
	"testing"
)

// TestMultiLanguageIntegration tests the complete multi-language syntax highlighting system
func TestMultiLanguageIntegration(t *testing.T) {
	// Test samples for different languages
	samples := map[string]string{
		"go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		"javascript": `function greet(name) {
	console.log("Hello, " + name + "!");
}

greet("World");`,
		"python": `def greet(name):
	print(f"Hello, {name}!")

if __name__ == "__main__":
	greet("World")`,
		"java": `public class HelloWorld {
	public static void main(String[] args) {
		System.out.println("Hello, World!");
	}
}`,
		"html": `<!DOCTYPE html>
<html>
<head>
	<title>Hello World</title>
</head>
<body>
	<h1>Hello, World!</h1>
</body>
</html>`,
		"css": `body {
	font-family: Arial, sans-serif;
	background-color: #f0f0f0;
	color: #333;
}

h1 {
	color: #007acc;
	text-align: center;
}`,
		"json": `{
	"name": "Hello World",
	"version": "1.0.0",
	"description": "A simple hello world example",
	"main": "index.js",
	"keywords": ["hello", "world", "example"]
}`,
	}

	// Test each language
	for language, code := range samples {
		t.Run(language, func(t *testing.T) {
			// Test basic highlighting
			segments := HighlightCode(code, language)
			if len(segments) == 0 {
				t.Errorf("Expected segments for %s code", language)
				return
			}

			// Verify text preservation
			var reconstructed string
			for _, segment := range segments {
				reconstructed += segment.Textual()
			}

			if reconstructed != code {
				t.Errorf("Text not preserved for %s", language)
			}

			// Test that we get more than just plain text (indicating actual highlighting)
			if len(segments) == 1 {
				t.Logf("Warning: Only one segment for %s (might be plain text)", language)
			}
		})
	}
}

// TestThemeIntegration tests theme switching with different languages
func TestThemeIntegration(t *testing.T) {
	code := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`

	themes := []string{"github", "monokai", "vs", "dracula"}
	
	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			segments := HighlightCodeWithTheme(code, "go", theme)
			
			if len(segments) == 0 {
				t.Errorf("Expected segments with theme %s", theme)
				return
			}

			// Verify text preservation
			var reconstructed string
			for _, segment := range segments {
				reconstructed += segment.Textual()
			}

			if reconstructed != code {
				t.Errorf("Text not preserved with theme %s", theme)
			}
		})
	}
}

// TestFileDetectionIntegration tests automatic language detection from filenames
func TestFileDetectionIntegration(t *testing.T) {
	testFiles := map[string]string{
		"main.go":           "package main\n\nfunc main() {}",
		"script.js":         "function test() { return true; }",
		"app.py":            "def hello():\n    print('Hello')",
		"Component.tsx":     "const App = () => <div>Hello</div>;",
		"style.css":         "body { margin: 0; }",
		"data.json":         `{"key": "value"}`,
		"config.xml":        "<root><item>value</item></root>",
		"docker-compose.yml": "version: '3'\nservices:\n  app:",
		"README.md":         "# Hello\n\nThis is a test.",
	}

	for filename, code := range testFiles {
		t.Run(filename, func(t *testing.T) {
			segments := DetectAndHighlight(code, filename)
			
			if len(segments) == 0 {
				t.Errorf("Expected segments for file %s", filename)
				return
			}

			// Verify text preservation
			var reconstructed string
			for _, segment := range segments {
				reconstructed += segment.Textual()
			}

			if reconstructed != code {
				t.Errorf("Text not preserved for file %s", filename)
			}
		})
	}
}

// TestPerformanceBaseline provides a basic performance test
func TestPerformanceBaseline(t *testing.T) {
	// Large code sample
	largeCode := ""
	for i := 0; i < 100; i++ {
		largeCode += `package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No arguments provided")
		return
	}
	
	for i, arg := range args {
		fmt.Printf("Arg %d: %s\n", i, strings.ToUpper(arg))
	}
}

`
	}

	// Test highlighting performance
	segments := HighlightCode(largeCode, "go")
	
	if len(segments) == 0 {
		t.Error("Expected segments for large code")
		return
	}

	// Verify text preservation
	var reconstructed string
	for _, segment := range segments {
		reconstructed += segment.Textual()
	}

	if reconstructed != largeCode {
		t.Error("Text not preserved for large code")
	}

	t.Logf("Successfully highlighted %d characters in %d segments", len(largeCode), len(segments))
}

// TestErrorHandling tests error handling and fallback behavior
func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		language string
	}{
		{"empty code", "", "go"},
		{"invalid language", "some code", "invalid-language-xyz"},
		{"malformed code", "this is not valid go code }{][", "go"},
		{"very long line", strings.Repeat("a", 10000), "text"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			segments := HighlightCode(tc.code, tc.language)
			
			// For empty code, expect empty segments
			if tc.code == "" {
				if len(segments) != 0 {
					t.Errorf("Expected 0 segments for empty code, got %d", len(segments))
				}
				return
			}

			// For non-empty code, should get at least one segment
			if len(segments) == 0 {
				t.Errorf("Expected at least one segment for non-empty code")
				return
			}

			// Verify text preservation
			var reconstructed string
			for _, segment := range segments {
				reconstructed += segment.Textual()
			}

			if reconstructed != tc.code {
				t.Errorf("Text not preserved in error case %s", tc.name)
			}
		})
	}
}