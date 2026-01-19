package editor

import "texteditor/internal/keys"
import "strings"

type Action int

const (
	TAB_STOP = 4
)

const (
	ActionNone Action = iota
	ActionQuit
	ActionSave
)

type EditorState struct {
	cursorX    int
	cursorY    int
	preferredX int
	lines      [][]rune
	width      int
	height     int
	rowOffset  int
	colOffset  int
	isDirty    bool
}

func New(w, h int, text string) *EditorState {
	return &EditorState{
		cursorX:    0,
		cursorY:    0,
		preferredX: 0,
		lines:      initializeText(text),
		width:      w,
		height:     h,
		rowOffset:  0,
		colOffset:  0,
		isDirty:    false,
	}
}

func (es *EditorState) Width() int {
	return es.width
}

func (es *EditorState) Height() int {
	return es.height
}

func (es *EditorState) Cursor() (int, int) {
	return es.cursorX, es.cursorY
}

func (es *EditorState) CursorRX() int {
	line := es.lines[es.cursorY]
	return cxToRx(line, es.cursorX, TAB_STOP)
}

func (es *EditorState) UpdateSize(w, h int) {
	es.width = w
	es.height = h
	es.clampCursor()
	es.adjustRowOffset()
	es.adjustColOffset()
}

func (es *EditorState) ContentHeight() int {
	if es.height <= 1 {
		return 0
	}
	return es.height - 1
}

func (es *EditorState) ContentWidth() int {
	if es.width <= 0 {
		return 0
	}

	return es.width
}

func (es *EditorState) RowOffset() int {
	return es.rowOffset
}

func (es *EditorState) ColOffset() int {
	return es.colOffset
}

func (es *EditorState) LineCount() int {
	return len(es.lines)
}

func (es *EditorState) Line(y int) []rune {
	return es.lines[y]
}

func (es *EditorState) Text() string {
	var b strings.Builder

	for i, line := range es.lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(string(line))
	}

	return b.String()
}

func (es *EditorState) MarkSaved() {
	es.isDirty = false
}

func (es *EditorState) IsDirty() bool {
	return es.isDirty
}

func (es *EditorState) HandleKey(keyEvent keys.KeyEvent) Action {
	if keyEvent.Kind == keys.KeyChar && keyEvent.Char == keys.CTRL_Q {
		return ActionQuit
	}

	switch keyEvent.Kind {
	case keys.KeyLeft:
		es.moveLeft()
	case keys.KeyRight:
		es.moveRight()
	case keys.KeyUp:
		es.moveUp()
	case keys.KeyDown:
		es.moveDown()
	case keys.KeyHome:
		es.home()
	case keys.KeyEnd:
		es.end()
	case keys.KeyPageDown:
		es.pageDown()
	case keys.KeyPageUp:
		es.pageUp()
	case keys.KeyChar:
		if keyEvent.Char >= 32 && keyEvent.Char <= 126 {
			es.insert(rune(keyEvent.Char))
			es.isDirty = true
		} else if keyEvent.Char == keys.BACKSPACE {
			es.backspace()
			es.isDirty = true
		} else if keyEvent.Char == keys.ENTER {
			es.enter()
			es.isDirty = true
		} else if keyEvent.Char == keys.CTRL_S {
			return ActionSave
		}
	}

	es.clampCursor()
	es.adjustRowOffset()
	es.adjustColOffset()
	return ActionNone
}

func initializeText(text string) [][]rune {
	var lines [][]rune

	readLines := strings.Split(text, "\n")

	for index, readLine := range readLines {
		if index == len(readLines)-1 && readLine == "" {
			break
		}
		lines = append(lines, []rune(readLine))
	}
	if len(lines) == 0 {
		lines = [][]rune{[]rune{}}
	}

	return lines
}

func (es *EditorState) adjustRowOffset() {
	if es.ContentHeight() <= 0 {
		es.rowOffset = 0
		return
	}

	maxTop := max(0, es.LineCount()-es.ContentHeight())

	if es.cursorY < es.rowOffset {
		es.rowOffset = es.cursorY
	} else if es.cursorY >= es.rowOffset+es.ContentHeight() {
		es.rowOffset = es.cursorY - es.ContentHeight() + 1
	}

	if es.rowOffset > maxTop {
		es.rowOffset = maxTop
	}
}

func (es *EditorState) adjustColOffset() {
	if es.ContentWidth() <= 0 {
		es.colOffset = 0
		return
	}

	line := es.lines[es.cursorY]
	rx := cxToRx(line, es.cursorX, TAB_STOP)

	if rx < es.colOffset {
		es.colOffset = rx
	} else if rx >= es.colOffset+es.ContentWidth() {
		es.colOffset = rx - es.ContentWidth() + 1
	}

	if es.colOffset < 0 {
		es.colOffset = 0
	}
}

