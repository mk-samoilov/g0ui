package g0ui

import "syscall"

// Key represents a keyboard input.
type Key int

const (
	KeyNone Key = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyEnter
	KeySpace
	KeyTab
	KeyEsc
	KeyCtrlC
	KeyBackspace
	KeyRune
)

// InputEvent holds the result of a single keypress.
type InputEvent struct {
	Key  Key
	Rune rune
}

func readInput() InputEvent {
	var buf [16]byte
	n, err := syscall.Read(syscall.Stdin, buf[:])
	if err != nil || n == 0 {
		return InputEvent{Key: KeyNone}
	}

	// Single byte
	if n == 1 {
		switch buf[0] {
		case 3: // Ctrl+C
			return InputEvent{Key: KeyCtrlC}
		case 13, 10: // Enter
			return InputEvent{Key: KeyEnter}
		case 32: // Space
			return InputEvent{Key: KeySpace}
		case 9: // Tab
			return InputEvent{Key: KeyTab}
		case 27: // Esc (standalone)
			return InputEvent{Key: KeyEsc}
		case 127, 8: // Backspace
			return InputEvent{Key: KeyBackspace}
		default:
			if buf[0] >= 32 && buf[0] < 127 {
				return InputEvent{Key: KeyRune, Rune: rune(buf[0])}
			}
			return InputEvent{Key: KeyNone}
		}
	}

	// Escape sequences
	if n >= 3 && buf[0] == 27 && buf[1] == '[' {
		switch buf[2] {
		case 'A':
			return InputEvent{Key: KeyUp}
		case 'B':
			return InputEvent{Key: KeyDown}
		case 'C':
			return InputEvent{Key: KeyRight}
		case 'D':
			return InputEvent{Key: KeyLeft}
		}
	}

	// Multi-byte UTF-8
	if buf[0] >= 0x80 {
		r := decodeUTF8(buf[:n])
		if r != 0 {
			return InputEvent{Key: KeyRune, Rune: r}
		}
	}

	return InputEvent{Key: KeyNone}
}

func decodeUTF8(b []byte) rune {
	if len(b) == 0 {
		return 0
	}
	// Simple UTF-8 decode
	c := b[0]
	switch {
	case c < 0x80:
		return rune(c)
	case c < 0xE0 && len(b) >= 2:
		return rune(c&0x1F)<<6 | rune(b[1]&0x3F)
	case c < 0xF0 && len(b) >= 3:
		return rune(c&0x0F)<<12 | rune(b[1]&0x3F)<<6 | rune(b[2]&0x3F)
	case c < 0xF8 && len(b) >= 4:
		return rune(c&0x07)<<18 | rune(b[1]&0x3F)<<12 | rune(b[2]&0x3F)<<6 | rune(b[3]&0x3F)
	}
	return 0
}
