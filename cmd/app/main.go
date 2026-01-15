package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
	"texteditor/internal/keys"
	"texteditor/internal/editor"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	editorState := editor.New(terminalWidth, terminalHeight)
	redraw(keys.KeyEvent{}, editorState)

	for {
		terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
		editorState.UpdateSize(terminalWidth, terminalHeight)

		keyEvent := keys.ReadKey()
		shouldQuit := editorState.HandleKey(keyEvent)
		if shouldQuit {
			break
		}

		redraw(keyEvent, editorState)
	}
}

func redraw(last keys.KeyEvent, editorState *editor.EditorState) {
	var stringBuilder strings.Builder
	stringBuilder.Write([]byte("\x1b[2J"))
	stringBuilder.Write([]byte("\x1b[H"))

	for y := 0; y < editorState.Height()-1; y++ {
		for x := 0; x < editorState.Width(); x++ {
			if char, ok := editorState.Cell(x, y); ok {
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
		if last.Char == keys.ENTER {
			fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "ENTER", editorState.Height(), editorState.Width())
		}else if last.Char < 32 && last.Char != keys.ENTER {
			fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %d | Screen size: %d:%d", last.Char, editorState.Height(), editorState.Width())
		} else if last.Char < 127 {
			fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %c | Screen size: %d:%d", last.Char, editorState.Height(), editorState.Width())
		} else if last.Char == keys.BACKSPACE {
			fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "BACKSPACE", editorState.Height(), editorState.Width())
		}
	case keys.KeyUp:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Up", editorState.Height(), editorState.Width())
	case keys.KeyDown:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Down", editorState.Height(), editorState.Width())
	case keys.KeyLeft:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Left", editorState.Height(), editorState.Width())
	case keys.KeyRight:
		fmt.Fprintf(&stringBuilder, "Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Right", editorState.Height(), editorState.Width())
	}

	cursorX, cursorY:= editorState.Cursor()
	fmt.Fprintf(&stringBuilder, "\x1b[%d;%dH", cursorY+1, cursorX+1) // ANSI cursor positions are 1-based
	os.Stdout.WriteString(stringBuilder.String())
}
