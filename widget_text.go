package g0ui

// Text adds a non-focusable text widget.
func (W) Text(s string) {
	ctx.widgets = append(ctx.widgets, widget{
		kind:      widgetText,
		label:     s,
		focusable: false,
		focusID:   -1,
	})
}
