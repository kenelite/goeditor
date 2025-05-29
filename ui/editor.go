package ui

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kenelite/goeditor/backend"
	"github.com/kenelite/goeditor/ui/syntax"
	"os"
)

type Editor struct {
	TextWidget *widget.Entry
	State      *backend.State
}

func NewEditor() *Editor {
	e := &Editor{
		TextWidget: widget.NewMultiLineEntry(),
		State:      backend.NewState(),
	}
	return e
}

func EditorView(path string) fyne.CanvasObject {
	code, err := os.ReadFile(path)
	if err != nil {
		return widget.NewLabel("Failed to open file: " + err.Error())
	}

	highlighted := syntax.HighlightGoCode(string(code))
	richText := widget.NewRichText(highlighted...)
	scroll := container.NewScroll(richText)

	return scroll
}
