package g0ui

import (
	"os"
	"strings"
)

// Run initializes the terminal, runs fn in a loop, and restores on exit.
func Run(fn func()) {
	if err := enableRawMode(); err != nil {
		os.Stderr.WriteString("g0ui: failed to enable raw mode: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer disableRawMode()

	enterAltScreen()
	defer exitAltScreen()

	hideCursor()
	defer showCursor()

	clearScreen()

	ctx.running = true
	ctx.quit = false
	ctx.firstFrame = true

	for ctx.running && !ctx.quit {
		fn()
	}
}

// Quit signals the main loop to stop.
func Quit() {
	ctx.quit = true
}

// Begin starts a new frame. It reads input and processes navigation.
func Begin(title string) {
	// Get terminal size
	ctx.termW, ctx.termH = getTermSize()
	ctx.title = title

	// Reset widgets for this frame
	ctx.widgets = ctx.widgets[:0]
	ctx.pressed = -1

	// Read input (blocking) — skip on first frame and after button press to render immediately
	if ctx.firstFrame {
		ctx.input = InputEvent{Key: KeyNone}
		ctx.firstFrame = false
	} else if ctx.needRedraw {
		ctx.input = InputEvent{Key: KeyNone}
		ctx.needRedraw = false
	} else {
		ctx.input = readInput()
	}

	// Handle global keys
	switch ctx.input.Key {
	case KeyCtrlC:
		ctx.quit = true
		return
	}

	// Handle navigation
	switch ctx.input.Key {
	case KeyUp:
		if ctx.focusIndex > 0 {
			ctx.focusIndex--
		}
	case KeyDown:
		if ctx.focusIndex < ctx.focusCount-1 {
			ctx.focusIndex++
		}
	case KeyTab:
		if ctx.focusIndex < ctx.focusCount-1 {
			ctx.focusIndex++
		} else {
			ctx.focusIndex = 0
		}
	case KeyEnter, KeySpace:
		ctx.pressed = ctx.focusIndex
	}
}

// End finishes the frame: performs layout and renders.
func End() {
	if ctx.quit {
		return
	}

	// Count focusable widgets for next frame's navigation bounds
	count := 0
	for _, w := range ctx.widgets {
		if w.focusable {
			count++
		}
	}
	ctx.focusCount = count

	// Clamp focus
	if ctx.focusIndex >= ctx.focusCount {
		ctx.focusIndex = ctx.focusCount - 1
	}
	if ctx.focusIndex < 0 {
		ctx.focusIndex = 0
	}

	// After button press, schedule immediate redraw next frame
	if ctx.pressed >= 0 {
		ctx.needRedraw = true
	}

	// Skip render if there's more buffered input
	if !hasBufferedInput() || ctx.pressed >= 0 {
		renderFrame(&ctx)
	}
}

// Text adds a non-focusable text widget.
func Text(s string) {
	ctx.widgets = append(ctx.widgets, widget{
		kind:      widgetText,
		label:     s,
		focusable: false,
		focusID:   -1,
	})
}

// Button adds a focusable button widget. Returns true if pressed this frame.
func Button(label string) bool {
	// Assign focusID
	fid := 0
	for _, w := range ctx.widgets {
		if w.focusable {
			fid++
		}
	}

	ctx.widgets = append(ctx.widgets, widget{
		kind:      widgetButton,
		label:     label,
		focusable: true,
		focusID:   fid,
	})

	return ctx.pressed == fid
}

// Break adds an empty line.
func Break() {
	Text("")
}

// Separation adds a horizontal line of dashes. Default length is 6.
func Separation(el ...int) {
	n := 6
	if len(el) > 0 && el[0] > 0 {
		n = el[0]
	}
	Text(strings.Repeat("--", n))
}

// BeginGroup starts a horizontal group of widgets.
func BeginGroup() {
	ctx.widgets = append(ctx.widgets, widget{
		kind:      widgetGroupStart,
		focusable: false,
		focusID:   -1,
	})
}

// EndGroup ends a horizontal group of widgets.
func EndGroup() {
	ctx.widgets = append(ctx.widgets, widget{
		kind:      widgetGroupEnd,
		focusable: false,
		focusID:   -1,
	})
}
