package g0ui

// frameContext holds all state for one frame cycle.
type frameContext struct {
	title   string
	widgets []widget

	focusIndex    int // which focusable widget is focused
	focusCount    int // total number of focusable widgets (from previous frame)
	scrollY       int // vertical scroll offset in content lines
	termW, termH  int
	pressed       int        // focusID that was activated this frame (-1 = none)
	quit          bool
	running       bool
	input         InputEvent
	firstFrame    bool
}

var ctx frameContext

func init() {
	ctx.pressed = -1
	ctx.firstFrame = true
}
