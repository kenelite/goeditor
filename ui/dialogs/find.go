package dialogs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kenelite/goeditor/backend"
)

// FindDialog represents the find dialog
type FindDialog struct {
	dialog        dialog.Dialog
	searchEntry   *widget.Entry
	optionsCheck  map[string]*widget.Check
	resultLabel   *widget.Label
	nextButton    *widget.Button
	prevButton    *widget.Button
	closeButton   *widget.Button
	
	// References
	editor        EditorInterface
	searchManager *backend.SearchManager
	window        fyne.Window
	
	// State
	isVisible     bool
	lastPattern   string
}

// EditorInterface defines the interface that the editor must implement
type EditorInterface interface {
	GetContent() string
	SetContent(content string)
	// TODO: Add methods for cursor positioning and text selection when available
	// SetCursorPosition(line, col int)
	// SelectText(start, end backend.Position)
}

// NewFindDialog creates a new find dialog
func NewFindDialog(editor EditorInterface, searchManager *backend.SearchManager, window fyne.Window) *FindDialog {
	fd := &FindDialog{
		editor:        editor,
		searchManager: searchManager,
		window:        window,
		isVisible:     false,
	}
	
	fd.createDialog()
	return fd
}

// createDialog creates the dialog UI
func (fd *FindDialog) createDialog() {
	// Search entry
	fd.searchEntry = widget.NewEntry()
	fd.searchEntry.SetPlaceHolder("Enter search text...")
	fd.searchEntry.OnChanged = func(text string) {
		fd.performSearch(text)
	}
	fd.searchEntry.OnSubmitted = func(text string) {
		fd.FindNext()
	}

	// Options checkboxes
	fd.optionsCheck = make(map[string]*widget.Check)
	
	fd.optionsCheck["caseSensitive"] = widget.NewCheck("Case sensitive", func(checked bool) {
		fd.updateSearchOptions()
	})
	
	fd.optionsCheck["wholeWord"] = widget.NewCheck("Whole word", func(checked bool) {
		fd.updateSearchOptions()
	})
	
	fd.optionsCheck["regex"] = widget.NewCheck("Regular expression", func(checked bool) {
		fd.updateSearchOptions()
	})
	
	fd.optionsCheck["wrapAround"] = widget.NewCheck("Wrap around", func(checked bool) {
		fd.updateSearchOptions()
	})
	
	// Set default values
	options := fd.searchManager.GetOptions()
	fd.optionsCheck["caseSensitive"].SetChecked(options.CaseSensitive)
	fd.optionsCheck["wholeWord"].SetChecked(options.WholeWord)
	fd.optionsCheck["regex"].SetChecked(options.RegularExpression)
	fd.optionsCheck["wrapAround"].SetChecked(options.WrapAround)

	// Result label
	fd.resultLabel = widget.NewLabel("Enter text to search")
	fd.resultLabel.Wrapping = fyne.TextWrapWord

	// Buttons
	fd.nextButton = widget.NewButton("Next", func() {
		fd.FindNext()
	})
	fd.nextButton.Disable()

	fd.prevButton = widget.NewButton("Previous", func() {
		fd.FindPrevious()
	})
	fd.prevButton.Disable()

	fd.closeButton = widget.NewButton("Close", func() {
		fd.Hide()
	})

	// Layout
	searchRow := container.NewBorder(nil, nil, widget.NewLabel("Find:"), nil, fd.searchEntry)
	
	optionsRow1 := container.NewHBox(
		fd.optionsCheck["caseSensitive"],
		fd.optionsCheck["wholeWord"],
	)
	
	optionsRow2 := container.NewHBox(
		fd.optionsCheck["regex"],
		fd.optionsCheck["wrapAround"],
	)
	
	buttonRow := container.NewHBox(
		fd.prevButton,
		fd.nextButton,
		widget.NewSeparator(),
		fd.closeButton,
	)

	content := container.NewVBox(
		searchRow,
		widget.NewSeparator(),
		widget.NewLabel("Options:"),
		optionsRow1,
		optionsRow2,
		widget.NewSeparator(),
		fd.resultLabel,
		widget.NewSeparator(),
		buttonRow,
	)

	// Create dialog
	fd.dialog = dialog.NewCustom("Find", "", content, fd.window)
	fd.dialog.Resize(fyne.NewSize(400, 300))
}

// Show displays the find dialog
func (fd *FindDialog) Show() {
	if fd.isVisible {
		return
	}
	
	fd.isVisible = true
	fd.dialog.Show()
	
	// Focus on search entry
	fd.window.Canvas().Focus(fd.searchEntry)
	
	// If there's selected text in editor, use it as search term
	// TODO: Implement when editor supports text selection
}

// Hide hides the find dialog
func (fd *FindDialog) Hide() {
	if !fd.isVisible {
		return
	}
	
	fd.isVisible = false
	fd.dialog.Hide()
	
	// Clear highlights
	fd.clearHighlights()
}

