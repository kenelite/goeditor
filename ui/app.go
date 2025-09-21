package ui

import (
	"path/filepath"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
)

func StartApp() {
	a := app.New()
	w := a.NewWindow("Goeditor")

	editor := NewEditor()
	
	// Initialize dialogs after window is created
	editor.InitializeDialogs(w)
	
	menu := NewMenu(w, editor)

	// Set up editor callbacks
	editor.OnFileChanged = func(path string) {
		updateWindowTitle(w, editor)
		w.SetMainMenu(NewMenu(w, editor))
	}
	
	editor.OnModified = func(modified bool) {
		updateWindowTitle(w, editor)
	}

	// Apply configuration
	editor.ApplyConfiguration()
	
	// Set up window from configuration
	config := editor.ConfigManager.GetUIConfig()
	w.Resize(fyne.NewSize(float32(config.WindowWidth), float32(config.WindowHeight)))
	
	w.SetMainMenu(menu)
	w.SetContent(container.NewMax(editor.TextWidget))

	setupShortcuts(w, editor)
	updateWindowTitle(w, editor)
	w.ShowAndRun()
}

// updateWindowTitle updates the window title based on current file and state
func updateWindowTitle(w fyne.Window, editor *Editor) {
	title := "Goeditor"
	
	if editor.GetCurrentFile() != "" {
		filename := filepath.Base(editor.GetCurrentFile())
		title = filename + " - Goeditor"
		
		if editor.IsModified() {
			title = "*" + title
		}
	} else if editor.IsModified() {
		title = "*Untitled - Goeditor"
	}
	
	w.SetTitle(title)
}

func setupShortcuts(w fyne.Window, editor *Editor) {
	// New file
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		editor.NewFile()
	})

	// Open file
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		openDialog := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			if r == nil {
				return
			}
			path := r.URI().Path()
			if err := editor.LoadFile(path); err != nil {
				dialog.ShowError(err, w)
			}
		}, w)
		openDialog.Show()
	})

	// Save file
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		if editor.GetCurrentFile() != "" {
			if err := editor.SaveFile(editor.GetCurrentFile()); err != nil {
				dialog.ShowError(err, w)
			}
		} else {
			saveDialog := dialog.NewFileSave(func(wr fyne.URIWriteCloser, err error) {
				if wr == nil {
					return
				}
				path := wr.URI().Path()
				if err := editor.SaveFile(path); err != nil {
					dialog.ShowError(err, w)
				}
			}, w)
			saveDialog.Show()
		}
	})

	// Save as
	w.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
	}, func(sc fyne.Shortcut) {
		dialog.NewFileSave(func(wr fyne.URIWriteCloser, err error) {
			if wr == nil {
				return
			}
			path := wr.URI().Path()
			if err := editor.SaveFile(path); err != nil {
				dialog.ShowError(err, w)
			}
		}, w).Show()
	})

	// Undo
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyZ, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		editor.Undo()
	})

	// Redo
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyY, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		editor.Redo()
	})

	// Find
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		editor.ShowFindDialog()
	})

	// Replace
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyH, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		editor.ShowReplaceDialog()
	})

	// Find Next
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF3}, func(sc fyne.Shortcut) {
		editor.FindNext()
	})

	// Find Previous
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF3, Modifier: fyne.KeyModifierShift}, func(sc fyne.Shortcut) {
		editor.FindPrevious()
	})

	// Quit
	w.Canvas().AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}, func(sc fyne.Shortcut) {
		w.Close()
	})
}
