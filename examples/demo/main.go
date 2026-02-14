package main

import (
	"fmt"
	"g0ui"
)

func main() {
	counter := 0

	g0ui.Run(func() {
		g0ui.Begin("My program dashboard")

		g0ui.Spacing()

		g0ui.Text("Hello from g0ui example")
		g0ui.Text(fmt.Sprintf("Counter: %d", counter))

		g0ui.Spacing()

		if g0ui.Button("Increment") {
			counter++
		}
		if g0ui.Button("Decrement") {
			counter--
		}

		g0ui.Spacing()
		g0ui.Separation()
		g0ui.Spacing()

		if g0ui.Button("Quit") {
			g0ui.Quit()
		}

		g0ui.End()
	})
}
