package g0ui

import "strings"

// renderLine represents one line of rendered content.
type renderLine struct {
	cells     []cell
	focusID   int // focusID of the widget this line belongs to (-1 if none)
	lineInWidget int // which line within the widget (for scroll-to-focus)
}

// layoutWidgets converts widgets into render lines.
func layoutWidgets(c *frameContext) []renderLine {
	innerW := c.termW - 4 // 2 border + 2 padding
	if innerW < 4 {
		innerW = 4
	}

	var lines []renderLine
	inGroup := false
	var groupRow []groupBlock

	for _, w := range c.widgets {
		switch w.kind {
		case widgetText:
			if inGroup {
				// Text in a group: treat as a block
				textLines := wrapText(w.label, innerW)
				block := groupBlock{
					lines:   textLinesToCells(textLines, -1),
					width:   maxLineWidth(textLines),
					focusID: -1,
				}
				groupRow = append(groupRow, block)
			} else {
				textLines := wrapText(w.label, innerW)
				for i, tl := range textLines {
					cells := stringToCells(tl, innerW)
					lines = append(lines, renderLine{cells: cells, focusID: -1, lineInWidget: i})
				}
			}

		case widgetButton:
			btnW := runeCount(w.label) + 4 // "│ " + label + " │"
			if btnW > innerW {
				btnW = innerW
			}
			block := makeButtonBlock(w.label, btnW, w.focusID)
			if inGroup {
				groupRow = append(groupRow, block)
			} else {
				for _, bl := range block.lines {
					lines = append(lines, bl)
				}
			}

		case widgetGroupStart:
			inGroup = true
			groupRow = nil

		case widgetGroupEnd:
			inGroup = false
			lines = append(lines, layoutGroupRow(groupRow, innerW)...)
			groupRow = nil
		}
	}

	// If group was not closed, flush it
	if inGroup && len(groupRow) > 0 {
		lines = append(lines, layoutGroupRow(groupRow, innerW)...)
	}

	return lines
}

// groupBlock is a rectangular block of cells for horizontal layout.
type groupBlock struct {
	lines   []renderLine
	width   int
	focusID int
}

// makeButtonBlock creates a 3-line button block.
func makeButtonBlock(label string, width int, focusID int) groupBlock {
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
	mid[1] = cell{ch: ' '}
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

	return groupBlock{
		lines: []renderLine{
			{cells: top, focusID: focusID, lineInWidget: 0},
			{cells: mid, focusID: focusID, lineInWidget: 1},
			{cells: bot, focusID: focusID, lineInWidget: 2},
		},
		width:   width,
		focusID: focusID,
	}
}

// layoutGroupRow arranges group blocks horizontally with wrapping.
func layoutGroupRow(blocks []groupBlock, maxW int) []renderLine {
	if len(blocks) == 0 {
		return nil
	}

	// Find max height of blocks in a row segment
	type rowSegment struct {
		blocks []groupBlock
	}

	var segments []rowSegment
	var current []groupBlock
	currentW := 0

	for _, b := range blocks {
		bw := b.width + 1 // +1 for spacing between blocks
		if len(current) > 0 && currentW+bw-1 > maxW {
			// Wrap to new row
			segments = append(segments, rowSegment{blocks: current})
			current = nil
			currentW = 0
		}
		current = append(current, b)
		if len(current) == 1 {
			currentW = b.width
		} else {
			currentW += bw
		}
	}
	if len(current) > 0 {
		segments = append(segments, rowSegment{blocks: current})
	}

	var result []renderLine
	for _, seg := range segments {
		// Find max height in this segment
		maxH := 0
		for _, b := range seg.blocks {
			if len(b.lines) > maxH {
				maxH = len(b.lines)
			}
		}

		// Merge blocks line by line
		for row := 0; row < maxH; row++ {
			merged := make([]cell, 0, maxW)
			mergedFocusID := -1
			for bi, b := range seg.blocks {
				if bi > 0 {
					merged = append(merged, cell{ch: ' '}) // gap
				}
				if row < len(b.lines) {
					merged = append(merged, b.lines[row].cells...)
					if b.lines[row].focusID >= 0 {
						mergedFocusID = b.lines[row].focusID
					}
				} else {
					// Pad with spaces
					for i := 0; i < b.width; i++ {
						merged = append(merged, cell{ch: ' '})
					}
				}
			}
			result = append(result, renderLine{cells: merged, focusID: mergedFocusID, lineInWidget: row})
		}
	}

	return result
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

func textLinesToCells(lines []string, focusID int) []renderLine {
	var result []renderLine
	for i, l := range lines {
		cells := make([]cell, len([]rune(l)))
		for j, r := range []rune(l) {
			cells[j] = cell{ch: r}
		}
		result = append(result, renderLine{cells: cells, focusID: focusID, lineInWidget: i})
	}
	return result
}

func runeCount(s string) int {
	return len([]rune(s))
}

func maxLineWidth(lines []string) int {
	max := 0
	for _, l := range lines {
		n := runeCount(l)
		if n > max {
			max = n
		}
	}
	return max
}
