package g0ui

import "strings"

// renderLine represents one line of rendered content.
type renderLine struct {
	cells        []cell
	focusID      int // focusID of the widget this line belongs to (-1 if none)
	lineInWidget int // which line within the widget (for scroll-to-focus)
}

// layoutWidgets converts widgets into render lines.
func layoutWidgets(c *frameContext) []renderLine {
	innerW := c.termW - 4 // 2 border + 2 padding
	if innerW < 4 {
		innerW = 4
	}

	var lines []renderLine

	for _, w := range c.widgets {
		switch w.kind {
		case widgetText:
			textLines := wrapText(w.label, innerW)
			for i, tl := range textLines {
				cells := stringToCells(tl, innerW)
				lines = append(lines, renderLine{cells: cells, focusID: -1, lineInWidget: i})
			}

		case widgetButton:
			btnW := runeCount(w.label) + 4 // "│ " + label + " │"
			if btnW > innerW {
				btnW = innerW
			}
			block := makeButtonBlock(w.label, btnW, w.focusID)
			for _, bl := range block.lines {
				lines = append(lines, bl)
			}
		}
	}

	return lines
}

// buttonBlock is a rectangular block of cells for a button.
type buttonBlock struct {
	lines   []renderLine
	width   int
	focusID int
}

// makeButtonBlock creates a 3-line button block.
func makeButtonBlock(label string, width int, focusID int) buttonBlock {
	innerW := width - 2 // subtract "│" on each side
	if innerW < 1 {
		innerW = 1
	}

	// Top: ┌─────┐
	top := make([]cell, width)
	top[0] = cell{ch: '┌'}
	for i := 1; i < width-1; i++ {
		top[i] = cell{ch: '─'}
	}
	top[width-1] = cell{ch: '┐'}

	// Middle: │ label │
	mid := make([]cell, width)
	mid[0] = cell{ch: '│'}
	mid[1] = cell{ch: ' ', focusInverse: true}
	labelRunes := []rune(label)
	for i := 0; i < innerW-2; i++ {
		if i < len(labelRunes) {
			mid[i+2] = cell{ch: labelRunes[i], focusInverse: true}
		} else {
			mid[i+2] = cell{ch: ' ', focusInverse: true}
		}
	}
	mid[width-2] = cell{ch: ' ', focusInverse: true}
	mid[width-1] = cell{ch: '│'}

	// Bottom: └─────┘
	bot := make([]cell, width)
	bot[0] = cell{ch: '└'}
	for i := 1; i < width-1; i++ {
		bot[i] = cell{ch: '─'}
	}
	bot[width-1] = cell{ch: '┘'}

	return buttonBlock{
		lines: []renderLine{
			{cells: top, focusID: focusID, lineInWidget: 0},
			{cells: mid, focusID: focusID, lineInWidget: 1},
			{cells: bot, focusID: focusID, lineInWidget: 2},
		},
		width:   width,
		focusID: focusID,
	}
}

// wrapText splits text into lines respecting \n and wrapping at maxW.
func wrapText(s string, maxW int) []string {
	if maxW < 1 {
		maxW = 1
	}
	var result []string
	paragraphs := strings.Split(s, "\n")
	for _, para := range paragraphs {
		if len(para) == 0 {
			result = append(result, "")
			continue
		}
		runes := []rune(para)
		for len(runes) > 0 {
			if len(runes) <= maxW {
				result = append(result, string(runes))
				break
			}
			// Find a break point (space)
			breakAt := maxW
			for i := maxW; i > 0; i-- {
				if runes[i] == ' ' {
					breakAt = i
					break
				}
			}
			result = append(result, string(runes[:breakAt]))
			runes = runes[breakAt:]
			// Skip leading space after break
			if len(runes) > 0 && runes[0] == ' ' {
				runes = runes[1:]
			}
		}
	}
	return result
}

// stringToCells converts a string to a fixed-width cell slice.
func stringToCells(s string, width int) []cell {
	cells := make([]cell, width)
	runes := []rune(s)
	for i := 0; i < width; i++ {
		if i < len(runes) {
			cells[i] = cell{ch: runes[i]}
		} else {
			cells[i] = cell{ch: ' '}
		}
	}
	return cells
}

func runeCount(s string) int {
	return len([]rune(s))
}
