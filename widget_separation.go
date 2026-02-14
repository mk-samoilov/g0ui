package g0ui

import "strings"

// Separation adds a horizontal line of dashes. Default length is 14.
func (w W) Separation(el ...int) {
	n := 14
	if len(el) > 0 && el[0] > 0 {
		n = el[0]
	}

	w.Text(strings.Repeat("-", n))
}
