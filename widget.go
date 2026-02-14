package g0ui

type widgetKind int

const (
	widgetText widgetKind = iota
	widgetButton
	widgetGroupStart
	widgetGroupEnd
)

type widget struct {
	kind      widgetKind
	label     string
	focusable bool
	focusID   int // index among focusable widgets (-1 if not focusable)
}
