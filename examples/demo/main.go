package main

import (
	"fmt"

	"g0ui"
)

func main() {
	counter := 0
	w := g0ui.Widgets

	g0ui.Run(func() {
		g0ui.Begin("My program dashboard")

		w.Spacing()

		w.Text("Hello from g0ui example")
		w.Text(fmt.Sprintf("Counter: %d", counter))

		w.Spacing()

		if w.Button("Increment") {
			counter++
		}
		if w.Button("Decrement") {
			counter--
		}

		w.Spacing()
		w.Separation()
		w.Spacing()

		if w.Button("Quit") {
			g0ui.Quit()
		}

		g0ui.End()
	})
}
