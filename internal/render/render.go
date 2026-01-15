package render

import (
	"fmt"
	"io"
	"strings"

	"texteditor/internal/editor"
	"texteditor/internal/keys"
)

func Draw(w io.Writer, es *editor.EditorState, last keys.KeyEvent) {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("\x1b[2J")
	stringBuilder.WriteString("\x1b[H")

	screenW := es.Width()
	screenH := es.Height()
	contentW := screenW
	contentH := screenH - 1

	for y := 0; y < contentH; y++ {
		for x := 0; x < contentW; x++ {
			if char, ok := es.Cell(x, y); ok {
				stringBuilder.WriteByte(char)
			} else if x == 0 {
				stringBuilder.WriteByte('~')
			} else {
				stringBuilder.WriteByte(' ')
			}
		}

		stringBuilder.WriteString("\r\n")
	}

	fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", keys.KeyLabel(last), screenW, screenH)

	cursorX, cursorY := es.Cursor()
	fmt.Fprintf(&stringBuilder, "\x1b[%d;%dH", cursorY+1, cursorX+1) // ANSI cursor positions are 1-based
	io.WriteString(w, stringBuilder.String())
}
