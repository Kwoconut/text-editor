package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
	"texteditor/internal/keys"
)

type EditorState struct {
	cursorX int
	cursorY int
	cells   map[[2]int]byte
	width   int
	height  int
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	editorState := EditorState{cursorX: 0, cursorY: 0, cells: make(map[[2]int]byte), width: terminalWidth, height: terminalHeight}
	redraw(keys.KeyEvent{}, editorState)

	for {
		keyEvent := keys.ReadKey()

		if keyEvent.Kind == keys.KeyChar && keyEvent.Char == keys.CTRL_Q {
			break
		}

		terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
		editorState.width = terminalWidth
		editorState.height = terminalHeight

		switch keyEvent.Kind {
		case keys.KeyLeft:
			editorState.cursorX--
			if editorState.cursorX < 0 {
				editorState.cursorX = 0
			}
		case keys.KeyRight:
			editorState.cursorX++
			if editorState.cursorX > editorState.width-1 {
				editorState.cursorX = editorState.width - 1
			}
		case keys.KeyUp:
			editorState.cursorY--
			if editorState.cursorY < 0 {
				editorState.cursorY = 0
			}
		case keys.KeyDown:
			editorState.cursorY++
			if editorState.cursorY > editorState.height-2 {
				editorState.cursorY = editorState.height - 2
			}
		case keys.KeyChar:
			if keyEvent.Char >= 32 && keyEvent.Char <= 126 {
				editorState.cells[[2]int{editorState.cursorX, editorState.cursorY}] = keyEvent.Char
				editorState.cursorX++
				if editorState.cursorX >= editorState.width-1 {
					editorState.cursorX = 0
					if editorState.cursorY < editorState.height-2 {
						editorState.cursorY++
					}
				}
			} else if keyEvent.Char == keys.BACKSPACE {
				delete(editorState.cells, [2]int{editorState.cursorX, editorState.cursorY})
				if (editorState.cursorX != 0 || editorState.cursorY != 0) {
					if editorState.cursorX >= 0 {
						editorState.cursorX--
						if editorState.cursorX < 0 {
							editorState.cursorY--
							editorState.cursorX = editorState.width - 1
						}
					}
				}
			}
		}

		redraw(keyEvent, editorState)
	}
}

func redraw(last keys.KeyEvent, editorState EditorState) {
	var stringBuilder strings.Builder
	stringBuilder.Write([]byte("\x1b[2J"))
	stringBuilder.Write([]byte("\x1b[H"))

	for y := 0; y < editorState.height-1; y++ {
		for x := 0; x < editorState.width; x++ {
			if char, ok := editorState.cells[[2]int{x, y}]; ok {
				stringBuilder.WriteByte(char)
			} else if x == 0 {
				stringBuilder.WriteByte('~')
			} else {
				stringBuilder.WriteByte(' ')
			}
		}

		stringBuilder.WriteString("\r\n")
	}

	switch last.Kind {
	case keys.KeyChar:
		if last.Char < 32 {
			fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %d | Screen size: %d:%d", last.Char, editorState.height, editorState.width)
		} else if last.Char < 127 {
			fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %c | Screen size: %d:%d", last.Char, editorState.height, editorState.width)
		} else if last.Char == 127 {
			fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "BACKSPACE", editorState.height, editorState.width)
		}
	case keys.KeyUp:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Up", editorState.height, editorState.width)
	case keys.KeyDown:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Down", editorState.height, editorState.width)
	case keys.KeyLeft:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Left", editorState.height, editorState.width)
	case keys.KeyRight:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Right", editorState.height, editorState.width)
	}

	fmt.Fprintf(&stringBuilder, "\x1b[%d;%dH", editorState.cursorY+1, editorState.cursorX+1) // ANSI cursor positions are 1-based
	os.Stdout.WriteString(stringBuilder.String())
}
