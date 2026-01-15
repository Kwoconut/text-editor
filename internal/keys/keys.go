package keys

import (
	"fmt"
	"io"
)

type Key int

type KeyEvent struct {
	Kind Key
	Char byte
}

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
	BACKSPACE          = 127
	ENTER 						 = 13
)

func ReadKey(r io.Reader) KeyEvent {
	var b [1]byte

	_, err := r.Read(b[:])

	if err != nil {
		return KeyEvent{Kind: KeyUnknown}
	}

	if b[0] != ESCAPE {
		return KeyEvent{Kind: KeyChar, Char: b[0]}
	} else {
		_, err := r.Read(b[:])

		if err != nil {
			return KeyEvent{Kind: KeyUnknown}
		}

		if b[0] == SQUARE_PARENTHESES {
			_, err := r.Read(b[:])

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

func KeyLabel(keyEvent KeyEvent) string {
	switch keyEvent.Kind {
	case KeyChar:
		if keyEvent.Char == ENTER {
			return "ENTER"
		} else if keyEvent.Char < 32 && keyEvent.Char != ENTER {
			return fmt.Sprintf("%d", keyEvent.Char)
		} else if keyEvent.Char < 127 {
			return string(keyEvent.Char)
		} else if keyEvent.Char == BACKSPACE {
			return "BACKSPACE"
		}
	case KeyUp:
		return "UP"
	case KeyDown:
		return "DOWN"
	case KeyLeft:
		return "LEFT"
	case KeyRight:
		return "RIGHT"
	}
	return "UNKNOWN"
}