package main

import (
	"fmt"
	"time"

	"g0ui"
)

func main() {
	w := g0ui.Widgets

	counter := 0
	var clickLog []string

	g0ui.Run(func() {
		g0ui.Begin("g0ui Demo App")

		w.Spacing()
		w.Text("Welcome to g0ui demo application!")
		w.Text(fmt.Sprintf("Current time: %s", time.Now().Format("15:04:05")))
		w.Spacing()

		w.Separation(30)
		w.Spacing()

		w.Text(fmt.Sprintf("Counter value: %d", counter))
		w.Spacing()

		if w.Button("+ Increment") {
			counter++
			clickLog = append(clickLog,
				fmt.Sprintf("[%s] Incremented to %d", time.Now().Format("15:04:05"), counter))
		}
		if w.Button("- Decrement") {
			counter--
			clickLog = append(clickLog,
				fmt.Sprintf("[%s] Decremented to %d", time.Now().Format("15:04:05"), counter))
		}
		if w.Button("Reset") {
			counter = 0
			clickLog = append(clickLog,
				fmt.Sprintf("[%s] Reset to 0", time.Now().Format("15:04:05")))
		}

		w.Spacing()
		w.Separation(32)
		w.Spacing()

		w.Text("Action log:")
		if len(clickLog) == 0 {
			w.Text("  (no actions yet)")
		} else {
			start := 0
			if len(clickLog) > 5 {
				start = len(clickLog) - 5
			}
			for _, entry := range clickLog[start:] {
				w.Text("  " + entry)
			}
		}

		w.Spacing()
		w.Separation(32)
		w.Spacing()

		if w.Button("Quit") {
			g0ui.Quit()
		}

		g0ui.End()
	})
}
