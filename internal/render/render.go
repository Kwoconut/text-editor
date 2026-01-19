package render

import (
	"fmt"
	"io"
	"strings"

	"texteditor/internal/editor"
	"texteditor/internal/keys"
)

func Draw(w io.Writer, es *editor.EditorState, last keys.KeyEvent, statusMsg string) {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("\x1b[2J")
	stringBuilder.WriteString("\x1b[H")

	screenW := es.Width()
	screenH := es.Height()
	contentW := screenW
	contentH := screenH - 1
	rowOffset := es.RowOffset()
	colOffset := es.ColOffset()

	for y := 0; y < contentH; y++ {
		docY := rowOffset + y
		if docY < es.LineCount() {
			currentLine := es.Line(docY)

			currentX := colOffset
			for currentX < contentW + colOffset {
				if currentX < len(currentLine) {
					stringBuilder.WriteRune(currentLine[currentX])
				} else {
					stringBuilder.WriteByte(' ')
				}
				currentX++
			}
		} else {
			currentX := colOffset
			for currentX < contentW + colOffset {
				if currentX == colOffset {
					stringBuilder.WriteByte('~')
				} else {
					stringBuilder.WriteByte(' ')
				}
				currentX++
			}
		}

		stringBuilder.WriteString("\r\n")
	}

	fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d |", keys.KeyLabel(last), screenW, screenH)
	if es.IsDirty() {
		fmt.Fprint(&stringBuilder, " * ")
	} else {
		fmt.Fprint(&stringBuilder, "   ")
	}
	fmt.Fprintf(&stringBuilder, "| %s", statusMsg)

	cursorX, cursorY := es.Cursor()
	screenCursorY := cursorY - rowOffset
	screenCursorX := cursorX - colOffset
	fmt.Fprintf(&stringBuilder, "\x1b[%d;%dH", screenCursorY+1, screenCursorX+1) // ANSI cursor positions are 1-based
	io.WriteString(w, stringBuilder.String())
}
