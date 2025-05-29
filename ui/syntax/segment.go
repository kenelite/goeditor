package syntax

import (
	"image/color"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

const averageCharWidth = 8

type SyntaxSegment struct {
	Text        string
	Color       color.Color
	selected    bool
	selectedRun string
}

func (s *SyntaxSegment) Copy() widget.RichTextSegment {
	return &SyntaxSegment{
		Text:  s.Text,
		Color: s.Color,
	}
}

func (s *SyntaxSegment) Inline() bool {
	return true
}

func (s *SyntaxSegment) Visual() fyne.CanvasObject {
	col := s.Color
	if s.selected {
		col = color.White
	}
	t := canvas.NewText(s.Text, col)
	t.TextSize = 14
	return t
}

func (s *SyntaxSegment) Select(start, end fyne.Position) {
	s.selected = true

	// Simply convert the X pixel coordinate to a character index (approximately)
	startIdx := int(start.X / averageCharWidth)
	endIdx := int(end.X / averageCharWidth)

	// Border protection
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(s.Text) {
		endIdx = len(s.Text)
	}
	if startIdx > endIdx {
		startIdx, endIdx = endIdx, startIdx
	}

	s.selectedRun = s.Text[startIdx:endIdx]
}

func (s *SyntaxSegment) Unselect() {
	s.selected = false
	s.selectedRun = ""
}

func (s *SyntaxSegment) SelectedText() string {
	return s.selectedRun
}

func (s *SyntaxSegment) Textual() string {
	return s.Text
}

func (s *SyntaxSegment) Update(obj fyne.CanvasObject) {

}

func NewSyntaxSegment(text string, col color.Color) *SyntaxSegment {
	return &SyntaxSegment{
		Text:  text,
		Color: col,
	}
}
