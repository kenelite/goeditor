package ui

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/kenelite/goeditor/backend"
	"os"
)

func NewMenu(win fyne.Window, editor *Editor) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", func() {
		editor.TextWidget.SetText("")
		editor.State.CurrentFile = ""
	})

	openItem := fyne.NewMenuItem("Open", func() {
		dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
			if err != nil || r == nil {
				return
			}
			content := backend.ReadFile(r.URI().Path())
			editor.TextWidget.SetText(content)
			editor.State.CurrentFile = r.URI().Path()
		}, win)
	})

	saveItem := fyne.NewMenuItem("Save", func() {
		if editor.State.CurrentFile != "" {
			backend.SaveFile(editor.State.CurrentFile, editor.TextWidget.Text)
		} else {
			dialog.ShowFileSave(func(wr fyne.URIWriteCloser, err error) {
				if err != nil || wr == nil {
					return
				}
				backend.SaveFile(wr.URI().Path(), editor.TextWidget.Text)
				editor.State.CurrentFile = wr.URI().Path()
			}, win)
		}
	})

	saveAsItem := fyne.NewMenuItem("Save As", func() {
		dialog.ShowFileSave(func(wr fyne.URIWriteCloser, err error) {
			if err != nil || wr == nil {
				return
			}
			backend.SaveFile(wr.URI().Path(), editor.TextWidget.Text)
			editor.State.CurrentFile = wr.URI().Path()
		}, win)
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		os.Exit(0)
	})

	fileMenu := fyne.NewMenu("File", newItem, openItem, saveItem, saveAsItem, quitItem)
	return fyne.NewMainMenu(fileMenu)
}
