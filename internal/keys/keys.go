package keys

import "os"

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
)

func ReadKey() KeyEvent {
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