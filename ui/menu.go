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

	// Enable/disable menu items based on state
	saveItem.Disabled = !editor.IsModified()

	fileMenu := fyne.NewMenu("File", newItem, openItem, saveItem, saveAsItem, quitItem)
	return fyne.NewMainMenu(fileMenu)
}
