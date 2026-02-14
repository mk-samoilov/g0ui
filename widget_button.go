package g0ui

// Button adds a focusable button widget. Returns true if pressed this frame.
func (W) Button(label string) bool {
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
