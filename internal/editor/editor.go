package editor

import "texteditor/internal/keys"

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
	if char, ok := es.cells[[2]int{x, y}]; ok {
		return char, ok
	}

	return 0, false
}

func (es *EditorState) UpdateSize(w, h int) {
	es.width = w
	es.height = h
}

func (es *EditorState) HandleKey(keyEvent keys.KeyEvent) bool {
	if keyEvent.Kind == keys.KeyChar && keyEvent.Char == keys.CTRL_Q {
		return true
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

	return false
}

func (es *EditorState) moveLeft() {
	es.cursorX--
	if es.cursorX < 0 {
		es.cursorX = 0
	}
}

func (es *EditorState) moveRight() {
	es.cursorX++
	if es.cursorX > es.width-1 {
		es.cursorX = es.width - 1
	}
}

func (es *EditorState) moveUp() {
	es.cursorY--
	if es.cursorY < 0 {
		es.cursorY = 0
	}
}

func (es *EditorState) moveDown() {
	es.cursorY++
	if es.cursorY > es.height-2 {
		es.cursorY = es.height - 2
	}
}

func (es *EditorState) insert(ch byte) {
	es.cells[[2]int{es.cursorX, es.cursorY}] = ch
	es.cursorX++
	if es.cursorX >= es.width {
		es.cursorX = 0
		if es.cursorY < es.height-2 {
			es.cursorY++
		}
	}
}

func (es *EditorState) backspace() {
	if es.cursorX != 0 || es.cursorY != 0 {
		if es.cursorX >= 0 {
			es.cursorX--
			if es.cursorX < 0 && es.cursorY > 0 {
				es.cursorY--
				es.cursorX = es.width - 1
			}
		}
	}
	delete(es.cells, [2]int{es.cursorX, es.cursorY})
}

func (es *EditorState) enter() {
	es.cursorX = 0
	if es.cursorY < es.height-2 {
		es.cursorY++
	}
}
