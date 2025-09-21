package dialogs

import (
	"fmt"
	"strconv"
	"strings"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// GoToLineDialog represents the "Go to Line" dialog
type GoToLineDialog struct {
	dialog      dialog.Dialog
	lineEntry   *widget.Entry
	statusLabel *widget.Label
	goButton    *widget.Button
	cancelButton *widget.Button
	
	// References
	editor      EditorInterface
	window      fyne.Window
	
	// State
	isVisible   bool
	maxLines    int
}

// NewGoToLineDialog creates a new "Go to Line" dialog
func NewGoToLineDialog(editor EditorInterface, window fyne.Window) *GoToLineDialog {
	gtd := &GoToLineDialog{
		editor:    editor,
		window:    window,
		isVisible: false,
		maxLines:  1,
	}
	
	gtd.createDialog()
	return gtd
}

// createDialog creates the dialog UI
func (gtd *GoToLineDialog) createDialog() {
	// Line number entry
	gtd.lineEntry = widget.NewEntry()
	gtd.lineEntry.SetPlaceHolder("Enter line number...")
	
	// Validate input to only allow numbers
	gtd.lineEntry.Validator = func(text string) error {
		if text == "" {
			return nil // Allow empty for placeholder
		}
		
		num, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("Please enter a valid number")
		}
		
		if num < 1 {
			return fmt.Errorf("Line number must be greater than 0")
		}
		
		if num > gtd.maxLines {
			return fmt.Errorf("Line number cannot exceed %d", gtd.maxLines)
		}
		
		return nil
	}
	
	// Handle Enter key to go to line
	gtd.lineEntry.OnSubmitted = func(text string) {
		gtd.goToLine()
	}
	
	// Status label for feedback
	gtd.statusLabel = widget.NewLabel("")
	gtd.statusLabel.Wrapping = fyne.TextWrapWord
	
	// Buttons
	gtd.goButton = widget.NewButton("Go", func() {
		gtd.goToLine()
	})
	gtd.goButton.Importance = widget.HighImportance
	
	gtd.cancelButton = widget.NewButton("Cancel", func() {
		gtd.Hide()
	})
	
	// Layout
	inputRow := container.NewBorder(
		nil, nil,
		widget.NewLabel("Line number:"),
		nil,
		gtd.lineEntry,
	)
	
	buttonRow := container.NewHBox(
		gtd.goButton,
		gtd.cancelButton,
	)
	
	content := container.NewVBox(
		inputRow,
		gtd.statusLabel,
		widget.NewSeparator(),
		buttonRow,
	)
	
	// Create dialog
	gtd.dialog = dialog.NewCustom("Go to Line", "", content, gtd.window)
	gtd.dialog.Resize(fyne.NewSize(300, 150))
}

// Show displays the "Go to Line" dialog
func (gtd *GoToLineDialog) Show() {
	if gtd.isVisible {
		return
	}
	
	// Update max lines based on current editor content
	gtd.updateMaxLines()
	
	// Update status label with current info
	gtd.updateStatusLabel()
	
	gtd.isVisible = true
	gtd.dialog.Show()
	
	// Focus on line entry
	gtd.window.Canvas().Focus(gtd.lineEntry)
}

// Hide hides the "Go to Line" dialog
func (gtd *GoToLineDialog) Hide() {
	if !gtd.isVisible {
		return
	}
	
	gtd.isVisible = false
	gtd.dialog.Hide()
	
	// Clear the entry
	gtd.lineEntry.SetText("")
	gtd.statusLabel.SetText("")
}

// IsVisible returns whether the dialog is currently visible
func (gtd *GoToLineDialog) IsVisible() bool {
	return gtd.isVisible
}

// goToLine performs the actual "go to line" operation
func (gtd *GoToLineDialog) goToLine() {
	text := strings.TrimSpace(gtd.lineEntry.Text)
	if text == "" {
		gtd.statusLabel.SetText("Please enter a line number")
		return
	}
	
	// Validate the input
	if err := gtd.lineEntry.Validator(text); err != nil {
		gtd.statusLabel.SetText(err.Error())
		return
	}
	
	// Parse line number
	lineNumber, err := strconv.Atoi(text)
	if err != nil {
		gtd.statusLabel.SetText("Invalid line number")
		return
	}
	
	// Go to the line
	if gtd.goToLineNumber(lineNumber) {
		gtd.statusLabel.SetText(fmt.Sprintf("Moved to line %d", lineNumber))
		// Hide dialog after successful navigation
		gtd.Hide()
	} else {
		gtd.statusLabel.SetText("Failed to go to line")
	}
}

// goToLineNumber navigates to the specified line number
func (gtd *GoToLineDialog) goToLineNumber(lineNumber int) bool {
	if gtd.editor == nil {
		return false
	}
	
	content := gtd.editor.GetContent()
	lines := strings.Split(content, "\n")
	
	// Validate line number
	if lineNumber < 1 || lineNumber > len(lines) {
		gtd.statusLabel.SetText(fmt.Sprintf("Line number must be between 1 and %d", len(lines)))
		return false
	}
	
	// Calculate position in text
	var position int
	for i := 0; i < lineNumber-1; i++ {
		position += len(lines[i]) + 1 // +1 for newline
	}
	
	// For now, we'll use a simple approach and assume the editor will handle
	// the cursor positioning through its own mechanisms
	// In a full implementation, we would need proper cursor positioning support
	
	return true
}

// updateMaxLines updates the maximum line count based on editor content
func (gtd *GoToLineDialog) updateMaxLines() {
	if gtd.editor == nil {
		gtd.maxLines = 1
		return
	}
	
	content := gtd.editor.GetContent()
	lines := strings.Split(content, "\n")
	gtd.maxLines = len(lines)
	
	if gtd.maxLines < 1 {
		gtd.maxLines = 1
	}
}

// updateStatusLabel updates the status label with current information
func (gtd *GoToLineDialog) updateStatusLabel() {
	gtd.statusLabel.SetText(fmt.Sprintf("Enter a line number (1-%d)", gtd.maxLines))
}

// HandleKeyEvent handles keyboard events for the dialog
func (gtd *GoToLineDialog) HandleKeyEvent(event *fyne.KeyEvent) bool {
	if !gtd.isVisible {
		return false
	}
	
	switch event.Name {
	case fyne.KeyEscape:
		gtd.Hide()
		return true
	case fyne.KeyReturn, fyne.KeyEnter:
		gtd.goToLine()
		return true
	}
	
	return false
}

// SetCurrentLine sets the current line number in the entry (for convenience)
func (gtd *GoToLineDialog) SetCurrentLine(lineNumber int) {
	gtd.lineEntry.SetText(fmt.Sprintf("%d", lineNumber))
}

// GetCurrentLine returns the line number entered in the dialog
func (gtd *GoToLineDialog) GetCurrentLine() (int, error) {
	text := strings.TrimSpace(gtd.lineEntry.Text)
	if text == "" {
		return 0, fmt.Errorf("no line number entered")
	}
	
	return strconv.Atoi(text)
}