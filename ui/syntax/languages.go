package syntax

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
)

// LanguageManager manages syntax highlighting for different programming languages
type LanguageManager struct {
	lexers map[string]chroma.Lexer
	cache  map[string]chroma.Lexer
}

// NewLanguageManager creates a new language manager
func NewLanguageManager() *LanguageManager {
	return &LanguageManager{
		lexers: make(map[string]chroma.Lexer),
		cache:  make(map[string]chroma.Lexer),
	}
}

// GetLexer returns the appropriate lexer for the given language
func (lm *LanguageManager) GetLexer(language string) chroma.Lexer {
	// Check cache first
	if lexer, exists := lm.cache[language]; exists {
		return lexer
	}

	// Try to get lexer by name
	lexer := lexers.Get(language)
	if lexer == nil {
		// Try common aliases
		lexer = lm.getLexerByAlias(language)
	}

	// Fall back to plain text if no lexer found
	if lexer == nil {
		log.Printf("No lexer found for language: %s, using fallback", language)
		lexer = lexers.Fallback
	}

	// Cache the result
	lm.cache[language] = lexer
	return lexer
}

// getLexerByAlias tries to find a lexer using common aliases
func (lm *LanguageManager) getLexerByAlias(language string) chroma.Lexer {
	aliases := map[string]string{
		"js":         "javascript",
		"jsx":        "javascript",
		"ts":         "typescript",
		"tsx":        "typescript",
		"py":         "python",
		"pyw":        "python",
		"pyi":        "python",
		"c++":        "cpp",
		"cxx":        "cpp",
		"cc":         "cpp",
		"hpp":        "cpp",
		"hxx":        "cpp",
		"h":          "c",
		"htm":        "html",
		"yml":        "yaml",
		"sh":         "bash",
		"zsh":        "bash",
		"dockerfile": "docker",
		"text":       "text",
		"txt":        "text",
		"plain":      "text",
	}

	if alias, exists := aliases[strings.ToLower(language)]; exists {
		return lexers.Get(alias)
	}

	return nil
}

// GetSupportedLanguages returns a list of all supported languages
func (lm *LanguageManager) GetSupportedLanguages() []string {
	var languages []string
	
	// Get all available lexers from chroma
	for _, lexer := range lexers.Registry.Lexers {
		if lexer != nil {
			config := lexer.Config()
			if config != nil {
				languages = append(languages, config.Name)
				// Add aliases as well
				for _, alias := range config.Aliases {
					languages = append(languages, alias)
				}
			}
		}
	}

	return languages
}

// IsLanguageSupported checks if a language is supported
func (lm *LanguageManager) IsLanguageSupported(language string) bool {
	lexer := lm.GetLexer(language)
	return lexer != lexers.Fallback
}

// GetLanguageInfo returns information about a language
func (lm *LanguageManager) GetLanguageInfo(language string) *LanguageInfo {
	lexer := lm.GetLexer(language)
	if lexer == nil {
		return nil
	}

	config := lexer.Config()
	if config == nil {
		return nil
	}

	return &LanguageInfo{
		Name:        config.Name,
		Aliases:     config.Aliases,
		Filenames:   config.Filenames,
		MimeTypes:   config.MimeTypes,
		Description: fmt.Sprintf("Syntax highlighting for %s", config.Name),
	}
}

// LanguageInfo contains information about a programming language
type LanguageInfo struct {
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases"`
	Filenames   []string `json:"filenames"`
	MimeTypes   []string `json:"mimeTypes"`
	Description string   `json:"description"`
}

// DetectLanguageFromFilename attempts to detect the language from a filename
func (lm *LanguageManager) DetectLanguageFromFilename(filename string) string {
	// Try to match by filename patterns first
	for _, lexer := range lexers.Registry.Lexers {
		if lexer != nil {
			config := lexer.Config()
			if config != nil {
				for _, pattern := range config.Filenames {
					if matched, _ := filepath.Match(pattern, filename); matched {
						// Return the first alias if available, otherwise the name
						if len(config.Aliases) > 0 {
							return config.Aliases[0]
						}
						return strings.ToLower(config.Name)
					}
				}
			}
		}
	}

	// Fall back to extension-based detection
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".go":
		return "go"
	case ".js", ".jsx", ".mjs":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".py", ".pyw", ".pyi":
		return "python"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cpp", ".cxx", ".cc", ".hpp", ".hxx":
		return "cpp"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".json":
		return "json"
	case ".xml", ".xsd", ".xsl":
		return "xml"
	case ".yaml", ".yml":
		return "yaml"
	case ".md", ".markdown":
		return "markdown"
	case ".sh", ".bash", ".zsh":
		return "bash"
	case ".sql":
		return "sql"
	case ".dockerfile":
		return "docker"
	default:
		return "text"
	}
}

// ClearCache clears the lexer cache
func (lm *LanguageManager) ClearCache() {
	lm.cache = make(map[string]chroma.Lexer)
}