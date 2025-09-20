package backend

import (
	"path/filepath"
	"strings"
)

// FileType represents a supported file type
type FileType struct {
	Name       string   `json:"name"`
	Extensions []string `json:"extensions"`
	LexerName  string   `json:"lexerName"`
	MimeType   string   `json:"mimeType"`
}

// FileTypeManager manages file type detection and handling
type FileTypeManager struct {
	fileTypes map[string]FileType
}

// NewFileTypeManager creates a new file type manager
func NewFileTypeManager() *FileTypeManager {
	manager := &FileTypeManager{
		fileTypes: make(map[string]FileType),
	}
	
	// Register default file types
	manager.registerDefaultFileTypes()
	
	return manager
}

// registerDefaultFileTypes registers the default supported file types
func (ftm *FileTypeManager) registerDefaultFileTypes() {
	defaultTypes := []FileType{
		{
			Name:       "Go",
			Extensions: []string{".go"},
			LexerName:  "go",
			MimeType:   "text/x-go",
		},
		{
			Name:       "JavaScript",
			Extensions: []string{".js", ".jsx", ".mjs"},
			LexerName:  "javascript",
			MimeType:   "application/javascript",
		},
		{
			Name:       "TypeScript",
			Extensions: []string{".ts", ".tsx"},
			LexerName:  "typescript",
			MimeType:   "application/typescript",
		},
		{
			Name:       "Python",
			Extensions: []string{".py", ".pyw", ".pyi"},
			LexerName:  "python",
			MimeType:   "text/x-python",
		},
		{
			Name:       "Java",
			Extensions: []string{".java"},
			LexerName:  "java",
			MimeType:   "text/x-java-source",
		},
		{
			Name:       "C",
			Extensions: []string{".c", ".h"},
			LexerName:  "c",
			MimeType:   "text/x-c",
		},
		{
			Name:       "C++",
			Extensions: []string{".cpp", ".cxx", ".cc", ".hpp", ".hxx"},
			LexerName:  "cpp",
			MimeType:   "text/x-c++",
		},
		{
			Name:       "HTML",
			Extensions: []string{".html", ".htm"},
			LexerName:  "html",
			MimeType:   "text/html",
		},
		{
			Name:       "CSS",
			Extensions: []string{".css"},
			LexerName:  "css",
			MimeType:   "text/css",
		},
		{
			Name:       "JSON",
			Extensions: []string{".json"},
			LexerName:  "json",
			MimeType:   "application/json",
		},
		{
			Name:       "XML",
			Extensions: []string{".xml", ".xsd", ".xsl"},
			LexerName:  "xml",
			MimeType:   "application/xml",
		},
		{
			Name:       "YAML",
			Extensions: []string{".yaml", ".yml"},
			LexerName:  "yaml",
			MimeType:   "application/x-yaml",
		},
		{
			Name:       "Markdown",
			Extensions: []string{".md", ".markdown"},
			LexerName:  "markdown",
			MimeType:   "text/markdown",
		},
		{
			Name:       "Shell",
			Extensions: []string{".sh", ".bash", ".zsh"},
			LexerName:  "bash",
			MimeType:   "application/x-sh",
		},
		{
			Name:       "SQL",
			Extensions: []string{".sql"},
			LexerName:  "sql",
			MimeType:   "application/sql",
		},
		{
			Name:       "Dockerfile",
			Extensions: []string{".dockerfile"},
			LexerName:  "dockerfile",
			MimeType:   "text/x-dockerfile",
		},
		{
			Name:       "Plain Text",
			Extensions: []string{".txt", ".text"},
			LexerName:  "text",
			MimeType:   "text/plain",
		},
	}

	for _, fileType := range defaultTypes {
		ftm.RegisterFileType(fileType)
	}
}

// RegisterFileType registers a new file type
func (ftm *FileTypeManager) RegisterFileType(fileType FileType) {
	for _, ext := range fileType.Extensions {
		ftm.fileTypes[strings.ToLower(ext)] = fileType
	}
}

// DetectFileType detects the file type based on file extension
func (ftm *FileTypeManager) DetectFileType(filename string) FileType {
	ext := strings.ToLower(filepath.Ext(filename))
	
	if fileType, exists := ftm.fileTypes[ext]; exists {
		return fileType
	}
	
	// Return default plain text type if not found
	return FileType{
		Name:      "Plain Text",
		Extensions: []string{ext},
		LexerName: "text",
		MimeType:  "text/plain",
	}
}

// GetFileTypeByExtension returns the file type for a given extension
func (ftm *FileTypeManager) GetFileTypeByExtension(extension string) (FileType, bool) {
	ext := strings.ToLower(extension)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	
	fileType, exists := ftm.fileTypes[ext]
	return fileType, exists
}

// GetSupportedExtensions returns all supported file extensions
func (ftm *FileTypeManager) GetSupportedExtensions() []string {
	extensions := make([]string, 0, len(ftm.fileTypes))
	for ext := range ftm.fileTypes {
		extensions = append(extensions, ext)
	}
	return extensions
}

// GetAllFileTypes returns all registered file types
func (ftm *FileTypeManager) GetAllFileTypes() []FileType {
	typeMap := make(map[string]FileType)
	
	// Deduplicate by name
	for _, fileType := range ftm.fileTypes {
		typeMap[fileType.Name] = fileType
	}
	
	types := make([]FileType, 0, len(typeMap))
	for _, fileType := range typeMap {
		types = append(types, fileType)
	}
	
	return types
}

// IsSupported checks if a file extension is supported
func (ftm *FileTypeManager) IsSupported(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	_, exists := ftm.fileTypes[ext]
	return exists
}

// GetLexerName returns the lexer name for a file
func (ftm *FileTypeManager) GetLexerName(filename string) string {
	fileType := ftm.DetectFileType(filename)
	return fileType.LexerName
}

// GetMimeType returns the MIME type for a file
func (ftm *FileTypeManager) GetMimeType(filename string) string {
	fileType := ftm.DetectFileType(filename)
	return fileType.MimeType
}