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
	contentH :=  screenH - 1
	rowOffset := es.RowOffset()

	for y := 0; y < contentH; y++ {
		docY := rowOffset + y
		if docY < es.LineCount() {
			currentLine := es.Line(docY)

			currentX := 0
			for currentX < contentW {
				if currentX < len(currentLine) {
					stringBuilder.WriteRune(currentLine[currentX])
				} else {
					stringBuilder.WriteByte(' ')
				}
				currentX++
			}
		} else {
			currentX := 0
			for currentX < contentW {
				if currentX == 0 {
					stringBuilder.WriteByte('~')
				} else {
					stringBuilder.WriteByte(' ')
				}
				currentX++
			}
		}

		stringBuilder.WriteString("\r\n")
	}

	fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", keys.KeyLabel(last), screenW, screenH)
	cursorX, cursorY := es.Cursor()
	screenCursorY := cursorY - rowOffset
	fmt.Fprintf(&stringBuilder, "\x1b[%d;%dH", screenCursorY+1, cursorX+1) // ANSI cursor positions are 1-based
	io.WriteString(w, stringBuilder.String())
}
