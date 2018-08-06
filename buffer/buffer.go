package buffer

import (
	"github.com/sirupsen/logrus"
)

type Buffer struct {
	lines      []Line
	cursorX    uint16
	cursorY    uint16
	viewHeight uint16
	viewWidth  uint16
	cursorAttr CellAttributes
}

// NewBuffer creates a new terminal buffer
func NewBuffer(viewCols uint16, viewLines uint16) *Buffer {
	b := &Buffer{
		cursorX: 0,
		cursorY: 0,
		lines:   []Line{},
	}
	b.ResizeView(viewCols, viewLines)
	return b
}

// Column returns cursor column
func (buffer *Buffer) CursorColumn() uint16 {
	return buffer.cursorX
}

// Line returns cursor line
func (buffer *Buffer) CursorLine() uint16 {
	return buffer.cursorY
}

// translates the cursor line to the raw buffer line
func (buffer *Buffer) RawLine() uint64 {
	rawHeight := buffer.Height()
	if int(buffer.viewHeight) > rawHeight {
		return uint64(buffer.cursorY)
	}
	return uint64(int(buffer.cursorY) + (rawHeight - int(buffer.viewHeight)))
}

// Width returns the width of the buffer in columns
func (buffer *Buffer) Width() uint16 {
	return buffer.viewWidth
}

func (buffer *Buffer) ViewWidth() uint16 {
	return buffer.viewWidth
}

func (buffer *Buffer) Height() int {
	return len(buffer.lines)
}

func (buffer *Buffer) ViewHeight() uint16 {
	return buffer.viewHeight
}

func (buffer *Buffer) ensureLinesExistToRawHeight() {
	for int(buffer.RawLine()) >= len(buffer.lines) {
		buffer.lines = append(buffer.lines, newLine())
	}
}

// Write will write a rune to the terminal at the position of the cursor, and increment the cursor position
func (buffer *Buffer) Write(runes ...rune) {
	for _, r := range runes {
		buffer.ensureLinesExistToRawHeight()
		if r == 0x0a {
			buffer.NewLine()
			continue
		}
		line := &buffer.lines[buffer.RawLine()]
		for int(buffer.CursorColumn()) >= len(line.cells) {
			line.cells = append(line.cells, newCell())
		}
		cell := &line.cells[buffer.CursorColumn()]
		cell.setRune(r)
		cell.attr = buffer.cursorAttr
		buffer.incrementCursorPosition()
	}
}

func (buffer *Buffer) incrementCursorPosition() {

	if buffer.CursorColumn()+1 < buffer.Width() {
		buffer.cursorX++
	} else {
		if buffer.cursorY == buffer.viewHeight-1 { // if we're on the last line, we can't move the cursor down, we have to move the buffer up, i.e. add a new line
			line := newLine()
			line.setWrapped(true)
			buffer.lines = append(buffer.lines, line)
			buffer.cursorX = 0
		} else {
			buffer.cursorX = 0
			if buffer.Height() < int(buffer.ViewHeight()) {
				line := newLine()
				line.setWrapped(true)
				buffer.lines = append(buffer.lines, line)
				buffer.cursorY++
			} else {
				panic("no test for this yet - not sure if possible?")
				line := &buffer.lines[buffer.RawLine()]
				line.setWrapped(true)
			}
		}
	}
}

func (buffer *Buffer) NewLine() {
	// if we're at the beginning of a line which wrapped from the previous one, and we need a new line, we can effectively not add a new line, and set the current one to non-wrapped
	if buffer.cursorX == 0 {
		line := &buffer.lines[buffer.RawLine()]
		if line.wrapped {
			line.setWrapped(false)
			return
		}
	}

	if buffer.cursorY == buffer.viewHeight-1 {
		buffer.lines = append(buffer.lines, newLine())
		buffer.cursorX = 0
	} else {
		buffer.cursorX = 0
		buffer.cursorY++
	}
}

func (buffer *Buffer) MovePosition(x int16, y int16) {

	if int16(buffer.cursorX)+x < 0 {
		x = -int16(buffer.cursorX)
	}

	if int16(buffer.cursorY)+y < 0 {
		y = -int16(buffer.cursorY)
	}

	buffer.SetPosition(uint16(int16(buffer.cursorX)+x), uint16(int16(buffer.cursorY)+y))
}

func (buffer *Buffer) SetPosition(col uint16, line uint16) {
	if col >= buffer.ViewWidth() {
		col = buffer.ViewWidth() - 1
		logrus.Errorf("Cannot set cursor position: column %d is outside of the current view width (%d columns)", col, buffer.ViewWidth())
	}
	if line >= buffer.ViewHeight() {
		line = buffer.ViewHeight() - 1
		logrus.Errorf("Cannot set cursor position: line %d is outside of the current view height (%d lines)", line, buffer.ViewHeight())
	}
	buffer.cursorX = col
	buffer.cursorY = line
}

func (buffer *Buffer) GetVisibleLines() []Line {
	lines := []Line{}
	for i := buffer.Height() - int(buffer.ViewHeight()); i < buffer.Height(); i++ {
		if i >= 0 && i < len(buffer.lines) {
			lines = append(lines, buffer.lines[i])
		}
	}
	return lines
}

// tested to here

func (buffer *Buffer) Clear() {
	for i := 0; i < int(buffer.ViewHeight()); i++ {
		buffer.lines = append(buffer.lines, newLine())
	}
	buffer.SetPosition(0, 0)
}

func (buffer *Buffer) ResizeView(width uint16, height uint16) {
	buffer.viewWidth = width
	buffer.viewHeight = height

	// @todo wrap/unwrap
}
