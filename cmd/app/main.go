package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Key int

const (
	KeyChar Key = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyUnknown
)

const (
	ESCAPE             = 27
	SQUARE_PARENTHESES = 91
	CTRL_Q             = 17
)

type KeyEvent struct {
	Kind Key
	Char byte
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

readKeyLoop:
	for {
		keyEvent := readKey()

		switch keyEvent.Kind {
		case KeyChar:
			if keyEvent.Char == CTRL_Q {
				break readKeyLoop
			} else {
				fmt.Printf("%c\r\n", keyEvent.Char)
			}
		case KeyUp:
			fmt.Print("UP\r\n")
		case KeyDown:
			fmt.Print("DOWN\r\n")
		case KeyLeft:
			fmt.Print("LEFT\r\n")
		case KeyRight:
			fmt.Print("RIGHT\r\n")
		}
	}
}

func readKey() KeyEvent {
	var b [1]byte

	_, err := os.Stdin.Read(b[:])

	if err != nil {
		return KeyEvent{Kind: KeyUnknown}
	}

	if b[0] != ESCAPE {
		return KeyEvent{Kind: KeyChar, Char: b[0]}
	} else {
		_, err := os.Stdin.Read(b[:])

		if err != nil {
			return KeyEvent{Kind: KeyUnknown}
		}

		if b[0] == SQUARE_PARENTHESES {
			_, err := os.Stdin.Read(b[:])

			if err != nil {
				return KeyEvent{Kind: KeyUnknown}
			}

			switch b[0] {
			case 65:
				return KeyEvent{Kind: KeyUp}
			case 66:
				return KeyEvent{Kind: KeyDown}
			case 67:
				return KeyEvent{Kind: KeyRight}
			case 68:
				return KeyEvent{Kind: KeyLeft}
			}
		} else {
			return KeyEvent{Kind: KeyUnknown}
		}
	}

	return KeyEvent{Kind: KeyUnknown}
}

func redraw(last KeyEvent) {
	
}
