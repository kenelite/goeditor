package syntax

import (
	"image/color"
	"log"
	"sync"

	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/chroma"
)

var (
	// Global managers - initialized once
	languageManager *LanguageManager
	themeManager    *ThemeManager
	managerOnce     sync.Once
)

// initManagers initializes the global managers
func initManagers() {
	managerOnce.Do(func() {
		languageManager = NewLanguageManager()
		themeManager = NewThemeManager()
	})
}

// GetLanguageManager returns the global language manager
func GetLanguageManager() *LanguageManager {
	initManagers()
	return languageManager
}

// GetThemeManager returns the global theme manager
func GetThemeManager() *ThemeManager {
	initManagers()
	return themeManager
}

// chromaToRGBA converts chroma.Colour to color.Color
func chromaToRGBA(c chroma.Colour) color.Color {
	if !c.IsSet() {
		return color.Black
	}
	r := uint8((c >> 16) & 0xFF)
	g := uint8((c >> 8) & 0xFF)
	b := uint8(c & 0xFF)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// HighlightCode highlights source code for any supported language
func HighlightCode(source string, language string) []widget.RichTextSegment {
	return HighlightCodeWithTheme(source, language, "")
}

// HighlightCodeWithTheme highlights source code with a specific theme
func HighlightCodeWithTheme(source string, language string, themeName string) []widget.RichTextSegment {
	initManagers()
	
	// Handle empty source
	if source == "" {
		return []widget.RichTextSegment{}
	}

	// Get lexer for the language
	lexer := languageManager.GetLexer(language)
	if lexer == nil {
		log.Printf("No lexer found for language: %s", language)
		return createPlainTextSegments(source)
	}

	// Get theme
	var style *chroma.Style
	if themeName != "" {
		// Try to use specific theme
		if err := themeManager.SetTheme(themeName); err != nil {
			log.Printf("Failed to set theme %s: %v", themeName, err)
		}
	}
	style = themeManager.GetTheme()

	// Tokenize the source code
	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		log.Printf("Tokenization error for language %s: %v", language, err)
		return createPlainTextSegments(source)
	}

	// Convert tokens to segments
	var segments []widget.RichTextSegment
	for token := iterator(); token != chroma.EOF; token = iterator() {
		if token.Value == "" {
			continue
		}

		// Get color for this token type
		tokenStyle := style.Get(token.Type)
		var col color.Color
		if tokenStyle.Colour.IsSet() {
			col = chromaToRGBA(tokenStyle.Colour)
		} else {
			// Use default foreground color
			col = themeManager.GetForegroundColor()
		}

		seg := NewSyntaxSegment(token.Value, col)
		segments = append(segments, seg)
	}

	// If no segments were created, return plain text
	if len(segments) == 0 {
		return createPlainTextSegments(source)
	}

	return segments
}

// createPlainTextSegments creates plain text segments as fallback
func createPlainTextSegments(source string) []widget.RichTextSegment {
	initManagers()
	col := themeManager.GetForegroundColor()
	return []widget.RichTextSegment{
		NewSyntaxSegment(source, col),
	}
}

// HighlightGoCode maintains backward compatibility
func HighlightGoCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "go")
}

// HighlightJavaScriptCode highlights JavaScript code
func HighlightJavaScriptCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "javascript")
}

// HighlightPythonCode highlights Python code
func HighlightPythonCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "python")
}

// HighlightJavaCode highlights Java code
func HighlightJavaCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "java")
}

// HighlightHTMLCode highlights HTML code
func HighlightHTMLCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "html")
}

// HighlightCSSCode highlights CSS code
func HighlightCSSCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "css")
}

// HighlightJSONCode highlights JSON code
func HighlightJSONCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "json")
}

// HighlightXMLCode highlights XML code
func HighlightXMLCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "xml")
}

// HighlightYAMLCode highlights YAML code
func HighlightYAMLCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "yaml")
}

// HighlightMarkdownCode highlights Markdown code
func HighlightMarkdownCode(source string) []widget.RichTextSegment {
	return HighlightCode(source, "markdown")
}

// DetectAndHighlight detects language from filename and highlights accordingly
func DetectAndHighlight(source string, filename string) []widget.RichTextSegment {
	initManagers()
	language := languageManager.DetectLanguageFromFilename(filename)
	return HighlightCode(source, language)
}

// SetTheme sets the global syntax highlighting theme
func SetTheme(themeName string) error {
	initManagers()
	return themeManager.SetTheme(themeName)
}

// GetCurrentTheme returns the current theme name
func GetCurrentTheme() string {
	initManagers()
	return themeManager.GetCurrentThemeName()
}

// GetAvailableThemes returns all available theme names
func GetAvailableThemes() []string {
	initManagers()
	return themeManager.GetAvailableThemes()
}

// GetSupportedLanguages returns all supported languages
func GetSupportedLanguages() []string {
	initManagers()
	return languageManager.GetSupportedLanguages()
}

// IsLanguageSupported checks if a language is supported
func IsLanguageSupported(language string) bool {
	initManagers()
	return languageManager.IsLanguageSupported(language)
}
