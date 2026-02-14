package g0ui

import (
	"bytes"
	"fmt"
	"os"
)

// cell is one character on screen.
type cell struct {
	ch           rune
	inverse      bool // currently rendered as inverse
	focusInverse bool // should invert when this cell's widget is focused
}

// screenBuf is the full-screen character buffer.
type screenBuf struct {
	cells [][]cell
	w, h  int
}

func newScreenBuf(w, h int) *screenBuf {
	cells := make([][]cell, h)
	for y := 0; y < h; y++ {
		row := make([]cell, w)
		for x := 0; x < w; x++ {
			row[x] = cell{ch: ' '}
		}
		cells[y] = row
	}
	return &screenBuf{cells: cells, w: w, h: h}
}

func (s *screenBuf) set(x, y int, c cell) {
	if x >= 0 && x < s.w && y >= 0 && y < s.h {
		s.cells[y][x] = c
	}
}

func (s *screenBuf) setRune(x, y int, ch rune) {
	if x >= 0 && x < s.w && y >= 0 && y < s.h {
		s.cells[y][x].ch = ch
	}
}

// drawBorder draws the rounded window border with title.
func drawBorder(buf *screenBuf, title string) {
	w, h := buf.w, buf.h
	if w < 4 || h < 3 {
		return
	}

	// Top: ╭─── TITLE ────╮
	buf.setRune(0, 0, '╭')
	buf.setRune(w-1, 0, '╮')
	for x := 1; x < w-1; x++ {
		buf.setRune(x, 0, '─')
	}

	// Insert title
	if title != "" {
		titleRunes := []rune("─── " + title + " ")
		startX := 1
		for i, r := range titleRunes {
			if startX+i >= w-1 {
				break
			}
			buf.setRune(startX+i, 0, r)
		}
	}

	// Bottom: ╰──────────────╯
	buf.setRune(0, h-1, '╰')
	buf.setRune(w-1, h-1, '╯')
	for x := 1; x < w-1; x++ {
		buf.setRune(x, h-1, '─')
	}

	// Sides: │
	for y := 1; y < h-1; y++ {
		buf.setRune(0, y, '│')
		buf.setRune(w-1, y, '│')
	}
}

// renderFrame renders the current frame to the terminal.
func renderFrame(c *frameContext) {
	buf := newScreenBuf(c.termW, c.termH)
	drawBorder(buf, c.title)

	lines := layoutWidgets(c)

	innerH := c.termH - 2 // content area height (minus top/bottom border)
	if innerH < 1 {
		innerH = 1
	}
	innerW := c.termW - 4 // content area width
	if innerW < 1 {
		innerW = 1
	}

	// Adjust scroll so focused widget is visible
	adjustScroll(c, lines, innerH)

	// Draw visible lines
	endY := c.scrollY + innerH
	if endY > len(lines) {
		endY = len(lines)
	}
	for i := c.scrollY; i < endY; i++ {
		screenY := 1 + (i - c.scrollY) // +1 for top border
		line := lines[i]
		isFocused := line.focusID >= 0 && line.focusID == c.focusIndex

		for x, cl := range line.cells {
			if x >= innerW {
				break
			}
			drawCell := cl
			// If this cell's widget is focused and the cell is marked for focus-inverse
			if isFocused && cl.focusInverse {
				drawCell.inverse = true
			}
			buf.set(x+2, screenY, drawCell) // +2 for border + padding
		}
	}

	flush(buf)
}

// adjustScroll ensures the focused widget is visible.
func adjustScroll(c *frameContext, lines []renderLine, innerH int) {
	// Find first and last line of focused widget
	firstLine := -1
	lastLine := -1
	for i, l := range lines {
		if l.focusID == c.focusIndex {
			if firstLine < 0 {
				firstLine = i
			}
			lastLine = i
		}
	}
	if firstLine < 0 {
		return // no focused widget visible
	}

	// Scroll up if needed
	if firstLine < c.scrollY {
		c.scrollY = firstLine
	}
	// Scroll down if needed
	if lastLine >= c.scrollY+innerH {
		c.scrollY = lastLine - innerH + 1
	}
	if c.scrollY < 0 {
		c.scrollY = 0
	}
}

// flush writes the screen buffer to stdout in one write.
func flush(buf *screenBuf) {
	var out bytes.Buffer
	out.Grow(buf.w * buf.h * 4) // rough estimate

	out.WriteString("\x1b[H") // move cursor to top-left
	currentInverse := false

	for y := 0; y < buf.h; y++ {
		out.WriteString(fmt.Sprintf("\x1b[%d;1H", y+1)) // move to line start
		for x := 0; x < buf.w; x++ {
			c := buf.cells[y][x]
			if c.inverse && !currentInverse {
				out.WriteString("\x1b[7m")
				currentInverse = true
			} else if !c.inverse && currentInverse {
				out.WriteString("\x1b[0m")
				currentInverse = false
			}
			out.WriteRune(c.ch)
		}
	}

	if currentInverse {
		out.WriteString("\x1b[0m")
	}

	os.Stdout.Write(out.Bytes())
}