func cxToRx(line []rune, cx int, tabStop int) int {
	if cx < 0 {
		return 0
	}
	if cx > len(line) {
		cx = len(line)
	}

	rx := 0
	for i := 0; i < cx; i++ {
		if line[i] == '\t' {
			spaces := tabStop - (rx % tabStop)
			rx += spaces
		} else {
			rx += 1
		}
	}

	return rx
}

func (es *EditorState) clampCursor() {
	if len(es.lines) == 0 {
		es.lines = [][]rune{[]rune{}}
	}

	if es.cursorY < 0 {
		es.cursorY = 0
	}

	if es.cursorY >= len(es.lines) {
		es.cursorY = len(es.lines) - 1
	}

	if es.cursorX < 0 {
		es.cursorX = 0
	}

	if es.cursorX > len(es.lines[es.cursorY]) {
		es.cursorX = len(es.lines[es.cursorY])
	}
}

func (es *EditorState) moveLeft() {
	if es.cursorX > 0 {
		es.cursorX--
	} else if es.cursorY > 0 {
		es.cursorY--
		es.cursorX = len(es.lines[es.cursorY])
	}

	es.preferredX = es.cursorX
}

func (es *EditorState) moveRight() {
	if es.cursorX < len(es.lines[es.cursorY]) {
		es.cursorX++
	} else if es.cursorY < len(es.lines)-1 {
		es.cursorY++
		es.cursorX = 0
	}

	es.preferredX = es.cursorX
}

func (es *EditorState) moveUp() {
	es.cursorY--

	if es.cursorY < 0 {
		es.cursorY = 0
	}

	es.cursorX = min(es.preferredX, len(es.lines[es.cursorY]))
}

func (es *EditorState) moveDown() {
	es.cursorY++

	if es.cursorY >= len(es.lines) {
		es.cursorY = len(es.lines) - 1
	}

	es.cursorX = min(es.preferredX, len(es.lines[es.cursorY]))
}

func (es *EditorState) insert(ch rune) {
	line := es.lines[es.cursorY]

	x := es.cursorX
	line = append(line, 0)
	copy(line[x+1:], line[x:])
	line[x] = ch

	es.lines[es.cursorY] = line
	es.cursorX++
	es.preferredX = es.cursorX
}

func (es *EditorState) backspace() {
	if es.cursorX == 0 && es.cursorY == 0 {
		return
	}

	if es.cursorX > 0 {
		line := es.lines[es.cursorY]
		x := es.cursorX
		copy(line[x-1:], line[x:])
		line = line[:len(line)-1]
		es.lines[es.cursorY] = line
		es.cursorX--
		es.preferredX = es.cursorX
		return
	}

	previousLine := es.lines[es.cursorY-1]
	oldLen := len(previousLine)
	currentLine := es.lines[es.cursorY]
	newPreviousLine := append(previousLine, currentLine...)
	es.lines[es.cursorY-1] = newPreviousLine
	es.lines = append(es.lines[:es.cursorY], es.lines[es.cursorY+1:]...)
	es.cursorY--
	es.cursorX = oldLen
	es.preferredX = es.cursorX
}

func (es *EditorState) enter() {
	line := es.lines[es.cursorY]
	left := line[:es.cursorX]
	right := line[es.cursorX:]
	es.lines[es.cursorY] = left
	es.lines = append(es.lines[:es.cursorY+1], append([][]rune{right}, es.lines[es.cursorY+1:]...)...)
	es.cursorY++
	es.cursorX = 0
	es.preferredX = es.cursorX
}

func (es *EditorState) home() {
	line := es.lines[es.cursorY]

	if len(line) <= 0 {
		es.cursorX = 0
		es.preferredX = es.cursorX
		return
	}

	indentX := 0
	for indentX < len(line) {
		if line[indentX] != ' ' && line[indentX] != '\t' {
			break
		}
		indentX++
	}

	if es.cursorX == indentX {
		es.cursorX = 0
	} else {
		es.cursorX = indentX
	}

	es.preferredX = es.cursorX
}

func (es *EditorState) end() {
	line := es.lines[es.cursorY]
	es.cursorX = len(line)
	es.preferredX = es.cursorX
}

func (es *EditorState) pageDown() {
	contentH := es.ContentHeight()
	var pageSize int

	if contentH <= 1 {
		pageSize = 1
	} else {
		pageSize = contentH - 1
	}

	es.cursorY += pageSize

	if es.cursorY >= len(es.lines) {
		es.cursorY = len(es.lines) - 1
	}

	es.cursorX = min(es.preferredX, len(es.lines[es.cursorY]))
}

func (es *EditorState) pageUp() {
	contentH := es.ContentHeight()
	var pageSize int

	if contentH <= 1 {
		pageSize = 1
	} else {
		pageSize = contentH - 1
	}

	es.cursorY -= pageSize

	if es.cursorY < 0 {
		es.cursorY = 0
	}

	es.cursorX = min(es.preferredX, len(es.lines[es.cursorY]))
}
