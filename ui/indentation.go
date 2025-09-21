package ui

import (
	"strings"
	"unicode"
)

// IndentationManager handles tab and indentation operations
type IndentationManager struct {
	tabSize      int
	useSpaces    bool
	autoIndent   bool
}

// NewIndentationManager creates a new indentation manager
func NewIndentationManager() *IndentationManager {
	return &IndentationManager{
		tabSize:    4,     // Default tab size
		useSpaces:  true,  // Use spaces by default
		autoIndent: true,  // Auto-indent by default
	}
}

// SetTabSize sets the tab size
func (im *IndentationManager) SetTabSize(size int) {
	if size > 0 && size <= 16 {
		im.tabSize = size
	}
}

// GetTabSize returns the current tab size
func (im *IndentationManager) GetTabSize() int {
	return im.tabSize
}

// SetUseSpaces sets whether to use spaces instead of tabs
func (im *IndentationManager) SetUseSpaces(useSpaces bool) {
	im.useSpaces = useSpaces
}

// GetUseSpaces returns whether spaces are used instead of tabs
func (im *IndentationManager) GetUseSpaces() bool {
	return im.useSpaces
}

// SetAutoIndent sets whether to automatically indent new lines
func (im *IndentationManager) SetAutoIndent(autoIndent bool) {
	im.autoIndent = autoIndent
}

// GetAutoIndent returns whether auto-indent is enabled
func (im *IndentationManager) GetAutoIndent() bool {
	return im.autoIndent
}

// GetIndentString returns the string to use for one level of indentation
func (im *IndentationManager) GetIndentString() string {
	if im.useSpaces {
		return strings.Repeat(" ", im.tabSize)
	}
	return "\t"
}

// HandleTabKey handles the Tab key press
func (im *IndentationManager) HandleTabKey(content string, cursorPos int) (newContent string, newCursorPos int) {
	if content == "" {
		return im.GetIndentString(), len(im.GetIndentString())
	}
	
	// Insert indentation at cursor position
	indentStr := im.GetIndentString()
	newContent = content[:cursorPos] + indentStr + content[cursorPos:]
	newCursorPos = cursorPos + len(indentStr)
	
	return newContent, newCursorPos
}

// HandleShiftTabKey handles the Shift+Tab key press (unindent)
func (im *IndentationManager) HandleShiftTabKey(content string, cursorPos int) (newContent string, newCursorPos int) {
	if content == "" {
		return content, cursorPos
	}
	
	lines := strings.Split(content, "\n")
	
	// Find which line the cursor is on
	currentPos := 0
	lineIndex := 0
	for i, line := range lines {
		if currentPos+len(line) >= cursorPos {
			lineIndex = i
			break
		}
		currentPos += len(line) + 1 // +1 for newline
	}
	
	if lineIndex >= len(lines) {
		return content, cursorPos
	}
	
	// Remove indentation from current line
	originalLine := lines[lineIndex]
	newLine, removed := im.removeIndentation(originalLine)
	
	if removed == 0 {
		return content, cursorPos
	}
	
	// Reconstruct content
	lines[lineIndex] = newLine
	newContent = strings.Join(lines, "\n")
	
	// Adjust cursor position
	newCursorPos = cursorPos - removed
	if newCursorPos < 0 {
		newCursorPos = 0
	}
	
	return newContent, newCursorPos
}

// HandleEnterKey handles the Enter key press with auto-indentation
func (im *IndentationManager) HandleEnterKey(content string, cursorPos int) (newContent string, newCursorPos int) {
	if !im.autoIndent {
		// Just insert newline without auto-indent
		newContent = content[:cursorPos] + "\n" + content[cursorPos:]
		return newContent, cursorPos + 1
	}
	
	lines := strings.Split(content, "\n")
	
	// Find current line
	currentPos := 0
	lineIndex := 0
	for i, line := range lines {
		if currentPos+len(line) >= cursorPos {
			lineIndex = i
			break
		}
		currentPos += len(line) + 1
	}
	
	if lineIndex >= len(lines) {
		return content + "\n", len(content) + 1
	}
	
	// Get indentation of current line
	currentLine := lines[lineIndex]
	indentation := im.GetLineIndentation(currentLine)
	
	// Check if we need extra indentation (e.g., after opening braces)
	extraIndent := ""
	trimmedLine := strings.TrimSpace(currentLine)
	if strings.HasSuffix(trimmedLine, "{") || strings.HasSuffix(trimmedLine, ":") {
		extraIndent = im.GetIndentString()
	}
	
	// Insert newline with indentation
	newLineContent := "\n" + indentation + extraIndent
	newContent = content[:cursorPos] + newLineContent + content[cursorPos:]
	newCursorPos = cursorPos + len(newLineContent)
	
	return newContent, newCursorPos
}

