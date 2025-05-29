package syntax

import (
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"image/color"
)

// chromaToRGBA converts chroma.Colour to color.Color
func chromaToRGBA(c chroma.Colour) color.Color {
	r := uint8((c >> 16) & 0xFF)
	g := uint8((c >> 8) & 0xFF)
	b := uint8(c & 0xFF)
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// HighlightGoCode uses Chroma to highlight Go code and returns a RichTextSegment slice
func HighlightGoCode(source string) []widget.RichTextSegment {
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		// 出错时返回纯文本段
		return []widget.RichTextSegment{
			NewSyntaxSegment(source, color.Black),
		}
	}

	var segments []widget.RichTextSegment
	for token := iterator(); token != chroma.EOF; token = iterator() {
		sty := style.Get(token.Type)
		col := chromaToRGBA(sty.Colour)
		seg := NewSyntaxSegment(token.Value, col)
		segments = append(segments, seg)
	}

	return segments
}
