package gui

import (
	"fmt"

	"github.com/liamg/aminal/buffer"
)

func (gui *GUI) textbox(col uint16, row uint16, text string, fg [3]float32, bg [3]float32) {

	lines := []string{}
	line := ""
	word := ""

	maxWidth := int(gui.terminal.ActiveBuffer().ViewWidth()) - 4
	maxHeight := (int(gui.terminal.ActiveBuffer().ViewHeight()) / 2) - 2

	if maxHeight < 1 {
		return
	}

	longestLine := 0

	addWord := func() {
		if len(line)+len(word) <= maxWidth {
			line = fmt.Sprintf("%s%s", line, word)
			if len(line) < maxWidth {
				line = fmt.Sprintf("%s ", line)
			} else {
				lines = append(lines, line)
				line = ""
			}
		} else {
			lines = append(lines, line)
			line = word
			for len(line) > maxWidth {
				// break word into bits
			}
		}

		word = ""
	}

	addLine := func() bool {
		addWord()
		if len(line) > longestLine {
			longestLine = len(line)
		}
		lines = append(lines, line)
		if len(lines) >= maxHeight-1 {
			lines = append(lines, "...")
			return true
		}
		line = ""
		return false
	}

	var done = false

DONE:
	for _, c := range text {
		switch c {
		case 0x0d:
			continue
		case 0x0a:
			if done = addLine(); done {
				break DONE
			}
		case ' ':
			addWord()
		default:
			word = fmt.Sprintf("%s%c", word, c)
		}
	}
	if word != "" {
		addWord()
	}
	if line != "" && !done {
		addLine()
	}

	for hx := col; hx < col+uint16(longestLine)+1; hx++ {
		for hy := row - 1; hy < row+uint16(len(lines))+1; hy++ {
			gui.renderer.DrawCellBg(buffer.NewBackgroundCell(bg), uint(hx), uint(hy), nil, true)
		}
	}

	x := float32(col) * gui.renderer.cellWidth

	f := gui.fontMap.DefaultFont()
	f.SetColor(fg[0], fg[1], fg[2], 1)

	for i, line := range lines {
		y := float32(row+1+uint16(i))*gui.renderer.cellHeight + f.MinY()
		f.Print(x, y, fmt.Sprintf(" %s", line))
	}

}
