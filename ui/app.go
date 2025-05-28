package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/kenelite/goeditor/backend"
)

func StartApp() {
	a := app.New()
	w := a.NewWindow("Goeditor")

	editor := NewEditor()
	menu := NewMenu(w, editor)

	w.SetMainMenu(menu)
	w.SetContent(container.NewMax(editor.TextWidget))
	w.Resize(fyne.NewSize(800, 600))

	setupShortcuts(w, editor)
	w.ShowAndRun()
}

func setupShortcuts(w fyne.Window, editor *Editor) {
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		editor.TextWidget.SetText("")
		editor.State.CurrentFile = ""
	})

	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		openDialog := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r == nil {
				return
			}
			path := r.URI().Path()
			content := backend.ReadFile(path)
			editor.TextWidget.SetText(content)
			editor.State.CurrentFile = path
			w.SetMainMenu(NewMenu(w, editor))
		}, w)
		openDialog.Show()
	})

	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		if editor.State.CurrentFile != "" {
			backend.SaveFile(editor.State.CurrentFile, editor.TextWidget.Text)
		} else {
			saveDialog := dialog.NewFileSave(func(wr fyne.URIWriteCloser, err error) {
				if wr == nil {
					return
				}
				path := wr.URI().Path()
				backend.SaveFile(path, editor.TextWidget.Text)
				editor.State.CurrentFile = path
				w.SetMainMenu(NewMenu(w, editor))
			}, w)
			saveDialog.Show()
		}
	})

	w.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
	}, func(sc fyne.Shortcut) {
		dialog.NewFileSave(func(wr fyne.URIWriteCloser, err error) {
			if wr == nil {
				return
			}
			path := wr.URI().Path()
			backend.SaveFile(path, editor.TextWidget.Text)
			editor.State.CurrentFile = path
			w.SetMainMenu(NewMenu(w, editor))
		}, w).Show()
	})

	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		w.Close()
	})
}
