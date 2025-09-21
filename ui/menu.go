package ui

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"os"
)

func NewMenu(win fyne.Window, editor *Editor) *fyne.MainMenu {
	// File menu items
	newItem := fyne.NewMenuItem("New", func() {
		editor.NewFile()
	})
	// Shortcuts are handled by the setupShortcuts function

	openItem := fyne.NewMenuItem("Open", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil || r == nil {
				return
			}
			if err := editor.LoadFile(r.URI().Path()); err != nil {
				dialog.ShowError(err, win)
			}
		}, win)
	})
	// Shortcuts are handled by the setupShortcuts function

	saveItem := fyne.NewMenuItem("Save", func() {
		if editor.GetCurrentFile() != "" {
			if err := editor.SaveFile(editor.GetCurrentFile()); err != nil {
				dialog.ShowError(err, win)
			}
		} else {
			dialog.ShowFileSave(func(wr fyne.URIWriteCloser, err error) {
				if err != nil || wr == nil {
					return
				}
				if err := editor.SaveFile(wr.URI().Path()); err != nil {
					dialog.ShowError(err, win)
				}
			}, win)
		}
	})
	// Shortcuts are handled by the setupShortcuts function

	saveAsItem := fyne.NewMenuItem("Save As", func() {
		dialog.ShowFileSave(func(wr fyne.URIWriteCloser, err error) {
			if err != nil || wr == nil {
				return
			}
			if err := editor.SaveFile(wr.URI().Path()); err != nil {
				dialog.ShowError(err, win)
			}
		}, win)
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		// TODO: Check for unsaved changes before quitting
		os.Exit(0)
	})

	// Edit menu items
	undoItem := fyne.NewMenuItem("Undo", func() {
		editor.Undo()
	})
	// Shortcuts are handled by the setupShortcuts function

	redoItem := fyne.NewMenuItem("Redo", func() {
		editor.Redo()
	})
	// Shortcuts are handled by the setupShortcuts function

	// Search menu items
	findItem := fyne.NewMenuItem("Find", func() {
		editor.ShowFindDialog()
	})
	// Shortcuts are handled by the setupShortcuts function

	replaceItem := fyne.NewMenuItem("Replace", func() {
		editor.ShowReplaceDialog()
	})
	// Shortcuts are handled by the setupShortcuts function

	findNextItem := fyne.NewMenuItem("Find Next", func() {
		editor.FindNext()
	})
	// Shortcuts are handled by the setupShortcuts function

	findPrevItem := fyne.NewMenuItem("Find Previous", func() {
		editor.FindPrevious()
	})
	// Shortcuts are handled by the setupShortcuts function

	// Navigation menu items
	goToLineItem := fyne.NewMenuItem("Go to Line...", func() {
		editor.ShowGoToLineDialog()
	})
	// Shortcuts are handled by the setupShortcuts function

	// Format menu items
	indentItem := fyne.NewMenuItem("Indent Lines", func() {
		editor.IndentSelectedLines()
	})
	// Shortcuts are handled by the setupShortcuts function

	unindentItem := fyne.NewMenuItem("Unindent Lines", func() {
		editor.UnindentSelectedLines()
	})
	// Shortcuts are handled by the setupShortcuts function

	// Enable/disable menu items based on state
	saveItem.Disabled = !editor.IsModified()
	undoItem.Disabled = !editor.CanUndo()
	redoItem.Disabled = !editor.CanRedo()

	// Create menus - simplified to avoid crashes
	fileMenu := fyne.NewMenu("File", newItem, openItem, saveItem, saveAsItem, quitItem)
	editMenu := fyne.NewMenu("Edit", undoItem, redoItem, findItem, replaceItem, findNextItem, findPrevItem, goToLineItem)
	formatMenu := fyne.NewMenu("Format", indentItem, unindentItem)
	
	return fyne.NewMainMenu(fileMenu, editMenu, formatMenu)
}
