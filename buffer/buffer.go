package buffer

import (
	"fmt"
)

type Buffer struct {
	lines       []line
	x           int
	y           int
	columnCount int
	lineCount   int
}

// NewBuffer creates a new terminal buffer
func NewBuffer(columns int) *Buffer {
	return &Buffer{
		x:           0,
		y:           0,
		lines:       []line{},
		columnCount: columns,
	}
}

// Column returns cursor column
func (buffer *Buffer) Column() int {
	return buffer.x
}

// Line returns cursor line
func (buffer *Buffer) Line() int {
	return buffer.y
}

// Width returns the width of the buffer in columns
func (buffer *Buffer) Width() int {
	return buffer.columnCount
}

// Write will write a rune to the terminal at the position of the cursor, and increment the cursor position
func (buffer *Buffer) Write(r rune) {
	for buffer.Line() >= len(buffer.lines) {
		buffer.lines = append(buffer.lines, newLine())
	}
	line := &buffer.lines[buffer.Line()]
	for buffer.Column() >= len(line.cells) {
		line.cells = append(line.cells, newCell())
	}
	cell := line.cells[buffer.Column()]
	cell.setRune(r)
}

func (buffer *Buffer) incrementCursorPosition() {

	if buffer.Column()+1 < buffer.Width() {
		buffer.x++
	} else {
		buffer.y++
		buffer.x = 0
	}
}

func (buffer *Buffer) SetPosition(col int, line int) error {
	if buffer.x >= buffer.Width() {
		return fmt.Errorf("Cannot set cursor position: column %d is outside of the current buffer width (%d columns)", col, buffer.Width())
	}
	buffer.x = col
	buffer.y = line
	return nil
}

func (buffer *Buffer) SetSize(cols int, lines int) {

}
