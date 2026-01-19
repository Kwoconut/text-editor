package render

import (
	"fmt"
	"io"
	"strings"

	"texteditor/internal/editor"
	"texteditor/internal/keys"
)

func Draw(w io.Writer, es *editor.EditorState, last keys.KeyEvent, statusMsg string) {
	screenW := es.Width()
	screenH := es.Height()
	if screenW <= 0 || screenH <= 0 {
		return
	}

	contentW := screenW
	contentH := screenH - 1
	if contentH < 0 {
		contentH = 0
	}

	rowOffset := es.RowOffset()
	colOffset := es.ColOffset()

	var b strings.Builder
	// Clear + home. Hide cursor during redraw to reduce flicker.
	b.WriteString("\x1b[2J")
	b.WriteString("\x1b[H")

	// Draw content area
	for screenY := 0; screenY < contentH; screenY++ {
		docY := rowOffset + screenY
		if docY >= 0 && docY < es.LineCount() {
			drawLine(&b, es.Line(docY), colOffset, contentW)
		} else {
			drawTildeRow(&b, contentW)
		}
		b.WriteString("\r\n")
	}

	// Status bar (pad/truncate to full width)
	status := buildStatusLine(es, last, statusMsg)
	status = fitWidth(status, screenW)
	b.WriteString(status)

	_, cursorY := es.Cursor()
	screenCursorY := cursorY - rowOffset

	rx := es.CursorRX()
	screenCursorX := rx - colOffset

	if contentH > 0 {
		screenCursorY = clamp(screenCursorY, 0, contentH-1)
	} else {
		screenCursorY = 0
	}
	if contentW > 0 {
		screenCursorX = clamp(screenCursorX, 0, contentW-1)
	} else {
		screenCursorX = 0
	}

	fmt.Fprintf(&b, "\x1b[%d;%dH", screenCursorY+1, screenCursorX+1) // 1-based
	io.WriteString(w, b.String())
}

func expandTabs(line []rune, tabStop int) []rune {
	var out []rune
	col := 0
	for _, ch := range line {
		if ch == '\t' {
			spaces := tabStop - (col % tabStop)
			for i := 0; i < spaces; i++ {
				out = append(out, ' ')
				col++
			}
		} else {
			out = append(out, ch)
			col++
		}
	}
	return out
}

func drawLine(b *strings.Builder, line []rune, colOffset, width int) {
	rendered := expandTabs(line, 4)
	lineLen := len(rendered)

	for screenX := 0; screenX < width; screenX++ {
		docX := colOffset + screenX
		if docX >= 0 && docX < lineLen {
			b.WriteRune(rendered[docX])
		} else {
			b.WriteByte(' ')
		}
	}
}

func drawTildeRow(b *strings.Builder, width int) {
	for screenX := 0; screenX < width; screenX++ {
		if screenX == 0 {
			b.WriteByte('~')
		} else {
			b.WriteByte(' ')
		}
	}
}

func buildStatusLine(es *editor.EditorState, last keys.KeyEvent, statusMsg string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d | Seq: %s |",
		keys.KeyLabel(last), es.Width(), es.Height(), keys.SeqHex(last.Seq))

	if es.IsDirty() {
		sb.WriteString(" * ")
	} else {
		sb.WriteString("   ")
	}

	if statusMsg != "" {
		sb.WriteString("| ")
		sb.WriteString(statusMsg)
	}
	return sb.String()
}

func fitWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) > width {
		return string(r[:width])
	}
	if len(r) < width {
		var b strings.Builder
		b.WriteString(string(r))
		for i := 0; i < width-len(r); i++ {
			b.WriteByte(' ')
		}
		return b.String()
	}
	return s
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
