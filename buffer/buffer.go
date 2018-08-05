package buffer

import (
	"fmt"
)

type Buffer struct {
	lines       []line
	x           uint16
	y           uint16
	columnCount uint16
	viewHeight  uint16
}

// NewBuffer creates a new terminal buffer
func NewBuffer() *Buffer {
	return &Buffer{
		x:           0,
		y:           0,
		lines:       []line{},
		columnCount: 0,
	}
}

// Column returns cursor column
func (buffer *Buffer) Column() uint16 {
	return buffer.x
}

// Line returns cursor line
func (buffer *Buffer) Line() uint16 {
	return buffer.y
}

// Width returns the width of the buffer in columns
func (buffer *Buffer) Width() uint16 {
	return buffer.columnCount
}

// Write will write a rune to the terminal at the position of the cursor, and increment the cursor position
func (buffer *Buffer) Write(r rune) {
	for int(buffer.Line()) >= len(buffer.lines) {
		buffer.lines = append(buffer.lines, newLine())
	}
	line := &buffer.lines[buffer.Line()]
	for int(buffer.Column()) >= len(line.cells) {
		line.cells = append(line.cells, newCell())
	}
	cell := line.cells[buffer.Column()]
	cell.setRune(r)
	buffer.incrementCursorPosition()
}

func (buffer *Buffer) incrementCursorPosition() {

	if buffer.Column()+1 < buffer.Width() {
		buffer.x++
	} else {
		buffer.y++
		buffer.x = 0
	}
}

func (buffer *Buffer) SetPosition(col uint16, line uint16) error {
	if buffer.x >= buffer.Width() {
		return fmt.Errorf("Cannot set cursor position: column %d is outside of the current buffer width (%d columns)", col, buffer.Width())
	}
	buffer.x = col
	buffer.y = line
	return nil
}

func (buffer *Buffer) Resize(cols int, lines int) {

}
