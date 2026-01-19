package keys

import (
	"fmt"
	"io"
	"strings"
)

type Key int

type KeyEvent struct {
	Kind Key
	Char byte   // meaningful only when Kind == KeyChar
	Seq  []byte // always the exact bytes read for this key event
}

const (
	KeyChar Key = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyUnknown
)

const (
	ESCAPE    = 27
	CTRL_Q    = 17
	CTRL_S    = 19
	BACKSPACE = 127
	ENTER     = 13
)

const maxEscapeSeqLen = 8

var seqMap = map[string]Key{
	"\x1b[A":  KeyUp,
	"\x1b[B":  KeyDown,
	"\x1b[C":  KeyRight,
	"\x1b[D":  KeyLeft,
	"\x1b[H":  KeyHome,
	"\x1b[F":  KeyEnd,
	"\x1b[5~": KeyPageUp,
	"\x1b[6~": KeyPageDown,

	"\x1bOH": KeyHome,
	"\x1bOF": KeyEnd,

	"\x1b[1~": KeyHome,
	"\x1b[4~": KeyEnd,
	"\x1b[7~": KeyHome,
	"\x1b[8~": KeyEnd,
}

func ReadKey(r io.Reader) KeyEvent {
	seq, err := readOneKeySeq(r)
	if err != nil {
		return KeyEvent{Kind: KeyUnknown, Seq: seq}
	}
	if len(seq) == 0 {
		return KeyEvent{Kind: KeyUnknown, Seq: seq}
	}

	if seq[0] != ESCAPE {
		return KeyEvent{Kind: KeyChar, Char: seq[0], Seq: seq}
	}

	if k, ok := seqMap[string(seq)]; ok {
		return KeyEvent{Kind: k, Seq: seq}
	}
	return KeyEvent{Kind: KeyUnknown, Seq: seq}
}

func readOneKeySeq(r io.Reader) ([]byte, error) {
	first, err := readByte(r)
	if err != nil {
		return nil, err
	}

	seq := []byte{first}
	if first != ESCAPE {
		return seq, nil
	}

	for len(seq) < maxEscapeSeqLen {
		b, err := readByte(r)
		if err != nil {
			return seq, err
		}
		seq = append(seq, b)

		if b == '~' || isASCIILetter(b) {
			break
		}
	}

	return seq, nil
}

func readByte(r io.Reader) (byte, error) {
	var buf [1]byte
	_, err := io.ReadFull(r, buf[:])
	return buf[0], err
}

func isASCIILetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func KeyLabel(ev KeyEvent) string {
	switch ev.Kind {
	case KeyChar:
		switch ev.Char {
		case ENTER:
			return "ENTER"
		case CTRL_S:
			return "CTRL_S"
		case CTRL_Q:
			return "CTRL_Q"
		case BACKSPACE:
			return "BACKSPACE"
		default:
			if ev.Char < 32 {
				return fmt.Sprintf("%d", ev.Char)
			}
			return string(ev.Char)
		}
	case KeyUp:
		return "UP"
	case KeyDown:
		return "DOWN"
	case KeyLeft:
		return "LEFT"
	case KeyRight:
		return "RIGHT"
	case KeyHome:
		return "HOME"
	case KeyEnd:
		return "END"
	case KeyPageUp:
		return "PAGEUP"
	case KeyPageDown:
		return "PAGEDOWN"
	default:
		return "UNKNOWN"
	}
}

func SeqHex(seq []byte) string {
	var b strings.Builder
	for i, x := range seq {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%02x", x)
	}
	return b.String()
}
