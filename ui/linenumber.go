package ui

import (
	"fmt"
	"strings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LineNumberWidget displays line numbers for the editor
type LineNumberWidget struct {
	widget.BaseWidget
	editor     *Editor
	lineCount  int
	lineHeight float32
	fontHeight float32
	padding    float32
}

// NewLineNumberWidget creates a new line number widget
func NewLineNumberWidget(editor *Editor) *LineNumberWidget {
	ln := &LineNumberWidget{
		editor:     editor,
		lineCount:  1,
		lineHeight: 20.0, // Default line height
		fontHeight: 14.0, // Default font height
		padding:    5.0,  // Padding between numbers and editor
	}
	
	// Only extend base widget if we're in a UI context
	// This prevents crashes during testing or when UI isn't ready
	if ln != nil {
		ln.ExtendBaseWidget(ln)
	}
	return ln
}

// CreateRenderer creates the renderer for the line number widget
func (ln *LineNumberWidget) CreateRenderer() fyne.WidgetRenderer {
	return &lineNumberRenderer{
		widget: ln,
		texts:  []*canvas.Text{},
	}
}

// UpdateLineCount updates the number of lines to display
func (ln *LineNumberWidget) UpdateLineCount(count int) {
	if count < 1 {
		count = 1
	}
	
	if ln.lineCount != count {
		ln.lineCount = count
		ln.Refresh()
	}
}

// SetLineHeight sets the height of each line
func (ln *LineNumberWidget) SetLineHeight(height float32) {
	if height > 0 && ln.lineHeight != height {
		ln.lineHeight = height
		ln.Refresh()
	}
}

// SetFontHeight sets the font height for line numbers
func (ln *LineNumberWidget) SetFontHeight(height float32) {
	if height > 0 && ln.fontHeight != height {
		ln.fontHeight = height
		ln.Refresh()
	}
}

// GetLineCount returns the current line count
func (ln *LineNumberWidget) GetLineCount() int {
	return ln.lineCount
}

// Tapped handles tap events on line numbers
func (ln *LineNumberWidget) Tapped(ev *fyne.PointEvent) {
	if ln.editor == nil {
		return
	}
	
	// Calculate which line was clicked
	lineIndex := int(ev.Position.Y / ln.lineHeight)
	if lineIndex >= 0 && lineIndex < ln.lineCount {
		lineNumber := lineIndex + 1
		ln.selectLine(lineNumber)
	}
}

// selectLine selects the entire line at the given line number
func (ln *LineNumberWidget) selectLine(lineNumber int) {
	if ln.editor == nil || ln.editor.TextWidget == nil {
		return
	}
	
	content := ln.editor.GetContent()
	lines := strings.Split(content, "\n")
	
	if lineNumber < 1 || lineNumber > len(lines) {
		return
	}
	
	// Calculate start position for the line
	var startPos int
	
	// Calculate start position (sum of all previous lines + newlines)
	for i := 0; i < lineNumber-1; i++ {
		startPos += len(lines[i]) + 1 // +1 for newline character
	}
	
	// Select the text (this is a simplified implementation)
	// In a real implementation, we would need proper text selection support
	ln.editor.TextWidget.CursorRow = lineNumber - 1
	ln.editor.TextWidget.CursorColumn = 0
	
	// Update editor state
	ln.editor.State.SetCursorPosition(lineNumber, 1)
	
	// Notify cursor change callback
	if ln.editor.OnCursorChanged != nil {
		ln.editor.OnCursorChanged(lineNumber, 1)
	}
}

// GetPreferredWidth calculates the preferred width for the line number widget
func (ln *LineNumberWidget) GetPreferredWidth() float32 {
	// Calculate width based on the number of digits needed for the highest line number
	digits := len(fmt.Sprintf("%d", ln.lineCount))
	charWidth := ln.fontHeight * 0.6 // Approximate character width
	return float32(digits)*charWidth + ln.padding*2
}

// lineNumberRenderer renders the line number widget
type lineNumberRenderer struct {
	widget *LineNumberWidget
	texts  []*canvas.Text
}

// Layout arranges the line number texts
func (r *lineNumberRenderer) Layout(size fyne.Size) {
	// Update text objects for each line
	r.updateTexts()
	
	// Position each text object
	for i, text := range r.texts {
		y := float32(i) * r.widget.lineHeight
		text.Move(fyne.NewPos(r.widget.padding, y))
		text.Resize(fyne.NewSize(size.Width-r.widget.padding*2, r.widget.lineHeight))
	}
}

// MinSize returns the minimum size for the widget
func (r *lineNumberRenderer) MinSize() fyne.Size {
	width := r.widget.GetPreferredWidth()
	height := float32(r.widget.lineCount) * r.widget.lineHeight
	return fyne.NewSize(width, height)
}

// Refresh updates the renderer
func (r *lineNumberRenderer) Refresh() {
	r.updateTexts()
	for _, text := range r.texts {
		text.Refresh()
	}
}

// Objects returns all canvas objects
func (r *lineNumberRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(r.texts))
	for i, text := range r.texts {
		objects[i] = text
	}
	return objects
}

// Destroy cleans up the renderer
func (r *lineNumberRenderer) Destroy() {
	r.texts = nil
}

// updateTexts creates or updates text objects for line numbers
func (r *lineNumberRenderer) updateTexts() {
	// Clear existing texts if line count changed
	if len(r.texts) != r.widget.lineCount {
		r.texts = make([]*canvas.Text, r.widget.lineCount)
	}
	
	// Create or update text objects
	for i := 0; i < r.widget.lineCount; i++ {
		lineNumber := i + 1
		
		if r.texts[i] == nil {
			r.texts[i] = canvas.NewText(fmt.Sprintf("%d", lineNumber), theme.ForegroundColor())
			r.texts[i].TextSize = r.widget.fontHeight
			r.texts[i].Alignment = fyne.TextAlignTrailing // Right-align numbers
		} else {
			r.texts[i].Text = fmt.Sprintf("%d", lineNumber)
			r.texts[i].Color = theme.ForegroundColor()
			r.texts[i].TextSize = r.widget.fontHeight
		}
	}
}

// ScrollToLine scrolls the line number widget to show a specific line
func (ln *LineNumberWidget) ScrollToLine(line int) {
	// This would be implemented when integrated with a scroll container
	// For now, we just ensure the line is within bounds
	if line < 1 {
		line = 1
	}
	if line > ln.lineCount {
		line = ln.lineCount
	}
	
	// Update editor state
	if ln.editor != nil {
		ln.editor.State.SetCursorPosition(line, 1)
		if ln.editor.OnCursorChanged != nil {
			ln.editor.OnCursorChanged(line, 1)
		}
	}
}

// UpdateFromEditor updates line count based on editor content
func (ln *LineNumberWidget) UpdateFromEditor() {
	if ln.editor == nil {
		return
	}
	
	content := ln.editor.GetContent()
	lines := strings.Split(content, "\n")
	ln.UpdateLineCount(len(lines))
}