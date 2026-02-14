package g0ui

import (
	"syscall"
	"unsafe"
)

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

var inputBuf []byte

func hasBufferedInput() bool {
	return len(inputBuf) > 0
}

// stdinReady returns true if stdin has data available within timeoutMs.
func stdinReady(timeoutMs int) bool {
	ts := syscall.Timespec{
		Sec:  int64(timeoutMs / 1000),
		Nsec: int64((timeoutMs % 1000) * 1000000),
	}
	var readfds [128]byte // fd_set for pselect (1024 bits)
	readfds[0] = 1        // set bit 0 (stdin)
	_, _, errno := syscall.Syscall6(
		syscall.SYS_PSELECT6,
		1, // nfds
		uintptr(unsafe.Pointer(&readfds[0])),
		0, 0,
		uintptr(unsafe.Pointer(&ts)),
		0,
	)
	return errno == 0 && readfds[0]&1 != 0
}

const frameTimeMs = 50 // 20 FPS

func readInput() InputEvent {
	// If buffer is empty, wait up to frameTimeMs for input
	if len(inputBuf) == 0 {
		if !stdinReady(frameTimeMs) {
			return InputEvent{Key: KeyNone}
		}
		var buf [64]byte
		n, err := syscall.Read(syscall.Stdin, buf[:])
		if err != nil || n == 0 {
			return InputEvent{Key: KeyNone}
		}
		inputBuf = append(inputBuf[:0], buf[:n]...)
	}

	b := inputBuf

	// Escape sequences (must check before single-byte ESC)
	if b[0] == 27 && len(b) >= 3 && b[1] == '[' {
		inputBuf = inputBuf[3:]
		switch b[2] {
		case 'A':
			return InputEvent{Key: KeyUp}
		case 'B':
			return InputEvent{Key: KeyDown}
		case 'C':
			return InputEvent{Key: KeyRight}
		case 'D':
			return InputEvent{Key: KeyLeft}
		}
		return InputEvent{Key: KeyNone}
	}

	// Multi-byte UTF-8
	if b[0] >= 0xC0 {
		size := 1
		switch {
		case b[0] < 0xE0:
			size = 2
		case b[0] < 0xF0:
			size = 3
		default:
			size = 4
		}
		if len(b) >= size {
			r := decodeUTF8(b[:size])
			inputBuf = inputBuf[size:]
			if r != 0 {
				return InputEvent{Key: KeyRune, Rune: r}
			}
			return InputEvent{Key: KeyNone}
		}
		// Not enough bytes — discard lead byte
		inputBuf = inputBuf[1:]
		return InputEvent{Key: KeyNone}
	}

	// Single byte — consume one byte from buffer
	ch := b[0]
	inputBuf = inputBuf[1:]

	switch ch {
	case 3:
		return InputEvent{Key: KeyCtrlC}
	case 13, 10:
		return InputEvent{Key: KeyEnter}
	case 32:
		return InputEvent{Key: KeySpace}
	case 9:
		return InputEvent{Key: KeyTab}
	case 27:
		return InputEvent{Key: KeyEsc}
	case 127, 8:
		return InputEvent{Key: KeyBackspace}
	default:
		if ch >= 32 && ch < 127 {
			return InputEvent{Key: KeyRune, Rune: rune(ch)}
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
