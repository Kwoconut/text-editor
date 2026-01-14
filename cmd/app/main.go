package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
	"texteditor/internal/keys"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	cursorX := 0
	cursorY := 0
	redraw(keys.KeyEvent{}, cursorX, cursorY, terminalWidth, terminalHeight)

	for {
		keyEvent := keys.ReadKey()

		if keyEvent.Kind == keys.KeyChar && keyEvent.Char == keys.CTRL_Q {
			break
		}

		terminalWidth, terminalHeight, _ := term.GetSize(int(os.Stdout.Fd()))

		switch keyEvent.Kind {
		case keys.KeyUp:
			if cursorY >= 1 {
				cursorY--
			}
		case keys.KeyDown:
			if cursorY < terminalHeight-2 {
				cursorY++
			}
		case keys.KeyLeft:
			if cursorX >= 1 {
				cursorX--
			}
		case keys.KeyRight:
			if cursorX < terminalWidth-1 {
				cursorX++
			}
		}

		redraw(keyEvent, cursorX, cursorY, terminalWidth, terminalHeight)
	}
}

func redraw(last keys.KeyEvent, cursorX, cursorY, terminalWidth, terminalHeight int) {

	os.Stdout.Write([]byte("\x1b[2J"))
	os.Stdout.Write([]byte("\x1b[H"))

	for i := 0; i < terminalHeight-1; i++ {
		fmt.Print("~\r\n")
	}

	switch last.Kind {
	case keys.KeyChar:
		if last.Char < 32 {
			fmt.Printf("Ctrl+Q to quit | Last key: %d | Screen size: %d:%d", last.Char, terminalHeight, terminalWidth)
		} else {
			fmt.Printf("Ctrl+Q to quit | Last key: %c | Screen size: %d:%d", last.Char, terminalHeight, terminalWidth)
		}
	case keys.KeyUp:
		fmt.Printf("Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Up", terminalHeight, terminalWidth)
	case keys.KeyDown:
		fmt.Printf("Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Down", terminalHeight, terminalWidth)
	case keys.KeyLeft:
		fmt.Printf("Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Left", terminalHeight, terminalWidth)
	case keys.KeyRight:
		fmt.Printf("Ctrl+Q to quit | Last key: %s | Screen size: %d:%d", "Right", terminalHeight, terminalWidth)
	}

	moveCursorCommand := fmt.Sprintf("\x1b[%d;%dH", cursorY+1, cursorX+1) // ANSI cursor positions are 1-based
	os.Stdout.Write([]byte(moveCursorCommand))
}