// IsVisible returns whether the dialog is currently visible
func (fd *FindDialog) IsVisible() bool {
	return fd.isVisible
}

// FindNext finds the next occurrence
func (fd *FindDialog) FindNext() bool {
	if !fd.searchManager.HasMatches() {
		return false
	}

	match := fd.searchManager.NextMatch()
	if match != nil {
		fd.updateResultLabel()
		fd.highlightCurrentMatch()
		// TODO: Scroll to match position when editor supports it
		return true
	}
	
	return false
}

// FindPrevious finds the previous occurrence
func (fd *FindDialog) FindPrevious() bool {
	if !fd.searchManager.HasMatches() {
		return false
	}

	match := fd.searchManager.PreviousMatch()
	if match != nil {
		fd.updateResultLabel()
		fd.highlightCurrentMatch()
		// TODO: Scroll to match position when editor supports it
		return true
	}
	
	return false
}

// SetSearchText sets the search text programmatically
func (fd *FindDialog) SetSearchText(text string) {
	fd.searchEntry.SetText(text)
	fd.performSearch(text)
}

// GetSearchText returns the current search text
func (fd *FindDialog) GetSearchText() string {
	return fd.searchEntry.Text
}

// performSearch performs the actual search
func (fd *FindDialog) performSearch(pattern string) {
	if pattern == "" {
		fd.clearSearch()
		return
	}

	// Only search if pattern changed
	if pattern == fd.lastPattern {
		return
	}
	
	fd.lastPattern = pattern
	
	// Get editor content and perform search
	content := fd.editor.GetContent()
	matches := fd.searchManager.Find(content, pattern)
	
	// Update UI
	fd.updateResultLabel()
	fd.updateButtons()
	
	// Highlight matches
	fd.highlightMatches()
	
	// If we have matches, go to first one
	if len(matches) > 0 {
		fd.searchManager.SetCurrentMatch(0)
		fd.highlightCurrentMatch()
	}
}

// updateSearchOptions updates search options based on checkboxes
func (fd *FindDialog) updateSearchOptions() {
	options := backend.SearchOptions{
		CaseSensitive:     fd.optionsCheck["caseSensitive"].Checked,
		WholeWord:         fd.optionsCheck["wholeWord"].Checked,
		RegularExpression: fd.optionsCheck["regex"].Checked,
		WrapAround:        fd.optionsCheck["wrapAround"].Checked,
	}
	
	fd.searchManager.SetOptions(options)
	
	// Re-search with new options if we have a pattern
	if fd.searchEntry.Text != "" {
		fd.lastPattern = "" // Force re-search
		fd.performSearch(fd.searchEntry.Text)
	}
}

// updateResultLabel updates the result label with match information
func (fd *FindDialog) updateResultLabel() {
	matchCount := fd.searchManager.GetMatchCount()
	currentIndex := fd.searchManager.GetCurrentIndex()
	
	if matchCount == 0 {
		if fd.searchEntry.Text == "" {
			fd.resultLabel.SetText("Enter text to search")
		} else {
			fd.resultLabel.SetText("No matches found")
		}
	} else {
		fd.resultLabel.SetText(fmt.Sprintf("Match %d of %d", currentIndex+1, matchCount))
	}
}

// updateButtons updates button states based on search results
func (fd *FindDialog) updateButtons() {
	hasMatches := fd.searchManager.HasMatches()
	
	if hasMatches {
		fd.nextButton.Enable()
		fd.prevButton.Enable()
	} else {
		fd.nextButton.Disable()
		fd.prevButton.Disable()
	}
}

// clearSearch clears the search state
func (fd *FindDialog) clearSearch() {
	fd.searchManager.Clear()
	fd.lastPattern = ""
	fd.updateResultLabel()
	fd.updateButtons()
	fd.clearHighlights()
}

// highlightMatches highlights all matches in the editor
func (fd *FindDialog) highlightMatches() {
	// TODO: Implement text highlighting when editor supports it
	// This would involve creating a rich text representation with highlighted segments
}

// highlightCurrentMatch highlights the current match
func (fd *FindDialog) highlightCurrentMatch() {
	// TODO: Implement current match highlighting when editor supports it
	// This would involve selecting the text or using a different highlight color
}

// clearHighlights clears all highlights from the editor
func (fd *FindDialog) clearHighlights() {
	// TODO: Implement highlight clearing when editor supports it
}

// HandleKeyEvent handles keyboard events for the find dialog
func (fd *FindDialog) HandleKeyEvent(event *fyne.KeyEvent) bool {
	if !fd.isVisible {
		return false
	}

	switch event.Name {
	case fyne.KeyEscape:
		fd.Hide()
		return true
	case fyne.KeyReturn, fyne.KeyEnter:
		fd.FindNext()
		return true
	case fyne.KeyF3:
		fd.FindNext()
		return true
	}

	return false
}