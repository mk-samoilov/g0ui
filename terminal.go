package g0ui

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// termios mirrors the C struct termios for terminal configuration.
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

var origTermios termios
var termInitialized bool

func enableRawMode() error {
	if termInitialized {
		return nil
	}
	// tcgetattr
	if err := tcgetattr(syscall.Stdin, &origTermios); err != nil {
		return fmt.Errorf("tcgetattr: %w", err)
	}

	raw := origTermios
	// Input flags: disable BRKINT, ICRNL, INPCK, ISTRIP, IXON
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	// Output flags: disable OPOST
	raw.Oflag &^= syscall.OPOST
	// Control flags: set CS8
	raw.Cflag |= syscall.CS8
	// Local flags: disable ECHO, ICANON, IEXTEN, ISIG
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	// Min bytes and timeout
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	if err := tcsetattr(syscall.Stdin, &raw); err != nil {
		return fmt.Errorf("tcsetattr: %w", err)
	}
	termInitialized = true
	return nil
}

func disableRawMode() {
	if !termInitialized {
		return
	}
	tcsetattr(syscall.Stdin, &origTermios)
	termInitialized = false
}

func tcgetattr(fd int, t *termios) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(t)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

func tcsetattr(fd int, t *termios) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(t)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

func enterAltScreen() {
	os.Stdout.WriteString("\x1b[?1049h")
}

func exitAltScreen() {
	os.Stdout.WriteString("\x1b[?1049l")
}

func hideCursor() {
	os.Stdout.WriteString("\x1b[?25l")
}

func showCursor() {
	os.Stdout.WriteString("\x1b[?25h")
}

func clearScreen() {
	os.Stdout.WriteString("\x1b[2J\x1b[H")
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getTermSize() (w, h int) {
	var ws winsize
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 || ws.Col == 0 {
		return 80, 24 // fallback
	}
	return int(ws.Col), int(ws.Row)
}
