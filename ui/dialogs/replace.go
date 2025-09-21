package dialogs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kenelite/goeditor/backend"
)

// ReplaceDialog represents the find and replace dialog
type ReplaceDialog struct {
	*FindDialog
	
	// Additional UI elements for replace
	replaceEntry     *widget.Entry
	replaceButton    *widget.Button
	replaceAllButton *widget.Button
	
	// State
	lastReplaceText  string
}

// NewReplaceDialog creates a new find and replace dialog
func NewReplaceDialog(editor EditorInterface, searchManager *backend.SearchManager, window fyne.Window) *ReplaceDialog {
	rd := &ReplaceDialog{
		FindDialog: NewFindDialog(editor, searchManager, window),
	}
	
	// Override the dialog creation to include replace functionality
	rd.createReplaceDialog()
	return rd
}

// createReplaceDialog creates the replace dialog UI
func (rd *ReplaceDialog) createReplaceDialog() {
	// Reuse find dialog components but create new dialog
	
	// Search entry (reuse from FindDialog)
	rd.searchEntry.OnChanged = func(text string) {
		rd.performSearch(text)
	}
	rd.searchEntry.OnSubmitted = func(text string) {
		rd.FindNext()
	}

	// Replace entry
	rd.replaceEntry = widget.NewEntry()
	rd.replaceEntry.SetPlaceHolder("Enter replacement text...")
	rd.replaceEntry.OnChanged = func(text string) {
		rd.lastReplaceText = text
		rd.updateReplaceButtons()
	}

	// Replace buttons
	rd.replaceButton = widget.NewButton("Replace", func() {
		rd.ReplaceCurrent()
	})
	rd.replaceButton.Disable()

	rd.replaceAllButton = widget.NewButton("Replace All", func() {
		rd.ReplaceAll()
	})
	rd.replaceAllButton.Disable()

	// Update existing buttons
	rd.closeButton = widget.NewButton("Close", func() {
		rd.Hide()
	})

	// Layout
	searchRow := container.NewBorder(nil, nil, widget.NewLabel("Find:"), nil, rd.searchEntry)
	replaceRow := container.NewBorder(nil, nil, widget.NewLabel("Replace:"), nil, rd.replaceEntry)
	
	optionsRow1 := container.NewHBox(
		rd.optionsCheck["caseSensitive"],
		rd.optionsCheck["wholeWord"],
	)
	
	optionsRow2 := container.NewHBox(
		rd.optionsCheck["regex"],
		rd.optionsCheck["wrapAround"],
	)
	
	navigationRow := container.NewHBox(
		rd.prevButton,
		rd.nextButton,
	)
	
	replaceButtonRow := container.NewHBox(
		rd.replaceButton,
		rd.replaceAllButton,
	)
	
	closeRow := container.NewHBox(
		widget.NewSeparator(),
		rd.closeButton,
	)

	content := container.NewVBox(
		searchRow,
		replaceRow,
		widget.NewSeparator(),
		widget.NewLabel("Options:"),
		optionsRow1,
		optionsRow2,
		widget.NewSeparator(),
		rd.resultLabel,
		widget.NewSeparator(),
		navigationRow,
		replaceButtonRow,
		closeRow,
	)

	// Create new dialog
	rd.dialog = dialog.NewCustom("Find and Replace", "", content, rd.window)
	rd.dialog.Resize(fyne.NewSize(450, 400))
}

// ReplaceCurrent replaces the current match
func (rd *ReplaceDialog) ReplaceCurrent() {
	if !rd.searchManager.HasMatches() {
		return
	}

	searchText := rd.searchEntry.Text
	replaceText := rd.replaceEntry.Text
	
	if searchText == "" {
		return
	}

	// Get current editor content
	content := rd.editor.GetContent()
	
	// Perform replacement
	options := backend.ReplaceOptions{
		SearchOptions: rd.searchManager.GetOptions(),
		ReplaceAll:    false,
	}
	
	newContent, count := rd.searchManager.Replace(content, searchText, replaceText, options)
	
	if count > 0 {
		// Update editor content
		rd.editor.SetContent(newContent)
		
		// Update result label
		rd.resultLabel.SetText(fmt.Sprintf("Replaced 1 occurrence"))
		
		// Re-search to update matches
		rd.performSearch(searchText)
		
		// Move to next match if available
		if rd.searchManager.HasMatches() {
			rd.FindNext()
		}
	}
}

