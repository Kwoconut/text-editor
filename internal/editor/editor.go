package editor

import "texteditor/internal/keys"

type Action int

const (
	ActionNone Action = iota
	ActionQuit
)

type EditorState struct {
	cursorX int
	cursorY int
	cells   map[[2]int]byte
	width   int
	height  int
}

func New(w, h int) *EditorState {
	return &EditorState{
		cursorX: 0,
		cursorY: 0,
		cells:   make(map[[2]int]byte),
		width:   w,
		height:  h,
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

func (es *EditorState) Cell(x int, y int) (byte, bool) {
	char, ok := es.cells[[2]int{x, y}]
	return char, ok
}

func (es *EditorState) UpdateSize(w, h int) {
	es.width = w
	es.height = h
	es.clampCursor()
}

func (es *EditorState) ContentHeight() int {
	if es.height <= 1 {
		return 0
	}
	return es.height - 1
}

func (es *EditorState) MaxCursorY() int {
	h := es.ContentHeight()
	if h == 0 {
		return 0
	}
	return h - 1
}

func (es *EditorState) MaxCursorX() int {
	if es.width <= 0 {
		return 0
	}
	return es.width - 1
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
			es.insert(keyEvent.Char)
		} else if keyEvent.Char == keys.BACKSPACE {
			es.backspace()
		} else if keyEvent.Char == keys.ENTER {
			es.enter()
		}
	}

	es.clampCursor()
	return ActionNone
}

func (es *EditorState) clampCursor() {
	if es.cursorX < 0 {
		es.cursorX = 0
	}
	if es.cursorX > es.MaxCursorX() {
		es.cursorX = es.MaxCursorX()
	}
	if es.cursorY < 0 {
		es.cursorY = 0
	}
	if es.cursorY > es.MaxCursorY() {
		es.cursorY = es.MaxCursorY()
	}
}

func (es *EditorState) moveLeft() {
	es.cursorX--
}

func (es *EditorState) moveRight() {
	es.cursorX++
}

func (es *EditorState) moveUp() {
	es.cursorY--
}

func (es *EditorState) moveDown() {
	es.cursorY++
}

func (es *EditorState) insert(ch byte) {
	es.cells[[2]int{es.cursorX, es.cursorY}] = ch
	es.cursorX++
	if es.cursorX > es.MaxCursorX() {
		es.cursorX = 0
		if es.cursorY < es.MaxCursorY() {
			es.cursorY++
		}
	}
}

func (es *EditorState) backspace() {
	if es.cursorX > 0 {
		es.cursorX--
	} else if es.cursorY > 0 {
		es.cursorY--
		es.cursorX = es.MaxCursorX()
	}
	delete(es.cells, [2]int{es.cursorX, es.cursorY})
}

func (es *EditorState) enter() {
	es.cursorX = 0
	if es.cursorY < es.MaxCursorY() {
		es.cursorY++
	}
}
