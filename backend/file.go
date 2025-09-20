package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo holds information about a file
type FileInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"modTime"`
	IsDirectory  bool      `json:"isDirectory"`
	Permissions  string    `json:"permissions"`
	Extension    string    `json:"extension"`
}

// FileError represents a file operation error
type FileError struct {
	Operation string
	Path      string
	Err       error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("文件%s失败 '%s': %v", e.Operation, e.Path, e.Err)
}

// FileManager handles file operations
type FileManager struct {
	fileTypeManager *FileTypeManager
}

// NewFileManager creates a new file manager
func NewFileManager() *FileManager {
	return &FileManager{
		fileTypeManager: NewFileTypeManager(),
	}
}

// ReadFile reads a file and returns its content
func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// SaveFile saves content to a file
func SaveFile(path, content string) {
	_ = os.WriteFile(path, []byte(content), 0644)
}

// ReadFileWithInfo reads a file and returns content with file information
func (fm *FileManager) ReadFileWithInfo(path string) (string, *FileInfo, error) {
	// Get file info first
	info, err := fm.GetFileInfo(path)
	if err != nil {
		return "", nil, &FileError{
			Operation: "读取",
			Path:      path,
			Err:       err,
		}
	}

	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return "", info, &FileError{
			Operation: "读取",
			Path:      path,
			Err:       err,
		}
	}

	return string(data), info, nil
}

// SaveFileWithBackup saves content to a file with backup
func (fm *FileManager) SaveFileWithBackup(path, content string) error {
	// Create backup if file exists
	if _, err := os.Stat(path); err == nil {
		backupPath := path + ".bak"
		if err := fm.copyFile(path, backupPath); err != nil {
			// Log warning but continue with save
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	// Save the file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return &FileError{
			Operation: "保存",
			Path:      path,
			Err:       err,
		}
	}

	return nil
}

// GetFileInfo returns information about a file
func (fm *FileManager) GetFileInfo(path string) (*FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	return &FileInfo{
		Path:        absPath,
		Name:        stat.Name(),
		Size:        stat.Size(),
		ModTime:     stat.ModTime(),
		IsDirectory: stat.IsDir(),
		Permissions: stat.Mode().String(),
		Extension:   strings.ToLower(filepath.Ext(path)),
	}, nil
}

// FileExists checks if a file exists
func (fm *FileManager) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsReadable checks if a file is readable
func (fm *FileManager) IsReadable(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// IsWritable checks if a file is writable
func (fm *FileManager) IsWritable(path string) bool {
	// Try to open file for writing
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		// If file doesn't exist, check if we can create it
		if os.IsNotExist(err) {
			dir := filepath.Dir(path)
			return fm.isDirWritable(dir)
		}
		return false
	}
	file.Close()
	return true
}

// GetFileType returns the file type for a given path
func (fm *FileManager) GetFileType(path string) FileType {
	return fm.fileTypeManager.DetectFileType(path)
}

// GetLexerName returns the lexer name for syntax highlighting
func (fm *FileManager) GetLexerName(path string) string {
	return fm.fileTypeManager.GetLexerName(path)
}

// ValidateFilePath validates if a file path is valid
func (fm *FileManager) ValidateFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// Check for invalid characters (basic validation)
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("文件路径包含无效字符: %s", char)
		}
	}

	return nil
}

// CreateDirectory creates a directory if it doesn't exist
func (fm *FileManager) CreateDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return &FileError{
			Operation: "创建目录",
			Path:      path,
			Err:       err,
		}
	}
	return nil
}

// copyFile copies a file from src to dst
func (fm *FileManager) copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// isDirWritable checks if a directory is writable
func (fm *FileManager) isDirWritable(dir string) bool {
	testFile := filepath.Join(dir, ".write_test")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(testFile)
	return true
}

// GetRecentFiles returns a list of recently opened files
func (fm *FileManager) GetRecentFiles(maxCount int) ([]string, error) {
	// This will be implemented when we add recent files functionality
	// For now, return empty list
	return []string{}, nil
}
