package editor

import "texteditor/internal/keys"

type Action int

const (
	ActionNone Action = iota
	ActionQuit
)

type EditorState struct {
	cursorX   int
	cursorY   int
	lines     [][]rune
	width     int
	height    int
	rowOffset int
}

func New(w, h int) *EditorState {
	return &EditorState{
		cursorX:   0,
		cursorY:   0,
		lines:     [][]rune{[]rune{}},
		width:     w,
		height:    h,
		rowOffset: 0,
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

func (es *EditorState) UpdateSize(w, h int) {
	es.width = w
	es.height = h
	es.clampCursor()
	es.adjustRowOffset()
}

func (es *EditorState) ContentHeight() int {
	if es.height <= 1 {
		return 0
	}
	return es.height - 1
}

func (es *EditorState) RowOffset() int {
	return es.rowOffset
}

func (es *EditorState) LineCount() int {
	return len(es.lines)
}

func (es *EditorState) Line(y int) []rune {
	return es.lines[y]
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
	case keys.KeyChar:
		if keyEvent.Char >= 32 && keyEvent.Char <= 126 {
			es.insert(rune(keyEvent.Char))
		} else if keyEvent.Char == keys.BACKSPACE {
			es.backspace()
		} else if keyEvent.Char == keys.ENTER {
			es.enter()
		}
	}


	es.clampCursor()
	es.adjustRowOffset()
	return ActionNone
}

func (es *EditorState) adjustRowOffset() {
	if es.ContentHeight() <= 0 {
		es.rowOffset = 0
		return
	}

	maxTop := max(0, es.LineCount() - es.ContentHeight())


	if es.cursorY < es.rowOffset {
		es.rowOffset = es.cursorY
	} else if es.cursorY >= es.rowOffset+es.ContentHeight() {
		es.rowOffset = es.cursorY - es.ContentHeight() + 1
	}

	if es.rowOffset > maxTop {
		es.rowOffset = maxTop
	}
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
}

func (es *EditorState) moveRight() {
	if es.cursorX < len(es.lines[es.cursorY]) {
		es.cursorX++
	} else if es.cursorY < len(es.lines)-1 {
		es.cursorY++
		es.cursorX = 0
	}
}

func (es *EditorState) moveUp() {
	es.cursorY--
}

func (es *EditorState) moveDown() {
	es.cursorY++
}

func (es *EditorState) insert(ch rune) {
	line := es.lines[es.cursorY]

	x := es.cursorX
	line = append(line, 0)
	copy(line[x+1:], line[x:])
	line[x] = ch

	es.lines[es.cursorY] = line
	es.cursorX++
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
}

func (es *EditorState) enter() {
	line := es.lines[es.cursorY]
	left := line[:es.cursorX]
	right := line[es.cursorX:]
	es.lines[es.cursorY] = left
	es.lines = append(es.lines[:es.cursorY+1], append([][]rune{right}, es.lines[es.cursorY+1:]...)...)
	es.cursorY++
	es.cursorX = 0
}