// IndentLines indents the selected lines
func (im *IndentationManager) IndentLines(content string, startLine, endLine int) string {
	lines := strings.Split(content, "\n")
	
	if startLine < 0 {
		startLine = 0
	}
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	if startLine > endLine {
		return content
	}
	
	indentStr := im.GetIndentString()
	
	// Add indentation to each line in the range
	for i := startLine; i <= endLine; i++ {
		if i < len(lines) && strings.TrimSpace(lines[i]) != "" {
			lines[i] = indentStr + lines[i]
		}
	}
	
	return strings.Join(lines, "\n")
}

// UnindentLines removes indentation from the selected lines
func (im *IndentationManager) UnindentLines(content string, startLine, endLine int) string {
	lines := strings.Split(content, "\n")
	
	if startLine < 0 {
		startLine = 0
	}
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	if startLine > endLine {
		return content
	}
	
	// Remove indentation from each line in the range
	for i := startLine; i <= endLine; i++ {
		if i < len(lines) {
			lines[i], _ = im.removeIndentation(lines[i])
		}
	}
	
	return strings.Join(lines, "\n")
}

// GetLineIndentation returns the indentation string of a line
func (im *IndentationManager) GetLineIndentation(line string) string {
	indentation := ""
	for _, char := range line {
		if char == ' ' || char == '\t' {
			indentation += string(char)
		} else {
			break
		}
	}
	return indentation
}

// removeIndentation removes one level of indentation from a line
func (im *IndentationManager) removeIndentation(line string) (newLine string, removedChars int) {
	if line == "" {
		return line, 0
	}
	
	indentStr := im.GetIndentString()
	
	// If the line starts with our indent string, remove it
	if strings.HasPrefix(line, indentStr) {
		return line[len(indentStr):], len(indentStr)
	}
	
	// Otherwise, remove leading whitespace up to tab size
	removed := 0
	runes := []rune(line)
	
	for i, r := range runes {
		if r == ' ' {
			removed++
			if removed >= im.tabSize {
				return string(runes[i+1:]), removed
			}
		} else if r == '\t' {
			return string(runes[i+1:]), 1
		} else {
			break
		}
	}
	
	if removed > 0 {
		return string(runes[removed:]), removed
	}
	
	return line, 0
}

// ConvertTabsToSpaces converts all tabs in content to spaces
func (im *IndentationManager) ConvertTabsToSpaces(content string) string {
	spaceStr := strings.Repeat(" ", im.tabSize)
	return strings.ReplaceAll(content, "\t", spaceStr)
}

// ConvertSpacesToTabs converts leading spaces to tabs
func (im *IndentationManager) ConvertSpacesToTabs(content string) string {
	lines := strings.Split(content, "\n")
	spaceStr := strings.Repeat(" ", im.tabSize)
	
	for i, line := range lines {
		// Only convert leading spaces
		indentation := im.GetLineIndentation(line)
		if strings.Contains(indentation, " ") {
			// Convert groups of spaces to tabs
			newIndentation := strings.ReplaceAll(indentation, spaceStr, "\t")
			lines[i] = newIndentation + strings.TrimLeft(line, " \t")
		}
	}
	
	return strings.Join(lines, "\n")
}

// GetIndentationLevel returns the indentation level of a line
func (im *IndentationManager) GetIndentationLevel(line string) int {
	level := 0
	indentStr := im.GetIndentString()
	
	for strings.HasPrefix(line, indentStr) {
		level++
		line = line[len(indentStr):]
	}
	
	return level
}

// IsWhitespaceOnly returns true if the line contains only whitespace
func (im *IndentationManager) IsWhitespaceOnly(line string) bool {
	return strings.TrimSpace(line) == ""
}

// TrimTrailingWhitespace removes trailing whitespace from all lines
func (im *IndentationManager) TrimTrailingWhitespace(content string) string {
	lines := strings.Split(content, "\n")
	
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	
	return strings.Join(lines, "\n")
}