// ReplaceAll replaces all matches
func (rd *ReplaceDialog) ReplaceAll() {
	searchText := rd.searchEntry.Text
	replaceText := rd.replaceEntry.Text
	
	if searchText == "" {
		return
	}

	// Get current editor content
	content := rd.editor.GetContent()
	
	// Perform replacement
	options := backend.ReplaceOptions{
		SearchOptions: rd.searchManager.GetOptions(),
		ReplaceAll:    true,
	}
	
	newContent, count := rd.searchManager.Replace(content, searchText, replaceText, options)
	
	if count > 0 {
		// Update editor content
		rd.editor.SetContent(newContent)
		
		// Update result label
		rd.resultLabel.SetText(fmt.Sprintf("Replaced %d occurrences", count))
		
		// Clear search since all matches are replaced
		rd.searchManager.Clear()
		rd.updateButtons()
		rd.updateReplaceButtons()
	} else {
		rd.resultLabel.SetText("No matches found to replace")
	}
}

// SetReplaceText sets the replace text programmatically
func (rd *ReplaceDialog) SetReplaceText(text string) {
	rd.replaceEntry.SetText(text)
	rd.lastReplaceText = text
	rd.updateReplaceButtons()
}

// GetReplaceText returns the current replace text
func (rd *ReplaceDialog) GetReplaceText() string {
	return rd.replaceEntry.Text
}

// Show displays the replace dialog
func (rd *ReplaceDialog) Show() {
	if rd.isVisible {
		return
	}
	
	rd.isVisible = true
	rd.dialog.Show()
	
	// Focus on search entry
	rd.window.Canvas().Focus(rd.searchEntry)
	
	// Update button states
	rd.updateReplaceButtons()
}

// updateReplaceButtons updates replace button states
func (rd *ReplaceDialog) updateReplaceButtons() {
	hasMatches := rd.searchManager.HasMatches()
	hasSearchText := rd.searchEntry.Text != ""
	
	// Replace button is enabled when we have a current match
	if hasMatches && rd.searchManager.GetCurrentMatch() != nil {
		rd.replaceButton.Enable()
	} else {
		rd.replaceButton.Disable()
	}
	
	// Replace All button is enabled when we have search text and matches
	if hasSearchText && hasMatches {
		rd.replaceAllButton.Enable()
	} else {
		rd.replaceAllButton.Disable()
	}
}

// Override performSearch to update replace buttons
func (rd *ReplaceDialog) performSearch(pattern string) {
	// Call parent method
	rd.FindDialog.performSearch(pattern)
	
	// Update replace buttons
	rd.updateReplaceButtons()
}

// Override updateButtons to also update replace buttons
func (rd *ReplaceDialog) updateButtons() {
	// Call parent method
	rd.FindDialog.updateButtons()
	
	// Update replace buttons
	rd.updateReplaceButtons()
}

// HandleKeyEvent handles keyboard events for the replace dialog
func (rd *ReplaceDialog) HandleKeyEvent(event *fyne.KeyEvent) bool {
	if !rd.isVisible {
		return false
	}

	switch event.Name {
	case fyne.KeyEscape:
		rd.Hide()
		return true
	case fyne.KeyReturn, fyne.KeyEnter:
		// If focus is on replace entry, perform replace
		if rd.window.Canvas().Focused() == rd.replaceEntry {
			rd.ReplaceCurrent()
			return true
		}
		// Otherwise, handle as find
		rd.FindNext()
		return true
	case fyne.KeyF3:
		rd.FindNext()
		return true
	}

	return false
}

// ConfirmReplaceAll shows a confirmation dialog for replace all operation
func (rd *ReplaceDialog) ConfirmReplaceAll() {
	matchCount := rd.searchManager.GetMatchCount()
	if matchCount == 0 {
		return
	}

	message := fmt.Sprintf("Replace all %d occurrences of '%s' with '%s'?", 
		matchCount, rd.searchEntry.Text, rd.replaceEntry.Text)
	
	confirm := dialog.NewConfirm("Confirm Replace All", message, func(confirmed bool) {
		if confirmed {
			rd.ReplaceAll()
		}
	}, rd.window)
	
	confirm.Show()
}

// ShowReplaceAllConfirmation enables/disables replace all confirmation
func (rd *ReplaceDialog) ShowReplaceAllConfirmation(show bool) {
	if show {
		// Override replace all button to show confirmation
		rd.replaceAllButton.OnTapped = func() {
			rd.ConfirmReplaceAll()
		}
	} else {
		// Direct replace all without confirmation
		rd.replaceAllButton.OnTapped = func() {
			rd.ReplaceAll()
		}
	}
}