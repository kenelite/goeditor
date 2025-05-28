package ui

import (
	"fyne.io/fyne/v2/widget"
	"github.com/kenelite/goeditor/backend"
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
